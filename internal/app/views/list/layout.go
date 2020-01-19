package list

import (
	"time"

	"github.com/hanjunlee/argocui/pkg/clientset"
	"github.com/hanjunlee/argocui/pkg/view"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

const (
	viewName = "list"
)

var (
	listPeriod = 2 * time.Second
)

func layoutList(s g *gocui.Gui, x0, y0, x1, y1 int) error {
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
		v.SetCursor(0, clientset.ListTableHeadSize)

		// focus on
		g.SetCurrentView(viewName)

		// set refresh
		go view.RefreshViewPeriodic(g, v, listPeriod, func() error {
		})
	}

	return nil
}

func listWidths(width int) []int {
	var (
		ns       = 30
		wN       = width - 70
		status   = 20
		age      = 10
		duration = 10
	)

	if wN < 45 {
		wN = 45
	}
	// TODO: much more active
	return []int{ns, wN, status, age, duration}
}
