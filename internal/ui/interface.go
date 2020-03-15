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
	NewDector(g *gocui.Gui, init string) error
	ReturnDector(g *gocui.Gui) (string error)

	// Switcher
	NewSwitcher(g *gocui.Gui) error
	ReturnSwitcher(g *gocui.Gui) (svc.UseCase, error)

	// Informer
	NewInformer(g *gocui.Gui) error
	ReturnInformer(g *gocui.Gui) error

	// Remover
	NewRemover(g *gocui.Gui, key string) error
	ReturnRemover(g *gocui.Gui, delete bool) error
}

// Presentor present the resource looks pretty.
type Presentor interface {
	PresentCore(*gocui.View, []runtime.Object) error 
	PresentInformer(*gocui.View, runtime.Object) error
}
