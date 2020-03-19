package workflow

// TODO: write a unit-test for this workflow.

import (
	"bufio"
	"context"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	svc "github.com/hanjunlee/argocui/internal/runtime"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/argo/pkg/client/clientset/versioned"
	listerv1alpha1 "github.com/argoproj/argo/pkg/client/listers/workflow/v1alpha1"
	"github.com/argoproj/argo/workflow/util"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// Repo is the workflow repository it always syncronizes the workflows as the storage at background.
type Repo struct {
	ac versioned.Interface
	al listerv1alpha1.WorkflowLister
	kc kubernetes.Interface
}

// NewRepo create a new workflow repository.
func NewRepo(ac versioned.Interface, al listerv1alpha1.WorkflowLister, kc kubernetes.Interface) *Repo {
	repo := &Repo{
		ac: ac,
		al: al,
		kc: kc,
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
	sort.Sort(workflows(ws))

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

	err = r.ac.ArgoprojV1alpha1().Workflows(ns).Delete(n, &metav1.DeleteOptions{})
	if err != nil {
		return err
	}

	return nil
}

// Logs get the channel to recieve Logs from a Argo workflow.
func (r *Repo) Logs(ctx context.Context, key string) (<-chan svc.Log, error) {
	ns, n, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return nil, err
	}

	w, err := r.al.Workflows(ns).Get(n)
	if err != nil {
		return nil, err
	}

	var (
		ch = make(chan svc.Log, 100)
	)
	// get logs and send to the channel.
	err = r.logsWorkflow(ctx, ch, w)
	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (r *Repo) logsWorkflow(ctx context.Context, ch chan<- svc.Log, w *wf.Workflow) error {
	err := util.DecompressWorkflow(w)
	if err != nil {
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
		ns, n, dn := w.Namespace, n.ID, n.DisplayName

		// get logs from nodes at background.
		go func() {
			log.Tracef("log '%s' node.", n)
			err := r.logsPod(ctx, ch, ns, n, dn)
			if err != nil {
				log.Errorf("couldn't get logs from '%s' node: %s.", n, err)
				return
			}
		}()
	}

	return nil
}

func (r *Repo) logsPod(ctx context.Context, ch chan<- svc.Log, ns string, n string, dn string) error {
	const (
		mainContainerName = "main"
	)

	req := r.kc.CoreV1().Pods(ns).GetLogs(n, &corev1.PodLogOptions{
		Container:  mainContainerName,
		Follow:     true,
		Timestamps: true, // add an RFC3339 or RFC3339Nano timestamp at the beginning
	})
	// TODO: mocking for unit-test.
	if isUnitTest := reflect.DeepEqual(req, &restclient.Request{}); isUnitTest {
		return nil
	}

	s, err := req.Stream()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(s)
	for {
		select {
		case <-ctx.Done():
			log.Trace("the context is closed.")
			return nil

		default:
			if !scanner.Scan() {
				log.Trace("finished to logs.")
				return nil
			}

			line := scanner.Text()
			t, m := splitTimeAndMessage(line)

			time, err := time.Parse(time.RFC3339, t)
			if err != nil {
				log.Warnf("can't parse the timestamp: %s", err)
				continue
			}

			ch <- svc.Log{
				DisplayName: dn,
				Pod:         n,
				Message:     m,
				Time:        time,
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
