package app

import (
	"github.com/jroimartin/gocui"
)

// ConfigureGui is
func ConfigureGui(g *gocui.Gui) {
	// settings of gui
	g.Highlight = true
	g.SelFgColor = gocui.ColorYellow
	g.InputEsc = true
}

// Layout is
func Layout(g *gocui.Gui) error {
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

	err = layoutList(g, 0, maxY/4+1, maxX-1, maxY-1)
	if err != nil {
		return err
	}

	return nil
}

func getMainViewSize(g *gocui.Gui) (x0, y0, x1, y1 int) {
	maxX, maxY := g.Size()
	x0 = 0
	y0 = maxY/4 + 1
	x1 = maxX - 1
	y1 = maxY - 1
	return x0, y0, x1, y1
}
