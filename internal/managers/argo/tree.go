package argo

import (
	"fmt"
	"strings"
	"time"

	"github.com/hanjunlee/argocui/pkg/argo"
	tw "github.com/hanjunlee/argocui/pkg/table/tablewriter"
	"github.com/hanjunlee/argocui/pkg/table/tree"
	"github.com/hanjunlee/argocui/pkg/util/view"

	"github.com/asaskevich/EventBus"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

type treeManager struct {
	// the key of workflow.
	// if the key is an empty string it is deactivate state.
	key string

	uc  argo.UseCase
	bus EventBus.Bus

	log *log.Entry
}

func newTreeManager(uc argo.UseCase, bus EventBus.Bus) *treeManager {
	return &treeManager{
		uc:  uc,
		bus: bus,
		log: log.WithFields(log.Fields{
			"pkg":  "argo-manager",
			"file": "tree.go",
		}),
	}
}

func (t *treeManager) isActive() bool {
	if t.key == "" {
		return false
	}
	return true
}

const (
	treeViewName   = "tree"
	treeUpperBound = 0
)

// lay out the tree view.
func (t *treeManager) layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var (
		period = 1 * time.Second
	)
	v, err := g.SetView(treeViewName, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Frame = false
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack

		v.SetCursor(0, treeUpperBound)

		go view.RefreshViewPeriodic(g, v, period, func() error {
			v.Clear()

			if !t.isActive() {
				t.log.Debug("not active")
				return nil
			}

			w := t.uc.Get(t.key)
			if w == nil {
				t.log.Warnf("there isn't the workflow: %s.", t.key)
				return nil
			}

			// the information of a workflow.
			info := tree.GetInfo(w)
			if err := t.renderInfo(v, info); err != nil {
				t.log.Error("failed to render info.")
				return nil
			}

			width, _ := v.Size()
			fmt.Fprintln(v, strings.Repeat(" ", width))

			// the tree of a workflow.
			root, _ := tree.GetTreeRoot(w)
			if err = t.renderRoot(v, root); err != nil {
				t.log.Error("failed to render root.")
				return nil
			}

			exit, _ := tree.GetTreeExit(w)
			if err = t.renderExit(v, exit); err != nil {
				t.log.Error("failed to render exit.")
				return nil
			}

			return nil
		})

		t.keybinding(g)
		t.subscribe(g)
	}
	return nil
}

func (t *treeManager) renderInfo(v *gocui.View, datas [][]string) error {
	var width int

	// set widths for each column.
	width, _ = v.Size()

	w := tw.NewTableWriter(v)

	w.SetColumnWidths([]int{40, width - 40})
	w.AppendBulk(datas)
	return w.Render()
}

func (t *treeManager) renderRoot(v *gocui.View, datas [][]string) error {
	width, _ := v.Size()

	w := tw.NewTableWriter(v)

	w.SetColumns([]string{"STEP", "PODNAME", "DURATION", "MESSAGE"})
	w.SetColumnWidths([]int{40, 40, 15, width - 95})
	w.AppendBulk(datas)
	return w.Render()
}

func (t *treeManager) renderExit(v *gocui.View, datas [][]string) error {
	width, _ := v.Size()

	w := tw.NewTableWriter(v)

	w.SetColumnWidths([]int{40, 40, 15, width - 95})
	w.AppendBulk(datas)
	return w.Render()
}

func (t *treeManager) keybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding(treeViewName, 'k', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return view.MoveCursorUp(g, v, treeUpperBound)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeViewName, 'j', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return view.MoveCursorDown(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeViewName, 'H', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return view.MoveCursorTop(g, v, treeUpperBound)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeViewName, 'L', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return view.MoveCursorBottom(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeViewName, gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			t.log.Info("deactivate the workflow.")
			t.key = ""
			t.bus.Publish(eventCollectionSetView)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

const (
	eventTreeSetView     = "tree:set-view"
	eventTreeGetWorkflow = "tree:get-workflow"
)

func (t *treeManager) subscribe(g *gocui.Gui) error {
	if err := t.bus.Subscribe(eventTreeSetView, func() {
		t.log.Info("set the tree current view.")
		g.SetViewOnTop(treeViewName)
		g.SetCurrentView(treeViewName)
	}); err != nil {
		return err
	}

	if err := t.bus.Subscribe(eventTreeGetWorkflow, func(key string) {
		t.log.Infof("get the workflow: %s.", key)
		t.key = key

		t.log.Infof("init the cursor.")
		v, _ := g.View(treeViewName)
		v.SetCursor(0, treeUpperBound)
		v.SetOrigin(0, 0)
	}); err != nil {
		return err
	}

	return nil
}
