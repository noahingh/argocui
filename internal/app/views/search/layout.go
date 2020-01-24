package search

import (
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"
	"github.com/hanjunlee/argocui/pkg/argo"
	"github.com/jroimartin/gocui"
	"github.com/asaskevich/EventBus"
)

const (
	viewName = "search"
)

// Layout lay out the search view.
func (c *Config) Layout(g *gocui.Gui, s argo.UseCase, bus EventBus.Bus, x0, y0, x1, y1 int) error {
	v, err := g.SetView(viewName, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// settings of search view
		v.Title = "Search"
		v.FgColor = gocui.ColorYellow
		v.Editable = true
		v.Editor = gocui.EditorFunc(viewutil.LineEditor)

		c.keybinding(g, s, bus)
		c.subscribe(g, bus)
	}

	return nil
}
