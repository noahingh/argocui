package repo

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"time"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/pkg/client/clientset/versioned"
	informers "github.com/argoproj/argo/pkg/client/informers/externalversions/workflow/v1alpha1"
	"github.com/argoproj/argo/workflow/util"
	"github.com/hanjunlee/argocui/pkg/argo"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// ArgoRepository is the repository of Argo, it always syncronizes the workflows into storage at background.
type ArgoRepository struct {
	c   *controller
	s   *storage
	ac  versioned.Interface
	kc  kubernetes.Interface
	log *log.Entry
}

// NewArgoRepository create a new Argo repository.
func NewArgoRepository(
	argoClientset versioned.Interface, argoInformer informers.WorkflowInformer, kubeClientset kubernetes.Interface) *ArgoRepository {
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

// Logs get the channel to recieve Logs from a Argo workflow.
func (a *ArgoRepository) Logs(ctx context.Context, key string) (<-chan argo.Log, error) {
	w, err := a.s.GetWorkflow(key)
	if err != nil {
		return nil, err
	}

	var (
		ch = make(chan argo.Log, 100)
	)
	// get logs and send to the channel.
	err = a.logsWorkflow(ctx, ch, w)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (a *ArgoRepository) logsWorkflow(ctx context.Context, ch chan<- argo.Log, w *wf.Workflow) error {
	err := util.DecompressWorkflow(w)
	if err != nil {
		a.log.Error(err)
		return err
	}

	// node is the unit of the executed step.
	var nodes []wf.NodeStatus
	for _, n := range w.Status.Nodes {
		if n.Type == wf.NodeTypePod && n.Phase != wf.NodeError {
			nodes = append(nodes, n)
		}
	}

	for _, n := range nodes {
		ns, n := w.Namespace, n.ID

		// get logs from nodes at background.
		go func() {
			a.log.Tracef("log '%s' node.", n)
			err := a.logsPod(ctx, ch, ns, n)
			if err != nil {
				a.log.Errorf("couldn't get logs from '%s' node: %s.", n, err)
				return
			}
		}()
	}

	return nil
}

func (a *ArgoRepository) logsPod(ctx context.Context, ch chan<- argo.Log, ns string, n string) error {
	const (
		mainContainerName = "main"
	)
	var (
		key = ns + "/" + n
	)

	s, err := a.kc.CoreV1().Pods(ns).GetLogs(n, &corev1.PodLogOptions{
		Container:  mainContainerName,
		Follow:     true,
		Timestamps: true, // add an RFC3339 or RFC3339Nano timestamp at the beginning
	}).Stream()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(s)
	for {
		select {
		case <-ctx.Done():
			a.log.WithField("key", key).Trace("the context is closed.")
			return nil

		default:
			if !scanner.Scan() {
				a.log.WithField("key", key).Trace("finished to logs.")
				return nil
			}

			line := scanner.Text()
			t, m := splitTimeAndMessage(line)

			time, err := time.Parse(time.RFC3339, t)
			if err != nil {
				a.log.WithField("key", key).Warnf("can't parse the timestamp: %s", err)
				continue
			}

			ch <- argo.Log{
				Pod:     n,
				Message: m,
				Time:    time,
			}
		}
	}
}

// splitTimeAndMessage split the log from Kubernetes into time and message.
func splitTimeAndMessage(l string) (string, string) {
	parts := strings.SplitN(l, " ", 2)
	return parts[0], parts[1]
}

func getNodeDisplayName(n wf.NodeStatus) string {
	dn := n.DisplayName
	if dn == "" {
		dn = n.Name
	}
	return dn
}
