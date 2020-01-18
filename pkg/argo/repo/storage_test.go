package repo

import (
	"reflect"
	"sync"
	"testing"
	"time"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/pkg/client/clientset/versioned/fake"
	informers "github.com/argoproj/argo/pkg/client/informers/externalversions"
	lister "github.com/argoproj/argo/pkg/client/listers/workflow/v1alpha1"
	log "github.com/sirupsen/logrus"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

func newLister(patterns []string) lister.WorkflowLister {
	var (
		wfs  = []*wf.Workflow{}
		objs = []runtime.Object{}
	)

	// create new workflows
	now := meta.Now()
	for delay, p := range patterns {
		namespace, name, _ := cache.SplitMetaNamespaceKey(p)

		w := &wf.Workflow{
			TypeMeta: meta.TypeMeta{
				APIVersion: "v1alpha1",
			},
			ObjectMeta: meta.ObjectMeta{
				Name:              name,
				Namespace:         namespace,
				CreationTimestamp: meta.NewTime(now.Add(time.Duration(delay) * time.Second)),
			},
		}
		wfs = append(wfs, w)
	}

	// create a new clientset
	for _, w := range wfs {
		objs = append(objs, w)
	}
	c := fake.NewSimpleClientset(objs...)

	// create a new informer
	const (
		noResync = 0
	)

	factory := informers.NewSharedInformerFactory(c, noResync)
	i := factory.Argoproj().V1alpha1().Workflows()

	for _, w := range wfs {
		i.Informer().GetIndexer().Add(w)
	}
	return i.Lister()
}

func Test_storage_List(t *testing.T) {
	type fields struct {
		list   []string
		lister lister.WorkflowLister
		mux    *sync.Mutex
		log    *log.Entry
	}
	type args struct {
		pattern string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
		{
			name: "wildcard",
			fields: fields{
				list: []string{
					"argo/first",
					"argo/second",
				},
				lister: newLister([]string{
					"argo/first",
					"argo/second",
				}),
				mux: &sync.Mutex{},
				log: log.NewEntry(log.New()),
			},
			args: args{
				pattern: "*",
			},
			want: []string{
				"argo/first",
				"argo/second",
			},
		},
		{
			name: "hello only",
			fields: fields{
				list: []string{
					"argo/foohello",
					"argo/hellofoo",
					"argo/bar",
				},
				lister: newLister([]string{
					"argo/foohello",
					"argo/hellofoo",
					"argo/bar",
				}),
				mux: &sync.Mutex{},
				log: log.NewEntry(log.New()),
			},
			args: args{
				pattern: "argo/*hello*",
			},
			want: []string{
				"argo/foohello",
				"argo/hellofoo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				list:   tt.fields.list,
				lister: tt.fields.lister,
				mux:    tt.fields.mux,
				log:    tt.fields.log,
			}
			if got := s.List(tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("storage.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_storage_ReplaceOrInsert(t *testing.T) {
	type fields struct {
		list   []string
		lister lister.WorkflowLister
		mux    *sync.Mutex
		log    *log.Entry
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
		{
			name: "insert a new",
			fields: fields{
				list: []string{
					"argo/second",
					"argo/first",
				},
				lister: newLister([]string{
					"argo/first",
					"argo/second",
					"argo/third",
				}),
				mux: &sync.Mutex{},
				log: log.NewEntry(log.New()),
			},
			args: args{
				key: "argo/third",
			},
			want: []string{
				"argo/third",
				"argo/second",
				"argo/first",
			},
		},
		{
			name: "insert sort",
			fields: fields{
				list: []string{
					"argo/third",
					"argo/first",
				},
				lister: newLister([]string{
					"argo/first",
					"argo/second",
					"argo/third",
				}),
				mux: &sync.Mutex{},
				log: log.NewEntry(log.New()),
			},
			args: args{
				key: "argo/second",
			},
			want: []string{
				"argo/third",
				"argo/second",
				"argo/first",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				list:   tt.fields.list,
				lister: tt.fields.lister,
				mux:    tt.fields.mux,
				log:    tt.fields.log,
			}
			if s.ReplaceOrInsert(tt.args.key); !reflect.DeepEqual(s.list, tt.want) {
				t.Errorf("storage.list = %v, want %v", s.list, tt.want)
			}
		})
	}
}

func Test_storage_Delete(t *testing.T) {
	type fields struct {
		list   []string
		lister lister.WorkflowLister
		mux    *sync.Mutex
		log    *log.Entry
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
		{
			name: "delete one",
			fields: fields{
				list: []string{
					"argo/third",
					"argo/second",
					"argo/first",
				},
				lister: newLister([]string{
					"argo/first",
					"argo/second",
					"argo/third",
				}),
				mux: &sync.Mutex{},
				log: log.NewEntry(log.New()),
			},
			args: args{
				key: "argo/third",
			},
			want: []string{
				"argo/second",
				"argo/first",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				list:   tt.fields.list,
				lister: tt.fields.lister,
				mux:    tt.fields.mux,
				log:    tt.fields.log,
			}
			if s.Delete(tt.args.key); !reflect.DeepEqual(s.list, tt.want) {
				t.Errorf("storage.list = %v, want %v", s.list, tt.want)
			}
		})
	}
}

func Test_storage_Has(t *testing.T) {
	type fields struct {
		list   []string
		lister lister.WorkflowLister
		mux    *sync.Mutex
		log    *log.Entry
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
		{
			name: "exist",
			fields: fields{
				list: []string{
					"argo/third",
					"argo/second",
					"argo/first",
				},
				lister: newLister([]string{
					"argo/first",
					"argo/second",
					"argo/third",
				}),
				mux: &sync.Mutex{},
				log: log.NewEntry(log.New()),
			},
			args: args{
				key: "argo/first",
			},
			want: true,
		},
		{
			name: "not exist",
			fields: fields{
				list: []string{
					"argo/third",
					"argo/second",
					"argo/first",
				},
				lister: newLister([]string{
					"argo/first",
					"argo/second",
					"argo/third",
				}),
				mux: &sync.Mutex{},
				log: log.NewEntry(log.New()),
			},
			args: args{
				key: "argo/foo",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &storage{
				list:   tt.fields.list,
				lister: tt.fields.lister,
				mux:    tt.fields.mux,
				log:    tt.fields.log,
			}
			if got := s.Has(tt.args.key); got != tt.want {
				t.Errorf("storage.Has() = %v, want %v", got, tt.want)
			}
		})
	}
}
