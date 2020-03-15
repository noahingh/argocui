package namespace

import (
	svc "github.com/hanjunlee/argocui/internal/runtime"
	tw "github.com/hanjunlee/argocui/pkg/tablewriter"
	err "github.com/hanjunlee/argocui/pkg/util/error"

	"github.com/jroimartin/gocui"
	corev1 "k8s.io/api/core/v1"
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

	t := tw.NewTableWriter(v)
	t.SetColumns([]string{"NAMESPACE"})
	t.SetColumnWidths([]int{w})
	t.SetHeaderBorder(true)
	t.AppendBulk(p.convertToRows(objs))
	return t.Render()
}

func (p *Presentor) convertToRows(objs []runtime.Object) [][]string {
	rows := make([][]string, 0)

	for _, o := range objs {
		n := o.(*corev1.Namespace)
		rows = append(rows, []string{n.GetName()})
	}
	return rows
}

// PresentInformer is not implemented.
func (p *Presentor) PresentInformer(v *gocui.View, objs runtime.Object) error {
	return err.ErrNotImplement
}

// PresentFollower is not implemented.
func (p *Presentor) PresentFollower(v *gocui.View, logs []svc.Log) error {
	return err.ErrNotImplement
}
