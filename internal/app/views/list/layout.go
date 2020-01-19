package list

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	h "github.com/argoproj/pkg/humanize"
	"github.com/hanjunlee/argocui/pkg/argo"
	tw "github.com/hanjunlee/argocui/pkg/table/tablewriter"
	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	"github.com/hanjunlee/argocui/pkg/view"
	"github.com/jroimartin/gocui"
)

const (
	viewName           = "list"
	upperBoundOfCursor = 2
)

var (
	// the configuration of the list view.
	conf config
	log = logrus.WithFields(logrus.Fields{
		"pkg": "list",
		"file": "layout.go",
	})
)

func init() {
	conf = config{
		namespace:   "*",
		namePattern: "*",
		cache:       []*wf.Workflow{},
	}
}

// Layout lay out the list view.
func Layout(s argo.UseCase, g *gocui.Gui, x0, y0, x1, y1 int) error {
	var (
		period = 1 * time.Second
	)

	v, err := g.SetView(viewName, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// settings of list view
		v.Highlight = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.Frame = false
		v.SetCursor(0, upperBoundOfCursor)

		// focus on
		g.SetCurrentView(viewName)

		// set refresh
		go view.RefreshViewPeriodic(g, v, period, func() error {
			v.Clear()

			wfs := s.Search(conf.pattern())
			conf.cache = wfs

			err := render(v, toRows(wfs))
			if err != nil {
				fmt.Fprintf(v, "error occurs: %s", err)
			}
			return nil
		})
	}

	return nil
}

// render present workflows as table-like format.
func render(v *gocui.View, datas [][]string) error {
	var width, nameWidth int

	// set widths for each column.
	width, _ = v.Size()
	nameWidth = width - 70
	if nameWidth < 45 {
		nameWidth = 45
	}

	t := tw.NewTableWriter(v)

	t.SetColumns([]string{"NAMESPACE", "NAME", "STATUS", "AGE", "DURATION"})
	t.SetColumnWidths([]int{30, nameWidth, 20, 10, 10})
	t.SetHeaderBorder(true)
	t.AppendBulk(datas)
	return t.Render()
}

func toRows(wfs []*wf.Workflow) [][]string {
	var (
		rows = make([][]string, 0)
	)

	for _, w := range wfs {
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
