package search

import (
	"strings"

	"github.com/asaskevich/EventBus"
	"github.com/hanjunlee/argocui/pkg/argo"
	"github.com/jroimartin/gocui"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{
		"pkg": "search",
	})
)

// Keybinding the keybinding of the search view.
func Keybinding(g *gocui.Gui, s argo.UseCase, bus EventBus.Bus) error {
	if err := g.SetKeybinding(viewName, gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			pattern, _ := v.Line(0)
			pattern = strings.TrimSpace(pattern)

			log.Debug("publish the event: list:set-name-pattern.")
			bus.Publish("list:set-name-pattern.", pattern)

			log.Debug("publish the event: list:set-view.")
			bus.Publish("list:set-view")
			return nil
		}); err != nil {
		return err
	}
	return nil
}
