package search 

import (
	"github.com/asaskevich/EventBus"
	"github.com/jroimartin/gocui"
)

const (
	eventSetView = "search:set-view"
)

// subscribe set events to be triggered in other views.
func (c *Config) subscribe(g *gocui.Gui, bus EventBus.Bus) error {
	if err := bus.Subscribe(eventSetView, func() {
		c.log.Info("set the current view search.")
		g.SetCurrentView(viewName)
	}); err != nil {
		return err
	}
	return nil
}
