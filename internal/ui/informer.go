package ui

import (
	"time"

	viewutil "github.com/hanjunlee/argocui/pkg/util/view"
	"github.com/hanjunlee/argocui/internal/ui/workflow"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

// NewInformer create a new view to display the information of a object.
func (m *Manager) NewInformer(g *gocui.Gui, key string) error {
	w, h := g.Size()

	// set view
	v, err := g.SetView(Informer, 0, h/5+3, w-1, h-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Frame = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, 0)
		g.SetCurrentView(Informer)
	}

	// set keybinding
	if err := g.SetKeybinding(Informer, gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return m.ReturnInformer(g)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Informer, 'k', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorUp(g, v, 0)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Informer, 'j', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorDown(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Informer, 'H', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorTop(g, v, 0)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Informer, 'L', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorBottom(g, v)
		}); err != nil {
		return err
	}

	// refresh the view
	go viewutil.RefreshViewPeriodic(g, v, 1*time.Second, func() error {
		v.Clear()
		o, err := m.svc.Get(key)
		if err != nil {
			log.Errorf("failed to get the object: %s", err)
			return nil
		}

		var p Presentor
		gvk, _, _ := objectKind(o)
		switch gvk.Kind {
		case "Workflow":
			p = workflow.NewPresentor()
		default:
			return nil
		}
		
		p.PresentInformer(v, o)
		return nil
	})

	return nil
}

// ReturnInformer switch to the Core.
func (m *Manager) ReturnInformer(g *gocui.Gui) error {
	defer g.SetCurrentView(Core)
	defer g.DeleteView(Informer)
	defer g.DeleteKeybindings(Informer)
	return nil
}
