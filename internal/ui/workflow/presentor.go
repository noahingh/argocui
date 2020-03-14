package workflow

import (
	"time"

	tw "github.com/hanjunlee/argocui/pkg/tablewriter"
	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	h "github.com/argoproj/pkg/humanize"
	"github.com/jroimartin/gocui"
	"k8s.io/apimachinery/pkg/runtime"
)

// Presentor is the presentor of mock.
type Presentor struct{}

// NewPresentor create a new presentor.
func NewPresentor() *Presentor {
	return &Presentor{}
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
