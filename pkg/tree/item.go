package tree

import (
	"fmt"
	"strings"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/hanjunlee/tree"
)

const (
	tabRetrySingleChild     = "──"
	tabStepGroupSingleChild = "──"
	tabStepGroupMultiChild  = "·─"
)

var (
	icons map[wf.NodePhase]string
)

func init() {
	icons = map[wf.NodePhase]string{
		wf.NodePending:   "◷",
		wf.NodeRunning:   "●",
		wf.NodeSucceeded: "✔",
		wf.NodeSkipped:   "○",
		wf.NodeFailed:    "✖",
		wf.NodeError:     "⚠",
	}
}

func items(nodes map[string]wf.NodeStatus) map[string]item {
	var (
		ret = make(map[string]item)
	)

	for id, ns := range nodes {
		ret[id] = item(ns)
	}
	return ret
}

// item is the implement of the Item interface of tree.
type item wf.NodeStatus

// string return the status of node as string.
func (i item) String() string {
	if i.Type == wf.NodeTypeStepGroup {
		return fmt.Sprintf("%s", icons[i.Phase])
	}

	var (
		ret         string
		prefix      string
		displayName = i.DisplayName
	)

	// preprocess for some edge cases, retry and step.
	if strings.Index(displayName, tabRetrySingleChild) == 0 {
		prefix = displayName[:len(tabRetrySingleChild)]
		displayName = displayName[len(tabRetrySingleChild):]

	}
	if strings.Index(displayName, tabStepGroupMultiChild) == 0 {
		prefix = displayName[:len(tabStepGroupMultiChild)]
		displayName = displayName[len(tabStepGroupMultiChild):]
	}

	ret = fmt.Sprintf("%s%s %s", prefix, icons[i.Phase], displayName)

	// append the suffix for the template.
	if i.TemplateRef != nil {
		ret = fmt.Sprintf("%s (%s/%s)", ret, i.TemplateRef.Name, i.TemplateRef.Template)
	} else if i.TemplateName != "" {
		ret = fmt.Sprintf("%s (%s)", ret, i.TemplateName)
	}
	return ret
}

// Less compare the start time first, and if the start times are equal it compare the display name.
func (i item) Less(comp tree.Item) bool {
	c := comp.(item)

	if equal := i.StartedAt.Equal(&c.StartedAt); !equal {
		return i.StartedAt.Before(&c.StartedAt)
	}

	if i.DisplayName != c.DisplayName {
		return i.DisplayName < c.DisplayName
	}
	return i.ID < c.ID
}
