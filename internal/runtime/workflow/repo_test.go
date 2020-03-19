package workflow

import (
	"context"
	"reflect"
	"testing"
	"time"

	svc "github.com/hanjunlee/argocui/internal/runtime"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/pkg/client/clientset/versioned"
	af "github.com/argoproj/argo/pkg/client/clientset/versioned/fake"
	ai "github.com/argoproj/argo/pkg/client/informers/externalversions"
	listerv1alpha1 "github.com/argoproj/argo/pkg/client/listers/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	kf "k8s.io/client-go/kubernetes/fake"
)

func newMockingRepo(ws ...*wf.Workflow) *Repo {
	const (
		noResyncPeriod = 0
	)
	var (
		neverStop <-chan struct{}
	)
	// clientset
	objs := []runtime.Object{}
	for _, w := range ws {
		objs = append(objs, w)
	}
	ac := af.NewSimpleClientset(objs...)

	// informer
	af := ai.NewSharedInformerFactory(ac, noResyncPeriod)
	i := af.Argoproj().V1alpha1().Workflows()
	for _, w := range ws {
		i.Informer().GetIndexer().Add(w)
	}
	go i.Informer().Run(neverStop)

	pods := make([]runtime.Object, 0)
	for _, w := range ws {
		for _, n := range w.Status.Nodes {
			if n.Type == wf.NodeTypePod && n.Phase != wf.NodeError {
				namespace, n := w.GetNamespace(), n.ID
				pods = append(pods, newPod(namespace, n))
			}
		}
	}
	kc := kf.NewSimpleClientset(pods...)

	return &Repo{
		ac: ac,
		al: i.Lister(),
		kc: kc,
	}
}

func newPod(namespace, name string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				corev1.Container{
					Name: "main",
				},
			},
		},
	}
}

