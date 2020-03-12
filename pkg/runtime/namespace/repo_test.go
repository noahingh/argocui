package namespace

import (
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes/fake"
	corev1lister "k8s.io/client-go/listers/core/v1"
)

func NewMockLister(objs ...runtime.Object) corev1lister.NamespaceLister {
	const (
		noResyncPeriod = 0
	)
	var (
		neverStop = make(<-chan struct{})
	)
	c := fake.NewSimpleClientset(objs...)
	f := informers.NewSharedInformerFactory(c, noResyncPeriod)

	i := f.Core().V1().Namespaces()
	for _, o := range objs {
		// add object into storage.
		i.Informer().GetIndexer().Add(o)
	}
	l := i.Lister()

	f.WaitForCacheSync(neverStop)

	return l
}

func NewNamespace(name string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func TestNewRepo(t *testing.T) {
	var (
		lister = NewMockLister(NewNamespace("foo"))
	)

	type args struct {
		l corev1lister.NamespaceLister
	}
	tests := []struct {
		name string
		args args
		want *Repo
	}{
		{
			name: "create a new",
			args: args{
				l: lister,
			},
			want: &Repo{
				n: lister,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRepo(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Get(t *testing.T) {
	type fields struct {
		n corev1lister.NamespaceLister
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    runtime.Object
		wantErr bool
	}{
		{
			name: "get a namespace",
			fields: fields{
				n: NewMockLister(NewNamespace("foo"), NewNamespace("bar")),
			},
			args: args{
				key: "foo",
			},
			want:    NewNamespace("foo"),
			wantErr: false,
		},
		{
			name: "not exist",
			fields: fields{
				n: NewMockLister(NewNamespace("foo"), NewNamespace("bar")),
			},
			args: args{
				key: "baz",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				n: tt.fields.n,
			}
			got, err := r.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Repo.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == true && tt.want == nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repo.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Search(t *testing.T) {
	type fields struct {
		n corev1lister.NamespaceLister
	}
	type args struct {
		namespace string
		pattern   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []runtime.Object
	}{
		{
			name: "all",
			fields: fields{
				n: NewMockLister(NewNamespace("foo"), NewNamespace("bar")),
			},
			args: args{
				namespace: "",
				pattern:   "",
			},
			want: []runtime.Object{NewNamespace("bar"), NewNamespace("foo")}, // sorted
		},
		{
			name: "pattern matched",
			fields: fields{
				n: NewMockLister(NewNamespace("foo"), NewNamespace("bar"), NewNamespace("baz")),
			},
			args: args{
				namespace: "",
				pattern:   "ba",
			},
			want: []runtime.Object{NewNamespace("bar"), NewNamespace("baz")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				n: tt.fields.n,
			}
			if got := r.Search(tt.args.namespace, tt.args.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repo.Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepo_Delete(t *testing.T) {
	type fields struct {
		n corev1lister.NamespaceLister
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "not implemented",
			fields: fields{
				n: NewMockLister(NewNamespace("foo")),
			},
			args: args{
				key: "foo",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repo{
				n: tt.fields.n,
			}
			if err := r.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Repo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
