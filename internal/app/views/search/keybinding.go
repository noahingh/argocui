package search

import (
	"strings"

	"github.com/asaskevich/EventBus"
	"github.com/hanjunlee/argocui/pkg/argo"
	"github.com/jroimartin/gocui"
)

// Keybinding the keybinding of the search view.
func (c *Config) Keybinding(g *gocui.Gui, s argo.UseCase, bus EventBus.Bus) error {
	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			const (
				eventNamePattern = "list:set-name-pattern"
				eventSetView     = "list:set-view"
			)
			pattern, _ := v.Line(0)
			pattern = strings.TrimSpace(pattern)

			c.log.Debug("publish the event: list: %s.", eventNamePattern)
			bus.Publish(eventNamePattern, pattern)

			c.log.Debug("publish the event: list: %s.", eventSetView)
			bus.Publish(eventSetView)
			return nil
		}); err != nil {
		return err
	}
	return nil
}
