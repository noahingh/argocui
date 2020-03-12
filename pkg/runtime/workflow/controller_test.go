package workflow

import (
	"reflect"
	"testing"
	"time"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/sirupsen/logrus"
	"github.com/argoproj/argo/pkg/client/clientset/versioned/fake"
	informers "github.com/argoproj/argo/pkg/client/informers/externalversions"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"

	"k8s.io/client-go/util/workqueue"
)

func newTestController(items []eventItem) *controller {
	var (
		wfs  = []*wf.Workflow{}
		objs = []runtime.Object{}
	)

	// create new workflows
	now := meta.Now()
	for delay, i := range items {
		if i.event == eventDelete {
			continue
		}

		namespace, name, _ := cache.SplitMetaNamespaceKey(i.key)

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

	// create a new queue
	q := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	for _, i := range items {
		q.Add(i)
	}

	return &controller{
		clientset: c,
		informer:  i,
		storage: newStorage(i),
		workqueue: q,
		log: logrus.NewEntry(logrus.New()),
	}
}

func Test_controller_nextItem(t *testing.T) {
	c := newTestController([]eventItem{
		eventItem{
			event: eventAdd,
			key:   "argo/first",
		},
		eventItem{
			event: eventAdd,
			key:   "argo/second",
		},
		eventItem{
			event: eventAdd,
			key:   "argo/third",
		},
	})

	tests := []struct {
		name   string
		controller *controller
		want   []string
	}{
		{
			name:   "add items",
			controller: c,
			want: []string{
				"argo/third",
				"argo/second",
				"argo/first",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.controller
			t.Log("start to work.")
			for cnt := c.workqueue.Len(); cnt > 0 ; cnt-- {
				c.nextItem()
			}
			if got := c.storage.list; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("c.storage.list = %v, want %v", got, tt.want)
			}
		})
	}
}
