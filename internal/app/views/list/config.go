package list

import (
	"fmt"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/sirupsen/logrus"
)

// Config is the configuration of list.
type Config struct {
	namespace   string
	namePattern string
	cache       []*wf.Workflow
	log *logrus.Entry
}

// NewConfig create a new config of list view.
func NewConfig() *Config {
	return &Config{
		namespace:   "*",
		namePattern: "",
		cache:       []*wf.Workflow{},
		log: logrus.WithFields(logrus.Fields{
			"pkg": "list",
		}),
	}
}

func (c *Config) pattern() string {
	return fmt.Sprintf("%s/*%s*", c.namespace, c.namePattern)
}
