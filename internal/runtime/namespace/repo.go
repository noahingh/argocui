package namespace

import (
	"context"
	"sort"
	"strings"

	svc "github.com/hanjunlee/argocui/internal/runtime"
	err "github.com/hanjunlee/argocui/pkg/util/error"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	corev1lister "k8s.io/client-go/listers/core/v1"
)

type namespaces []*corev1.Namespace

func (n namespaces) Len() int {
	return len(n)
}

func (n namespaces) Less(i, j int) bool {
	in, jn := n[i], n[j]
	return in.GetName() < jn.GetName()
}

func (n namespaces) Swap(i, j int) {
	tmp := n[i]
	n[i] = n[j]
	n[j] = tmp
}

// Repo is the namespace repo.
type Repo struct {
	n corev1lister.NamespaceLister
}

// NewRepo create a new namespace repository.
func NewRepo(l corev1lister.NamespaceLister) *Repo {
	return &Repo{
		n: l,
	}
}

// Get return the namespace.
func (r *Repo) Get(key string) (runtime.Object, error) {
	return r.n.Get(key)
}

// Search return namespaces which is matched with the pattern.
func (r *Repo) Search(namespace, pattern string) []runtime.Object {
	objs := make([]runtime.Object, 0)

	nss, _ := r.n.List(labels.Everything())
	sort.Sort(namespaces(nss))

	for _, ns := range nss {
		if pattern == "" {
			objs = append(objs, ns)
			continue
		}

		name := ns.GetName()

		if idx := strings.Index(name, pattern); idx != -1 {
			objs = append(objs, ns)
		}
	}
	return objs
}

// Delete is not implemented.
func (r *Repo) Delete(key string) error {
	return err.ErrNotImplement
}

// Logs is not implemented.
func (r *Repo) Logs(ctx context.Context, key string) (<-chan svc.Log, error) {
	return nil, err.ErrNotImplement
}
