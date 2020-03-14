package tree

import (
	"fmt"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/argoproj/pkg/humanize"
	"github.com/hanjunlee/argocui/pkg/util/argo"
	"github.com/hanjunlee/tree"
)

func init() {
	tree.TabChild = "├─"
	tree.TabGrandChild = "│ "
	tree.TabLastChild = "└─"
	tree.TabGrandLastChild = "  "
}

// GetInfo get the general information of a Argo workflow, such as name, namespace, service account and so on.
// It return rows consist of two columns, key and value.
func GetInfo(w *wf.Workflow) [][]string {
	var (
		ret = make([][]string, 0)
	)

	setKeyVal := func(key, val string) {
		ret = append(ret, []string{key, val})
	}

	setKeyVal("Name:", w.ObjectMeta.Name)
	setKeyVal("Namespace:", w.ObjectMeta.Namespace)

	sa := w.Spec.ServiceAccountName
	if sa == "" {
		sa = "default"
	}
	setKeyVal("ServiceAccount:", sa)
	setKeyVal("Status:", argo.WorkflowStatus(w))
	if w.Status.Message != "" {
		setKeyVal("Message:", w.Status.Message)
	}

	setKeyVal("Created:", humanize.Timestamp(w.ObjectMeta.CreationTimestamp.Time))
	if !w.Status.StartedAt.IsZero() {
		setKeyVal("Started:", humanize.Timestamp(w.Status.StartedAt.Time))
	}
	if !w.Status.FinishedAt.IsZero() {
		setKeyVal("Finished:", humanize.Timestamp(w.Status.FinishedAt.Time))
	}
	if !w.Status.StartedAt.IsZero() {
		setKeyVal("Duration:", humanize.RelativeDuration(w.Status.StartedAt.Time, w.Status.FinishedAt.Time))
	}

	if len(w.Spec.Arguments.Parameters) > 0 {
		setKeyVal("Parameters:", "")
		for _, p := range w.Spec.Arguments.Parameters {
			if p.Value == nil {
				continue
			}
			setKeyVal(fmt.Sprintf("  %s:", p.Name), *p.Value)
		}
	}
	if w.Status.Outputs != nil {
		if len(w.Status.Outputs.Parameters) > 0 {
			setKeyVal("Output Parameters:", "")
			for _, p := range w.Status.Outputs.Parameters {
				setKeyVal(fmt.Sprintf("  %s:", p.Name), *p.Value)
			}
		}
		if len(w.Status.Outputs.Artifacts) > 0 {
			setKeyVal("Output Artifacts:", "")
			for _, art := range w.Status.Outputs.Artifacts {
				if art.S3 != nil {
					setKeyVal(fmt.Sprintf("  %s:", art.Name), art.S3.String())
				} else if art.Artifactory != nil {
					setKeyVal(fmt.Sprintf("  %s:", art.Name), art.Artifactory.String())
				}
			}
		}
	}
	return ret
}

// GetTreeRoot return the tree of a root node in workflow, and it consist of step, podname, duration and message.
func GetTreeRoot(w *wf.Workflow) ([][]string, error) {
	rootName := w.GetName()
	return GetTreeNode(w, rootName)
}

// GetTreeExit return the tree of a on-exit node in workflow, and it consist of step, podname, duration and message.
func GetTreeExit(w *wf.Workflow) ([][]string, error) {
	name := w.GetName()
	onExitName := fmt.Sprintf("%s.onExit", name)
	return GetTreeNode(w, onExitName)
}

// GetTreeNode return the tree of a node in workflow, consist of step, podname, duration and message.
func GetTreeNode(w *wf.Workflow, nodeName string) ([][]string, error) {
	var (
		ret   [][]string
		nodes = w.Status.Nodes
	)
	if _, ok := nodes[nodeName]; nodes == nil || !ok {
		return [][]string{}, nil
	}

	// make the tree for steps.
	root := item(nodes[nodeName])
	t := tree.NewTree(root)

	initTree(t, items(nodes), root)

	steps := t.Render()
	items := t.RenderedItems()
	for idx := 0; idx < len(steps); idx++ {
		s, i := steps[idx], items[idx].(item)

		if i.Type == wf.NodeTypePod {
			row := []string{s, i.ID, humanize.RelativeDurationShort(i.StartedAt.Time, i.FinishedAt.Time), i.Message}
			ret = append(ret, row)
		} else {
			row := []string{s, "", "", i.Message}
			ret = append(ret, row)
		}
	}

	return ret, nil
}

// initTree make the tree of Argo workflow.
func initTree(t *tree.Tree, items map[string]item, root item) error {
	for _, id := range root.Children {
		child := items[id]

		err := move(t, items, child, root)
		if err != nil {
			return err
		}
	}
	return nil
}

func move(t *tree.Tree, items map[string]item, child item, parent item) error {
	// if the parent is a execution node the child have to be one of children of the boundary node
	// except a force move.
	if isExecutionNode(parent) {
		item := items[child.BoundaryID]
		parent = item
	}

	// some edge cases, "retry" and "step group", are skipped, and are replaced by the first grand child.
	if child.Type == wf.NodeTypeRetry && len(child.Children) == 1 {
		firstGrandChild := items[child.Children[0]]
		firstGrandChild.DisplayName = tabRetrySingleChild + firstGrandChild.DisplayName
		return move(t, items, firstGrandChild, parent)
	}

	if child.Type == wf.NodeTypeStepGroup && len(child.Children) == 1 {
		firstGrandChild := items[child.Children[0]]
		firstGrandChild.DisplayName = tabStepGroupSingleChild + firstGrandChild.DisplayName
		return move(t, items, firstGrandChild, parent)
	}

	err := t.Move(child, parent)
	if err != nil {
		return err
	}

	for _, id := range child.Children {
		grandChild := items[id]
		err = move(t, items, grandChild, child)
	}
	if err != nil {
		return err
	}

	return err
}

func isExecutionNode(i item) bool {
	if t := i.Type; t == wf.NodeTypePod || t == wf.NodeTypeSkipped || t == wf.NodeTypeSuspend {
		return true
	}
	return false
}
