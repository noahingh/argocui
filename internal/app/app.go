package app

import (
	"github.com/jroimartin/gocui"

	"github.com/asaskevich/EventBus"
	"github.com/hanjunlee/argocui/internal/app/views/info"
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

		err = info.LayoutInfo(g, 1, 0, maxX/5-1, maxY/4-1)
		if err != nil {
			return err
		}

		err = info.LayoutBrand(g, maxX/5, 0, maxX-1, maxY/4-1)
		if err != nil {
			return err
		}

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

// Keybinding bind keys on GUI.
func Keybinding(s argo.UseCase, g *gocui.Gui, bus *EventBus.EventBus) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}

	// if err := keybindingSearch(g); err != nil {
	// 	return err
	// }

	if err := list.Keybinding(s, g, bus); err != nil {
		return err
	}

	// keybinding works apart from the view.
	// if err := keybindingLogs(g); err != nil {
	// 	return err
	// }

	return nil
}
