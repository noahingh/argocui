package repo

import (
	"errors"
	"sync"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	informers "github.com/argoproj/argo/pkg/client/informers/externalversions/workflow/v1alpha1"
	lister "github.com/argoproj/argo/pkg/client/listers/workflow/v1alpha1"
	"github.com/ryanuber/go-glob"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/cache"
)

var (
	errNotExist error = errors.New("workflow does not exist")
)

type storage struct {
	list   []string
	lister lister.WorkflowLister
	mux    *sync.Mutex
	log    *log.Entry
}

func newStorage(informer informers.WorkflowInformer) *storage {
	return &storage{
		list:   []string{},
		lister: informer.Lister(),
		mux:    &sync.Mutex{},
		log: log.WithFields(log.Fields{
			"pkg":  "repo",
			"file": "storage.go",
		}),
	}
}

// List return keys which matchs with the pattern, it support the glob.
func (s *storage) List(pattern string) []string {
	s.mux.Lock()
	copied := make([]string, len(s.list))
	copy(copied, s.list)
	s.mux.Unlock()

	ret := make([]string, 0)
	for _, key := range copied {
		if !glob.Glob(pattern, key) {
			s.log.Tracef("the pattern doesn't match with '%s'.", key)
			continue
		}

		s.log.Tracef("append the '%s' key.", key)
		ret = append(ret, key)
	}
	return ret
}

// ReplaceOrInsert adds the index of given item to the list.
// otherwise it returns an empty string.
func (s *storage) ReplaceOrInsert(key string) string {
	s.mux.Lock()
	defer s.mux.Unlock()

	if s.has(key) {
		// delete the key and insert
		s.delete(key)
	}

	item, err := s.GetWorkflow(key)
	if err != nil {
		return ""
	}

	for i, wKey := range s.list {
		comp, err := s.GetWorkflow(wKey)
		if err != nil {
			// pass
			continue
		}

		if !less(item, comp) {
			continue
		}
		s.list = append(s.list, "")
		copy(s.list[i+1:], s.list[i:])
		s.list[i] = key
		return key
	}

	s.list = append(s.list, key)
	return key
}

func (s *storage) Delete(key string) string {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.delete(key)
}

// isExist confirm weather the same key is exist.
func (s *storage) Has(key string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.has(key)
}

func (s *storage) delete(key string) string {
	if !s.has(key) {
		return ""
	}
	i := s.index(key)
	s.list = append(s.list[:i], s.list[i+1:]...)

	return key
}

func (s *storage) has(key string) bool {
	if s.index(key) != -1 {
		return true
	}
	return false
}

// getIndex get the index of workflow.
func (s *storage) index(key string) int {
	ret := -1
	for i, wfName := range s.list {
		if key == wfName {
			ret = i
		}
	}
	return ret
}

func (s *storage) len() int { return len(s.list) }

func (s *storage) GetWorkflow(key string) (*wf.Workflow, error) {
	var namespace, name string
	var err error

	namespace, name, err = cache.SplitMetaNamespaceKey(key)
	w, err := s.lister.Workflows(namespace).Get(name)
	if err != nil {
		return nil, errNotExist
	}

	return w, nil
}

// less compare item and comp that the finished time of the item is close to now.
func less(item, comp *wf.Workflow) bool {
	iStart := item.ObjectMeta.CreationTimestamp
	iFinish := item.Status.FinishedAt
	cStart := comp.ObjectMeta.CreationTimestamp
	cFinish := comp.Status.FinishedAt

	if iFinish.IsZero() && cFinish.IsZero() {
		return cStart.Before(&iStart)
	}
	if iFinish.IsZero() && !cFinish.IsZero() {
		return true
	}
	if !iFinish.IsZero() && cFinish.IsZero() {
		return false
	}
	// comp finished eariler
	return cFinish.Before(&iFinish)
}