func TestRepo_Get(t *testing.T) {
	type fields struct {
		ac versioned.Interface
		al listerv1alpha1.WorkflowLister
		kc kubernetes.Interface
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		repo    *Repo
		args    args
		want    runtime.Object
		wantErr bool
	}{
		{
			name: "get a workflow",
			repo: newMockingRepo(&wf.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "hello-world-",
				},
				Status: wf.WorkflowStatus{
					Nodes: map[string]wf.NodeStatus{
						"hello-world-0": wf.NodeStatus{
							Type:  wf.NodeTypePod,
							Phase: wf.NodeSucceeded,
							ID:    "hello-world-0",
						},
					},
				},
			}),
			args: args{
				key: "default/hello-world-",
			},
			want: &wf.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "hello-world-",
				},
				Status: wf.WorkflowStatus{
					Nodes: map[string]wf.NodeStatus{
						"hello-world-0": wf.NodeStatus{
							Type:  wf.NodeTypePod,
							Phase: wf.NodeSucceeded,
							ID:    "hello-world-0",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "not exist",
			repo: newMockingRepo(&wf.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "hello-world-",
				},
				Status: wf.WorkflowStatus{
					Nodes: map[string]wf.NodeStatus{
						"hello-world-0": wf.NodeStatus{
							Type:  wf.NodeTypePod,
							Phase: wf.NodeSucceeded,
							ID:    "hello-world-0",
						},
					},
				},
			}),
			args: args{
				key: "default/not-exist",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			got, err := r.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Search(t *testing.T) {
	std := time.Now()
	type args struct {
		namespace string
		pattern   string
	}
	tests := []struct {
		name string
		repo *Repo
		args args
		want []runtime.Object
	}{
		{
			name: "search workflows",
			repo: newMockingRepo(
				&wf.Workflow{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:         "default",
						Name:              "wf-finished",
						CreationTimestamp: metav1.NewTime(std),
					},
					Status: wf.WorkflowStatus{
						FinishedAt: metav1.NewTime(std.Add(time.Minute)),
					},
				},
				&wf.Workflow{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:         "default",
						Name:              "wf-running",
						CreationTimestamp: metav1.NewTime(std),
					},
				},
			),
			args: args{
				namespace: "default",
				pattern:   "wf-",
			},
			want: []runtime.Object{
				&wf.Workflow{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:         "default",
						Name:              "wf-running",
						CreationTimestamp: metav1.NewTime(std),
					},
				},
				&wf.Workflow{
					ObjectMeta: metav1.ObjectMeta{
						Namespace:         "default",
						Name:              "wf-finished",
						CreationTimestamp: metav1.NewTime(std),
					},
					Status: wf.WorkflowStatus{
						FinishedAt: metav1.NewTime(std.Add(time.Minute)),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			if got := r.Search(tt.args.namespace, tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repo.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Delete(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		repo    *Repo
		args    args
		wantErr bool
	}{
		{
			name: "delete a workflow",
			repo: newMockingRepo(&wf.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "hello-world-",
				},
			}),
			args: args{
				key: "default/hello-world-",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			if err := r.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Repo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRepo_Logs(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name          string
		repo          *Repo
		args          args
		wantCntAction int
		wantErr       bool
	}{
		{
			name: "log multiple pods",
			repo: newMockingRepo(&wf.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "hello-world-",
				},
				Status: wf.WorkflowStatus{
					Nodes: map[string]wf.NodeStatus{
						"hello-world-0": wf.NodeStatus{
							Type:        wf.NodeTypePod,
							Phase:       wf.NodeSucceeded,
							ID:          "hello-world-0",
							DisplayName: "hello-world",
						},
						"hello-world-1": wf.NodeStatus{
							Type:        wf.NodeTypePod,
							Phase:       wf.NodeSucceeded,
							ID:          "hello-world-1",
							DisplayName: "hello-world",
						},
					},
				},
			}),
			args: args{
				ctx: context.Background(),
				key: "default/hello-world-",
			},
			wantCntAction: 2,
			wantErr:       false,
		},
		{
			name: "an error node exist",
			repo: newMockingRepo(&wf.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "hello-world-",
				},
				Status: wf.WorkflowStatus{
					Nodes: map[string]wf.NodeStatus{
						"hello-world-0": wf.NodeStatus{
							Type:        wf.NodeTypePod,
							Phase:       wf.NodeSucceeded,
							ID:          "hello-world-0",
							DisplayName: "hello-world",
						},
						"hello-world-1": wf.NodeStatus{
							Type:        wf.NodeTypePod,
							Phase:       wf.NodeError, // error exist
							ID:          "hello-world-1",
							DisplayName: "hello-world",
						},
					},
				},
			}),
			args: args{
				ctx: context.Background(),
				key: "default/hello-world-",
			},
			wantCntAction: 1,
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			_, err := r.Logs(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repo.Logs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// wait for logging pods.
			time.Sleep(500 * time.Millisecond)
			if actions := r.kc.(*kf.Clientset).Actions(); len(actions) != tt.wantCntAction {
				t.Errorf("Repo.Logs() = %v, want %v", len(actions), tt.wantCntAction)
			}
		})
	}
}

func TestRepo_logsPod(t *testing.T) {
	type fields struct {
		ac versioned.Interface
		al listerv1alpha1.WorkflowLister
		kc kubernetes.Interface
	}
	type args struct {
		ctx context.Context
		ch  chan<- svc.Log
		ns  string
		n   string
		dn  string
	}
	tests := []struct {
		name    string
		repo    *Repo
		args    args
		wantErr bool
	}{
		{
			name: "log a pod",
			repo: newMockingRepo(&wf.Workflow{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: "default",
					Name:      "hello-world-",
				},
				Status: wf.WorkflowStatus{
					Nodes: map[string]wf.NodeStatus{
						"hello-world-0": wf.NodeStatus{
							Type:        wf.NodeTypePod,
							Phase:       wf.NodeSucceeded,
							ID:          "hello-world-0",
							DisplayName: "hello-world",
						},
					},
				},
			}),
			args: args{
				ctx: context.Background(),
				ch:  make(chan<- svc.Log),
				ns:  "default",
				n:   "hello-world-0",
				dn:  "hello-world",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.repo
			if err := r.logsPod(tt.args.ctx, tt.args.ch, tt.args.ns, tt.args.n, tt.args.dn); (err != nil) != tt.wantErr {
				t.Errorf("Repo.logsPod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if actions := r.kc.(*kf.Clientset).Actions(); len(actions) != 1 {
				t.Errorf("Repo.logsPod() len(actions) != 1, %v", actions)
			}
		})
	}
}

func Test_splitTimeAndMessage(t *testing.T) {
	type args struct {
		l string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitTimeAndMessage(tt.args.l)
			if got != tt.want {
				t.Errorf("splitTimeAndMessage() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("splitTimeAndMessage() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getNodeDisplayName(t *testing.T) {
	type args struct {
		n wf.NodeStatus
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getNodeDisplayName(tt.args.n); got != tt.want {
				t.Errorf("getNodeDisplayName() = %v, want %v", got, tt.want)
			}
		})
	}
}
