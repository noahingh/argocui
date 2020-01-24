package list

import (
	"github.com/asaskevich/EventBus"
	"github.com/jroimartin/gocui"
)

const (
	eventSetView = "list:set-view"
	eventSetNamePattern = "list:set-name-pattern"
)

// subscribe set events to be triggered in other views.
func (c *Config) subscribe(g *gocui.Gui, bus EventBus.Bus) error {
	if err := bus.Subscribe(eventSetView, func() {
		c.log.Info("set the current view list.")
		g.SetCurrentView(viewName)
	}); err != nil {
		return err
	}

	if err := bus.Subscribe(eventSetNamePattern, func(pattern string) {
		if pattern == c.namePattern {
			return
		}

		c.log.Infof("set the name of pattern %s.", pattern)
		c.namePattern = pattern

		c.log.Info("init cursor of the view.")
		v, _ := g.View(viewName)

		v.SetOrigin(0, 0)
		v.SetCursor(0, upperBoundOfCursor)
	}); err != nil {
		return err
	}
	return nil
}
