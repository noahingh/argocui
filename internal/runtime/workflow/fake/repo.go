package repo

import (
	"context"
	"fmt"
	"sort"
	"strings"

	svc "github.com/hanjunlee/argocui/internal/runtime"
	repo "github.com/hanjunlee/argocui/internal/runtime/workflow"
	err "github.com/hanjunlee/argocui/pkg/util/error"

	"github.com/argoproj/argo/pkg/client/clientset/versioned"
	listerv1alpha1 "github.com/argoproj/argo/pkg/client/listers/workflow/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

// Repo is the workflow repository it always syncronizes the workflows as the storage at background.
type Repo struct {
	ac versioned.Interface
	ai cache.SharedIndexInformer
	al listerv1alpha1.WorkflowLister
}

// NewRepo create a new workflow repository.
func NewRepo(ac versioned.Interface, ai cache.SharedIndexInformer, al listerv1alpha1.WorkflowLister) *Repo {
	repo := &Repo{
		ac: ac,
		ai: ai,
		al: al,
	}
	return repo
}

// Get get the workflow by the key, the format is "namespace/key", and if doesn't exist it return nil.
func (r *Repo) Get(key string) (runtime.Object, error) {
	ns, n, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return nil, err
	}

	w, err := r.al.Workflows(ns).Get(n)
	if err != nil {
		return nil, err
	}

	return w, nil
}

// Search get workflows which are matched with pattern.
func (r *Repo) Search(namespace, pattern string) []runtime.Object {
	var (
		ret = make([]runtime.Object, 0)
	)

	ws, err := r.al.Workflows(namespace).List(labels.Everything())
	if err != nil {
		return ret
	}
	sort.Sort(repo.Workflows(ws))

	for _, w := range ws {
		name := w.GetName()
		if pattern == "" {
			ret = append(ret, w)
			continue
		}

		if i := strings.Index(name, pattern); i != -1 {
			ret = append(ret, w)
		}
	}
	return ret
}

// Delete delete the workflow by the key.
func (r *Repo) Delete(key string) error {
	if key == "" {
		return fmt.Errorf("there is no key to delete")
	}

	ns, n, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	// delete in the client
	err = r.ac.ArgoprojV1alpha1().Workflows(ns).Delete(n, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	// delete in the indexer
	i, _, _ := r.ai.GetIndexer().GetByKey(key)
	if err := r.ai.GetIndexer().Delete(i); err != nil {
		return err
	}

	return nil
}

// Logs is not implemented.
func (r *Repo) Logs(ctx context.Context, key string) (<-chan svc.Log, error) {
	return nil, err.ErrNotImplement
}
