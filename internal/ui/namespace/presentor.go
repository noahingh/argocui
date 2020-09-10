package namespace

import (
	"fmt"
	"text/tabwriter"

	svc "github.com/hanjunlee/argocui/internal/runtime"
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
	width, _ := v.Size()

	w := tabwriter.NewWriter(v, width, 1, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "NAMESPACE\t")
	items := p.convertToRows(objs)
	for _, i := range items {
		fmt.Fprintf(w, "%s\t\n", i[0])
	}

	return w.Flush()
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
