package ui

import (
	svc "github.com/hanjunlee/argocui/pkg/runtime"

	"github.com/jroimartin/gocui"
	"k8s.io/apimachinery/pkg/runtime"
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
	ReturnSwitcher(g *gocui.Gui) (svc.UseCase, error)

	// Remover
	//
	// NewRemover switch to the remove with the key to delete.
	NewRemover(g *gocui.Gui, key string) error
	// ReturnRemover switch to the core view.
	ReturnRemover(g *gocui.Gui, delete bool) error
}

// Presentor present the resource looks pretty.
type Presentor interface {
	PresentCore(*gocui.View, []runtime.Object) error 
}
