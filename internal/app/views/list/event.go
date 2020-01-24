package list

import (
	"github.com/asaskevich/EventBus"
	"github.com/jroimartin/gocui"
)

const (
	eventSetView = "list:set-view"
	eventSetNamePattern = "list:set-name-pattern"
)

// Subscribe set events to be triggered in other views.
func Subscribe(g *gocui.Gui, bus EventBus.Bus) error {
	if err := bus.Subscribe(eventSetView, func() {
		log.Info("set the current view list.")
		g.SetCurrentView(viewName)
	}); err != nil {
		return err
	}

	if err := bus.Subscribe(eventSetNamePattern, func(pattern string) {
		if pattern == conf.namePattern {
			return
		}

		log.Infof("set the name of pattern %s.", pattern)
		conf.namePattern = pattern

		log.Info("init cursor of the view.")
		v, _ := g.View(viewName)

		v.SetOrigin(0, 0)
		v.SetCursor(0, upperBoundOfCursor)
	}); err != nil {
		return err
	}
	return nil
}
