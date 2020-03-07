package ui

import (
	"fmt"
	"strings"

	"github.com/hanjunlee/argocui/pkg/resource"
	"github.com/hanjunlee/argocui/pkg/runtime/mock"
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"
	"github.com/willf/pad"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

const (
	// Core is the core view.
	Core string = "core"
	// Dector is the dector view.
	Dector string = "dector"
	// Switcher is the switcher view.
	Switcher string = "switcher"
)

// Manager is the manager of UI.
type Manager struct {
	Svc        resource.UseCase
	SvcEntries map[string]resource.UseCase

	// Dected is the string dected by the Dector.
	Dected string
}

// Layout lay out the resource of service.
func (m *Manager) Layout(g *gocui.Gui) error {
	w, h := g.Size()

	v, err := g.SetView(Core, 0, h/4+3, w-1, h-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Frame = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack

		g.SetCurrentView(Core)
	}

	v.Clear()
	// TODO: print the content.
	switch m.Svc.GetRepoType() {
	case resource.Mock:
		w, _ := v.Size()
		for _, a := range m.Svc.Search(m.Dected) {
			s := pad.Right(string(a.(mock.Animal)), w, " ")
			fmt.Fprintln(v, s)
		}
	}

	return nil
}

// Keybinding keybinding of views in the manager.
func (m *Manager) Keybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}

	// Core keybinding
	if err := g.SetKeybinding(Core, 'k', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorUp(g, v, 0)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Core, 'j', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorDown(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Core, 'H', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorTop(g, v, 0)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Core, 'L', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorBottom(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Core, '/', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			log.Infof("new dector")
			return m.NewDector(g, m.Dected)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Core, ':', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			log.Infof("new switcher")
			return m.NewSwitcher(g)
		}); err != nil {
		return err
	}

	// Dector keybinding
	if err := g.SetKeybinding(Dector, gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			dected, err := m.ReturnDector(g)
			if err != nil {
				return err
			}
			m.Dected = dected
			log.Infof("detect and set the word: %s", dected)

			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Dector, gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			m.ReturnDector(g)
			m.Dected = ""
			log.Info("exit dector")

			return nil
		}); err != nil {
		return err
	}

	// Switcher keybinding
	if err := g.SetKeybinding(Switcher, gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			svc, err := m.ReturnSwitcher(g)
			if err != nil {
				log.Warnf("couldn't switch service: %s", err)
				return nil
			}
			m.Svc = svc
			log.Infof("switch the service")

			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Switcher, gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			m.ReturnSwitcher(g)
			log.Info("exit switcher")

			return nil
		}); err != nil {
		return err
	}

	return nil
}

// NewDector create and switch to the dector.
func (m *Manager) NewDector(g *gocui.Gui, init string) error {
	w, h := g.Size()
	v, err := g.SetView(Dector, 0, h/4, w-1, h/4+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Search"
		v.FgColor = gocui.ColorYellow
		v.Editable = true
		v.Editor = gocui.EditorFunc(inlineEditor)

		fmt.Fprint(v, init)
		v.SetCursor(len(init), 0)

		g.SetCurrentView(Dector)
	}

	return nil
}

// ReturnDector return the result from the dector and back to the Core.
func (m *Manager) ReturnDector(g *gocui.Gui) (string, error) {
	v, _ := g.View(Dector)
	s, _ := v.Line(0)
	s = strings.TrimSpace(s)

	g.DeleteView(Dector)

	v, _ = g.View(Core)
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
	g.SetCurrentView(Core)
	return s, nil
}

// NewSwitcher create and switch to the Switcher
func (m *Manager) NewSwitcher(g *gocui.Gui) error {
	w, h := g.Size()
	v, err := g.SetView(Switcher, 0, h/4, w-1, h/4+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Switch"
		v.FgColor = gocui.ColorCyan
		v.Editable = true
		v.Editor = gocui.EditorFunc(inlineEditor)

		g.SetCurrentView(Switcher)
	}

	return nil
}

// ReturnSwitcher return the service from the switcher and back to the Core.
func (m *Manager) ReturnSwitcher(g *gocui.Gui) (resource.UseCase, error) {
	v, _ := g.View(Switcher)
	s, _ := v.Line(0)
	s = strings.TrimSpace(s)

	g.DeleteView(Switcher)

	v, _ = g.View(Core)
	v.SetOrigin(0, 0)
	v.SetCursor(0, 0)
	g.SetCurrentView(Core)

	svc, ok := m.SvcEntries[s]
	if !ok {
		return nil, fmt.Errorf("there is no service: %s", s)
	}
	return svc, nil
}
