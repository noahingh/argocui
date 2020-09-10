package workflow

import (
	"fmt"
	"strings"

	// "strings"
	"text/tabwriter"
	"time"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	h "github.com/argoproj/pkg/humanize"
	svc "github.com/hanjunlee/argocui/internal/runtime"
	tw "github.com/hanjunlee/argocui/pkg/tablewriter"
	"github.com/hanjunlee/argocui/pkg/tree"
	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	colorutil "github.com/hanjunlee/argocui/pkg/util/color"
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
	width, _ := v.Size()

	w := tabwriter.NewWriter(v, width/5, 1, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "NAMESPACE\tNAME\tSTATUS\tAGE\tDURATION\t")
	items := p.convertToRows(objs)
	for _, i := range items {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", i[0], i[1], i[2], i[3], i[4])
	}

	return w.Flush()
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

	// print the general information
	tabw := tabwriter.NewWriter(v, width/4+1, 1, 1, ' ', tabwriter.TabIndent)

	infos := tree.GetInfo(w)
	for _, i := range infos {
		key, val := i[0], i[1]
		fmt.Fprintf(tabw, "%s\t%s\t\t\t\n", key, val)
	}

	if err := tabw.Flush(); err != nil {
		return err
	}

	fmt.Fprintln(tabw, strings.Repeat(" ", width))

	// print the workflow tree
	fmt.Fprintln(tabw, "STEP\tPODNAME\tDURATION\tMESSAGE\t")
	tr, err := tree.GetTreeRoot(w)
	if err != nil {
		return err
	}

	for _, i := range tr {
		fmt.Fprintf(tabw, "%s\t%s\t%s\t%s\t\n", i[0], i[1], i[2], i[3])
	}

	return tabw.Flush()
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
