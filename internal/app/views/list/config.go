package list

import (
	"fmt"
	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
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
