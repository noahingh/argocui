package argo

import (
	"time"
	"strings"

	"github.com/hanjunlee/argocui/pkg/kube"
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"

	"github.com/asaskevich/EventBus"
	tw "github.com/hanjunlee/argocui/pkg/table/tablewriter"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

const (
	namespaceViewName   = "namespace"
	namespaceUpperBound = 2
)

type namespaceManager struct {
	ku  kube.UseCase
	bus EventBus.Bus

	log *log.Entry
}

func newNamespaceManager(ku kube.UseCase, bus EventBus.Bus) *namespaceManager {
	return &namespaceManager{
		ku:  ku,
		bus: bus,
		log: log.WithFields(log.Fields{
			"pkg":  "argo-manager",
			"file": "namespace.go",
		}),
	}
}

// lay out the namespace view.
func (nm *namespaceManager) layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(namespaceViewName, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// settings of list view
		v.Highlight = true
		v.Frame = false
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, namespaceUpperBound)

		go viewutil.RefreshViewPeriodic(g, v, 3*time.Second, func() error {
			v.Clear()

			var (
				namespaces []string
			)
			if namespaces, err = nm.ku.GetNamespaces(); err != nil {
				nm.log.Errorf("failed to get namespaces: %s", err)
				return nil
			}

			nm.render(v, nm.toRows(namespaces))
			return nil
		})

		nm.keybinding(g)
		nm.subscribe(g)
	}

	return nil
}

func (nm *namespaceManager) render(v *gocui.View, datas [][]string) error {
	width, _ := v.Size()

	t := tw.NewTableWriter(v)

	t.SetColumns([]string{"NAMESPACE"})
	t.SetColumnWidths([]int{width})
	t.SetHeaderBorder(true)
	t.AppendBulk(datas)
	return t.Render()
}

func (nm *namespaceManager) toRows(namespaces []string) [][]string {
	ret := make([][]string, 0)
	for _, ns := range namespaces {
		ret = append(ret, []string{ns})
	}

	return ret
}

// keybindings of the namespace view.
func (nm *namespaceManager) keybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding(namespaceViewName, 'k', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorUp(g, v, namespaceUpperBound)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(namespaceViewName, 'j', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorDown(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(namespaceViewName, 'H', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorTop(g, v, namespaceUpperBound)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(namespaceViewName, 'L', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorBottom(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(namespaceViewName, gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, cy := v.Cursor() 

			ns, err := v.Line(cy)
			if err != nil {
				nm.log.Errorf("failed to get namespace.")
				return nil
			}
			ns = strings.TrimSpace(ns)

			nm.bus.Publish(eventCollectionSetNamespace, ns)
			nm.bus.Publish(eventCollectionSetView)
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(namespaceViewName, gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			nm.log.Info("back to the collection view.")
			nm.bus.Publish(eventCollectionSetView)
			return nil
		}); err != nil {
		return err
	}

	return nil
}

// subscribes events of the namespace view.
const (
	eventNamespaceSetView = "namespace:set-view"
)

func (nm *namespaceManager) subscribe(g *gocui.Gui) error {
	if err := nm.bus.Subscribe(eventNamespaceSetView, func() {
		nm.log.Info("set the current view namespace.")
		g.SetViewOnTop(namespaceViewName)
		g.SetCurrentView(namespaceViewName)
	}); err != nil {
		return err
	}
	return nil
}
