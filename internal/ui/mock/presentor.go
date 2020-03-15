package mock

import (
	svc "github.com/hanjunlee/argocui/pkg/runtime"
	mockrepo "github.com/hanjunlee/argocui/pkg/runtime/mock"
	tw "github.com/hanjunlee/argocui/pkg/tablewriter"
	err "github.com/hanjunlee/argocui/pkg/util/error"

	"k8s.io/apimachinery/pkg/runtime"
	"github.com/jroimartin/gocui"
)

// Presentor is the presentor of mock.
type Presentor struct {}

// NewPresentor create a new presentor.
func NewPresentor() *Presentor {
	return &Presentor{}
}

// PresentCore present the core view for Animal.
func (p *Presentor) PresentCore(v *gocui.View, objs []runtime.Object) error {
	const (
		namespaceWidth = 50
	)
	w, _ := v.Size()

	t := tw.NewTableWriter(v)
	t.SetColumns([]string{"NAMESPACE", "NAME"})
	t.SetColumnWidths([]int{namespaceWidth, w - namespaceWidth})
	t.SetHeaderBorder(true)
	t.AppendBulk(p.convertToRows(objs))
	return t.Render()
}

func (p *Presentor) convertToRows(objs []runtime.Object) [][]string {
	rows := make([][]string, 0)

	for _, o := range objs {
		a := o.(*mockrepo.Animal)
		rows = append(rows, []string{a.GetNamespace(), a.GetName()})
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
