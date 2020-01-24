package list

import (
	"github.com/asaskevich/EventBus"
	"github.com/jroimartin/gocui"
)

const (
	// SetView is the name of event to set the current view list.
	SetView = "list:set-view"
	// SetNamePattern is the name of event to set the name pattern of the config.
	SetNamePattern = "list:set-name-pattern"
)

// Subscribe set events to be triggered in other views.
func Subscribe(g *gocui.Gui, bus EventBus.Bus) error {
	if err := bus.Subscribe(SetView, func() {
		log.Info("set the current view list.")
		g.SetCurrentView(viewName)
	}); err != nil {
		return err
	}

	if err := bus.Subscribe(SetNamePattern, func(pattern string) {
		log.Infof("set the name of pattern %s.", pattern)
		conf.namePattern = pattern
	}); err != nil {
		return err
	}
	return nil
}
