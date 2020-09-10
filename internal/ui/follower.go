package ui

import (
	"context"
	"sort"
	"time"

	"github.com/hanjunlee/argocui/internal/runtime"
	"github.com/hanjunlee/argocui/internal/ui/workflow"
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

type logGroup []runtime.Log

func (g logGroup) Len() int {
	return len(g)
}

func (g logGroup) Less(i, j int) bool {
	return g[i].Time.Before(g[j].Time)
}

func (g logGroup) Swap(i, j int) {
	t := g[i]
	g[i] = g[j]
	g[j] = t
	return
}

// NewFollower create a new view to follow logs of a object.
func (m *Manager) NewFollower(g *gocui.Gui, key string) error {
	w, h := g.Size()

	// set view
	v, err := g.SetView(Follower, 0, h/5+3, w-1, h-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Frame = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, headerSize)
		g.SetCurrentView(Follower)
	}

	// set keybinding
	if err := g.SetKeybinding(Follower, gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return m.ReturnFollower(g)
		}); err != nil {
		return err
	}

	// TODO: display time
	if err := g.SetKeybinding(Follower, 't', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Follower, 'k', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorUp(g, v, headerSize)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Follower, 'j', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorDown(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Follower, 'H', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorTop(g, v, headerSize)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Follower, 'L', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorBottom(g, v)
		}); err != nil {
		return err
	}

	// follow a workflow
	m.logs = make([]runtime.Log, 0)
	ctx, cancel := context.WithCancel(context.Background())
	m.cancel = cancel
	ch, err := m.svc.Logs(ctx, key)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Infof("stop to follow: %s.", key)
				return
			case log := <-ch:
				m.logs = append(m.logs, log)
			}
		}
	}()

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

		sort.Sort(logGroup(m.logs))
		p.PresentFollower(v, m.logs)
		return nil
	})

	return nil
}

// ReturnFollower switch to the viewCore.
func (m *Manager) ReturnFollower(g *gocui.Gui) error {
	defer g.SetCurrentView(viewCore)
	defer g.DeleteView(Follower)
	defer g.DeleteKeybindings(Follower)

	m.cancel()
	return nil
}
