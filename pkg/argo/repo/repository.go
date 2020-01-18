package repo

import (
	"fmt"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/pkg/client/clientset/versioned"
	informers "github.com/argoproj/argo/pkg/client/informers/externalversions/workflow/v1alpha1"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// ArgoRepository is the repository of Argo, it always syncronizes the workflows into storage at background.
type ArgoRepository struct {
	c   *controller
	s   *storage
	ac  versioned.Interface
	kc  kubernetes.Clientset
	log *log.Entry
}

// NewArgoRepository create a new Argo repository.
func NewArgoRepository(
	argoClientset versioned.Interface, argoInformer informers.WorkflowInformer, kubeClientset kubernetes.Clientset) *ArgoRepository {
	var (
		neverStop = make(chan struct{}, 0)
	)
	s := newStorage(argoInformer)
	c := newController(argoClientset, argoInformer, s)

	// run the controller
	go c.run(neverStop)
	c.waitForSynced(neverStop)
	repo := &ArgoRepository{
		c:  c,
		s:  s,
		ac: argoClientset,
		kc: kubeClientset,
		log: log.WithFields(log.Fields{
			"pkg":  "repo",
			"file": "repository.go",
		}),
	}
	return repo
}

// Get get the workflow by the key, the format is "namespace/key", and if doesn't exist it return nil.
func (a *ArgoRepository) Get(key string) *wf.Workflow {
	w, err := a.s.GetWorkflow(key)
	if err != nil {
		a.log.Errorf("failed to get the '%s' workflow.", key)
		return nil
	}

	return w
}

// Search get workflows which are matched with pattern.
func (a *ArgoRepository) Search(pattern string) []*wf.Workflow {
	var (
		wfs = []*wf.Workflow{}
	)
	keys := a.s.List(pattern)
	for _, k := range keys {
		w, err := a.s.GetWorkflow(k)
		if err != nil {
			a.log.Errorf("failed to get the '%s' workflow.", k)
			return nil
		}

		wfs = append(wfs, w)
	}
	return wfs
}

// Delete delete the workflow by the key.
func (a *ArgoRepository) Delete(key string) error {
	if key == "" {
		return fmt.Errorf("there is no key to delete")
	}

	ns, n, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}

	a.log.Debugf("delete '%s' workflow", key)
	err = a.ac.ArgoprojV1alpha1().Workflows(ns).Delete(n, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Logs return the slice of log which type is string and the identifier and message are separated by delimeter.
// func (a *ArgoRepository) Logs(key string) (logs []string, delim string, err error) {
//
// }
