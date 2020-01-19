package app

import (
	"github.com/jroimartin/gocui"

	"github.com/hanjunlee/argocui/internal/app/views/list"
	"github.com/hanjunlee/argocui/pkg/argo"
)

// ConfigureGui is
func ConfigureGui(g *gocui.Gui) {
	// settings of gui
	g.Highlight = true
	g.SelFgColor = gocui.ColorYellow
	g.InputEsc = true
}

// ManagerFunc return the manager function.
func ManagerFunc(s argo.UseCase, g *gocui.Gui) func(g *gocui.Gui) error {
	return func(g *gocui.Gui) error {
		var (
			err error
		)

		maxX, maxY := g.Size()

		// err = layoutInfo(g, 1, 0, maxX/5-1, maxY/4-1)
		// if err != nil {
		// 	return err
		// }

		// err = layoutBrand(g, maxX/5, 0, maxX-1, maxY/4-1)
		// if err != nil {
		// 	return err
		// }

		// err = layoutSearch(g, 0, maxY/4-2, maxX-1, maxY/4)
		// if err != nil {
		// 	return err
		// }

		err = list.Layout(s, g, 0, maxY/4+1, maxX-1, maxY-1)
		if err != nil {
			return err
		}
		return nil
	}
}
