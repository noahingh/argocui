package ui

import (
	resource "github.com/hanjunlee/argocui/pkg/runtime"
	"github.com/jroimartin/gocui"
)

// ManagerIface is the interface of manager.
type ManagerIface interface {
	Layout(g *gocui.Gui) error
	Keybinding(g *gocui.Gui) error

	// Dector
	//
	// NewDector switch to the Dector.
	NewDector(g *gocui.Gui, init string) error
	// ReturnDector return the result and switch to the core view.
	ReturnDector(g *gocui.Gui) (string error)

	// Switcher
	//
	// NewSwitcher switch services.
	NewSwitcher(g *gocui.Gui) error
	// ReturnSwitcher return the service.
	ReturnSwitcher(g *gocui.Gui) (resource.UseCase, error)
}
