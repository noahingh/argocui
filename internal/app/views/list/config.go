package list

import (
	"fmt"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/sirupsen/logrus"
)

// Configuration of view.
type config struct {
	namespace   string
	namePattern string
	cache       []*wf.Workflow
}

func (c config) pattern() string {
	return fmt.Sprintf("%s/*%s*", c.namespace, c.namePattern)
}

var (
	// the configuration of the list view.
	conf config
	log = logrus.WithFields(logrus.Fields{
		"pkg": "list",
	})
)

func init() {
	conf = config{
		namespace:   "*",
		namePattern: "*",
		cache:       []*wf.Workflow{},
	}
}
