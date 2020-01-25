package argo

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hanjunlee/argocui/pkg/argo"
	"github.com/hanjunlee/argocui/pkg/util/view"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
	"github.com/asaskevich/EventBus"
)

const (
	followerViewName   = "follower"
	followerUpperBound = 2
)

type followerManager struct {
	// controller goroutine which follow pods.
	ctx    context.Context
	cancel context.CancelFunc

	// store logs which comes from pods.
	key      string
	logs     logGroup
	podColor map[string]gocui.Attribute

	uc argo.UseCase
	bus EventBus.Bus

	logger *log.Entry
}

func newFollowerManager(uc argo.UseCase, bus EventBus.Bus) *followerManager {
	return &followerManager{
		uc: uc,
		bus: bus,
		logger: log.WithFields(log.Fields{
			"pkg": "argo manager",
			"file": "follower.go",
		}),
	}
}

func (f *followerManager) layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var (
		period = 1 * time.Second
	)
	v, err := g.SetView(followerViewName, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Frame = false
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack

		v.SetCursor(0, followerUpperBound)

		go view.RefreshViewPeriodic(g, v, period, func() error {
			v.Clear()
			fmt.Fprintln(v, "follower")

			return nil
		})

		f.keybindingLogs(g)
		f.subscribe(g)
	}
	return nil

}

// keybinding of the follower.
func (f *followerManager) keybindingLogs(g *gocui.Gui) error {
	if err := g.SetKeybinding(followerViewName, 'k', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return view.MoveCursorUp(g, v, followerUpperBound)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(followerViewName, 'j', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return view.MoveCursorDown(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(followerViewName, 'H', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return view.MoveCursorTop(g, v, followerUpperBound)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(followerViewName, 'L', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return view.MoveCursorBottom(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(followerViewName, gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			f.bus.Publish(eventCollectionSetView)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

// subscribe events of follower.
const (
	eventFollowerSetView = "follower:set-view"
)

func (f *followerManager) subscribe(g *gocui.Gui) error {
	if err := f.bus.Subscribe(eventFollowerSetView, func() {
		f.logger.Info("set the follower current view.")
		g.SetViewOnTop(followerViewName)
		g.SetCurrentView(followerViewName)
	}); err != nil {
		return err
	}
	return nil
}

func (f *followerManager) appendLog(log argo.Log) {
	var (
		colorset = []gocui.Attribute{
			gocui.ColorDefault,
			gocui.ColorBlack,
			gocui.ColorRed,
			gocui.ColorGreen,
			gocui.ColorYellow,
			gocui.ColorBlue,
			gocui.ColorMagenta,
			gocui.ColorCyan,
		}
	)
	// set a color of pod.
	pod := log.Pod
	if _, has := f.podColor[log.Pod]; !has {
		mod := len(f.podColor) % len(colorset)
		f.podColor[pod] = colorset[mod]
	}

	f.logs = append(f.logs, log)
	sort.Sort(f.logs)
	return
}
