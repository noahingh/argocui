package search 

import (
	"github.com/asaskevich/EventBus"
	"github.com/jroimartin/gocui"
)

const (
	// SetView is the name of event to set the current view list.
	SetView = "search:set-view"
)

// Subscribe set events to be triggered in other views.
func Subscribe(g *gocui.Gui, bus *EventBus.EventBus) error {
	if err := bus.Subscribe(SetView, func() {
		log.Info("set the current view search.")
		g.SetCurrentView(viewName)
	}); err != nil {
		return err
	}
	return nil
}
