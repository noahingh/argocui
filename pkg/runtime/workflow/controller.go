package workflow

import (
	"fmt"
	"time"

	"github.com/argoproj/argo/pkg/client/clientset/versioned"
	informers "github.com/argoproj/argo/pkg/client/informers/externalversions/workflow/v1alpha1"
	log "github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	eventAdd    = "add"
	eventUpdate = "update"
	eventDelete = "delete"
	threadCount = 1
)

type eventItem struct {
	// the type of event.
	event string
	// the key of workflow.
	key string
}

type controller struct {
	clientset versioned.Interface
	informer  informers.WorkflowInformer
	storage   *storage
	workqueue workqueue.RateLimitingInterface
	log       *log.Entry
}

func newController(clientset versioned.Interface, informer informers.WorkflowInformer, storage *storage) *controller {
	c := &controller{
		clientset: clientset,
		informer:  informer,
		storage:   storage,
		workqueue: workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		log: log.WithFields(log.Fields{
			"pkg":  "repo",
			"file": "controller.go",
		}),
	}

	c.log.Debug("create a new controller and add events.")
	informer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.enqueue(obj, eventAdd)
		},
		UpdateFunc: func(old, new interface{}) {
			c.enqueue(new, eventUpdate)
		},
		DeleteFunc: func(obj interface{}) {
			c.enqueue(obj, eventDelete)
		},
	})
	return c
}

func (c *controller) enqueue(obj interface{}, event string) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		c.log.Debug("failed to enqueue the object.")
		return
	}

	c.log.WithField("key", key).Tracef("enqueue the '%s event for '%s'.", event, key)
	c.workqueue.Add(eventItem{
		key:   key,
		event: event,
	})
}

// run the controller.
func (c *controller) run(stopCh <-chan struct{}) error {
	defer c.workqueue.ShutDown()

	i := c.informer.Informer()
	if ok := cache.WaitForCacheSync(stopCh, i.HasSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	c.log.Debugf("run '%d' of workers.", threadCount)
	for j := 0; j < threadCount; j++ {
		wait.Until(c.work, time.Second, stopCh)
	}

	c.log.Debug("run the controller.")
	<-stopCh

	return nil
}

func (c *controller) waitForSynced(stopCh <-chan struct{}) error {
	i := c.informer.Informer()
	if ok := cache.WaitForCacheSync(stopCh, i.HasSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	return nil
}

func (c *controller) work() {
	for c.nextItem() {
	}
}

func (c *controller) nextItem() bool {
	obj, shutdown := c.workqueue.Get()
	if shutdown {
		return false
	}
	defer c.workqueue.Done(obj)

	e, ok := obj.(eventItem)
	if !ok {
		c.workqueue.Forget(obj)
		return true
	}

	switch e.event {
	case eventAdd:
		if c.storage.Has(e.key) {
			c.log.WithField("key", e.key).Tracef("'%s' already exist in the storage.", e.key)
			break
		}

		c.log.WithField("key", e.key).Tracef("insert '%s' into the storage.", e.key)
		c.storage.ReplaceOrInsert(e.key)

	case eventUpdate:
		c.storage.ReplaceOrInsert(e.key)

	case eventDelete:
		c.storage.Delete(e.key)
	}

	return true
}
