package argo

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/hanjunlee/argocui/pkg/argo"
	"github.com/hanjunlee/argocui/pkg/util/color"
	tw "github.com/hanjunlee/argocui/pkg/table/tablewriter"
	"github.com/hanjunlee/argocui/pkg/util/view"

	"github.com/asaskevich/EventBus"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
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
	// if the key is the empty string the state is unfollow.
	key      string
	logs     logGroup
	podColor map[string]gocui.Attribute

	uc  argo.UseCase
	bus EventBus.Bus

	logger *log.Entry
}

func newFollowerManager(uc argo.UseCase, bus EventBus.Bus) *followerManager {
	return &followerManager{
		uc:  uc,
		bus: bus,
		logger: log.WithFields(log.Fields{
			"pkg":  "argo manager",
			"file": "follower.go",
		}),
	}
}

func (f *followerManager) isFollowing() bool {
	if f.key != "" {
		return true
	}
	return false
}

// follow logs of the workflow until follower cancel it.
func (f *followerManager) follow(key string) {
	ctx, cancel := context.WithCancel(context.Background())
	f.ctx, f.cancel = ctx, cancel

	f.key = key
	f.logs = logGroup{}
	f.podColor = map[string]gocui.Attribute{}

	// follow
	ch, err := f.uc.Logs(ctx, key)
	if err != nil {
		f.logger.Errorf("failed to logs: %s.", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				f.logger.Infof("stop to follow: %s.", key)
				return
			case log := <-ch:
				f.appendLog(log)
			}
		}
	}()
}

func (f *followerManager) appendLog(log argo.Log) {
	var (
		colorset = []gocui.Attribute{
			gocui.ColorRed,
			gocui.ColorGreen,
			gocui.ColorYellow,
			gocui.ColorBlue,
			gocui.ColorMagenta,
			gocui.ColorCyan,
			gocui.ColorBlack,
		}
	)
	// set a color of pod.
	pod := log.Pod
	if _, has := f.podColor[pod]; !has {
		mod := len(f.podColor) % len(colorset)
		f.podColor[pod] = colorset[mod]
	}

	f.logs = append(f.logs, log)
	sort.Sort(f.logs)
	return
}

func (f *followerManager) unfollow() {
	f.logger.Debug("cancel the context.")
	f.ctx = nil
	f.cancel()

	f.key = ""
	f.logs = nil
	f.podColor = nil
}

// lay out the follower.
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

			if !f.isFollowing() {
				return nil
			}

			d := f.datas()
			if err := f.render(v, d); err != nil {
				fmt.Fprintln(v, "fail to render, see the log.")
				f.logger.Errorf("fail to render: %s", err)
			}

			return nil
		})

		f.keybindingLogs(g)
		f.subscribe(g)
	}
	return nil

}

// presentation layer to display logs.
func (f *followerManager) render(v *gocui.View, datas [][]string) error {
	var width int

	// set widths for each column.
	width, _ = v.Size()

	t := tw.NewTableWriter(v)

	t.SetColumns([]string{"NAME", "MESSAGE"})
	t.SetColumnWidths([]int{40, width - 40})
	t.SetHeaderBorder(true)
	t.AppendBulk(datas)
	return t.Render()
}

func (f *followerManager) datas() [][]string {
	d := [][]string{}

	for _, l := range f.logs {
		dn, pod, message := l.DisplayName, l.Pod, l.Message
		podColor := f.podColor[pod]
		d = append(d, []string{color.ChangeColor(dn+":", podColor), message})
	}

	return d
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
			f.logger.Info("unfollow the workflow.")
			f.unfollow()
			f.bus.Publish(eventCollectionSetView)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

// subscribe events of follower.
const (
	eventFollowerSetView        = "follower:set-view"
	eventFollowerFollowWorkflow = "follower:follow-workflow"
)

func (f *followerManager) subscribe(g *gocui.Gui) error {
	if err := f.bus.Subscribe(eventFollowerSetView, func() {
		f.logger.Info("set the follower current view.")
		g.SetViewOnTop(followerViewName)
		g.SetCurrentView(followerViewName)
	}); err != nil {
		return err
	}

	if err := f.bus.Subscribe(eventFollowerFollowWorkflow, func(key string) {
		f.logger.Infof("follow the workflow: %s.", key)
		f.follow(key)

		f.logger.Infof("init the cursor.")
		v, _ := g.View(followerViewName)
		v.SetCursor(0, followerUpperBound)
		v.SetOrigin(0, 0)
	}); err != nil {
		return err
	}

	return nil
}
