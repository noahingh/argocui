package workflow

import (
	"fmt"
	"strings"
	"time"

	svc "github.com/hanjunlee/argocui/internal/runtime"
	tw "github.com/hanjunlee/argocui/pkg/tablewriter"
	"github.com/hanjunlee/argocui/pkg/tree"
	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	colorutil "github.com/hanjunlee/argocui/pkg/util/color"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	h "github.com/argoproj/pkg/humanize"
	"github.com/jroimartin/gocui"
	"k8s.io/apimachinery/pkg/runtime"
)

// Presentor is the presentor of mock.
type Presentor struct {
	podcolor map[string]gocui.Attribute
}

// NewPresentor create a new presentor.
func NewPresentor() *Presentor {
	return &Presentor{
		podcolor: make(map[string]gocui.Attribute),
	}
}

// PresentCore present the core view for Animal.
func (p *Presentor) PresentCore(v *gocui.View, objs []runtime.Object) error {
	w, _ := v.Size()
	nameWidth := w - 90
	if nameWidth < 45 {
		nameWidth = 45
	}

	t := tw.NewTableWriter(v)
	t.SetColumns([]string{"NAMESPACE", "NAME", "STATUS", "AGE", "DURATION"})
	t.SetColumnWidths([]int{50, nameWidth, 20, 10, 10})
	t.SetHeaderBorder(true)
	t.AppendBulk(p.convertToRows(objs))
	return t.Render()
}

// TODO: have a unit-test
func (p *Presentor) convertToRows(objs []runtime.Object) [][]string {
	rows := make([][]string, 0)

	for _, o := range objs {
		w := o.(*wf.Workflow)
		var (
			namespace = w.GetNamespace()
			name      = w.GetName()
			status    = argoutil.WorkflowStatus(w)
			age       = h.RelativeDurationShort(w.ObjectMeta.CreationTimestamp.Time, time.Now())
			duration  = h.RelativeDurationShort(w.Status.StartedAt.Time, w.Status.FinishedAt.Time)
		)
		rows = append(rows, []string{namespace, name, status, age, duration})
	}
	return rows
}

// PresentInformer display the general information and the tree of nodes like "argo get".
func (p *Presentor) PresentInformer(v *gocui.View, obj runtime.Object) error {
	w := obj.(*wf.Workflow)
	width, _ := v.Size()

	// general information
	t := tw.NewTableWriter(v)
	t.SetColumnWidths([]int{40, width - 40})
	t.AppendBulk(tree.GetInfo(w))
	t.Render()
	fmt.Fprintln(v, strings.Repeat(" ", width))

	// tree
	t = tw.NewTableWriter(v)
	t.SetColumns([]string{"STEP", "PODNAME", "DURATION", "MESSAGE"})
	t.SetColumnWidths([]int{50, 50, 15, width - 95})

	tr, err := tree.GetTreeRoot(w)
	if err != nil {
		return err
	}
	t.AppendBulk(tr)

	te, err := tree.GetTreeExit(w)
	if err != nil {
		return err
	}
	t.AppendBulk(te)

	return t.Render()
}

// PresentFollower display logs and color the display name of node.
func (p *Presentor) PresentFollower(v *gocui.View, logs []svc.Log) error {
	w, _ := v.Size()
	t := tw.NewTableWriter(v)

	t.SetColumns([]string{"NAME", "MESSAGE"})
	t.SetColumnWidths([]int{50, w - 40})
	t.SetHeaderBorder(true)
	t.AppendBulk(p.convertLogsToRows(logs))
	return t.Render()
}

func (p *Presentor) convertLogsToRows(logs []svc.Log) [][]string {
	rows := make([][]string, 0)

	for _, l := range logs {
		pc := p.nodeColor(l)
		rows = append(rows, []string{colorutil.ChangeColor(l.DisplayName+":", pc), l.Message})
	}
	return rows
}

func (p *Presentor) nodeColor(log svc.Log) gocui.Attribute {
	var (
		colorset = []gocui.Attribute{
			gocui.ColorRed,
			gocui.ColorGreen,
			gocui.ColorYellow,
			gocui.ColorBlue,
			gocui.ColorMagenta,
			gocui.ColorCyan,
			gocui.ColorBlack,
		}
	)

	// set a color of pod.
	pod := log.Pod
	color, has := p.podcolor[pod]
	if has {
		return color
	}

	mod := len(p.podcolor) % len(colorset)
	c := colorset[mod]
	p.podcolor[pod] = c
	return c
}
