package ui

import (
	"fmt"
	"strings"

	runtime "github.com/hanjunlee/argocui/pkg/runtime"
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"
	"k8s.io/client-go/tools/cache"

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
	// Remover is the remover view.
	Remover string = "remover"
)

// Manager is the manager of UI.
type Manager struct {
	Svc        runtime.UseCase
	SvcEntries map[string]runtime.UseCase

	// namespace is the context of the manager.
	namespace string
	// Cache keys of runtime object after search query.
	cache []string

	// dected is the string dected by the Dector.
	dected string

	// removed is the key which is removed.
	removed string
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

	m.cache = make([]string, 0)
	for _, o := range m.Svc.Search(m.namespace, m.dected) {
		gvk := o.GetObjectKind().GroupVersionKind()
		switch gvk.Kind {
		case "Mock":
			key, _ := cache.MetaNamespaceKeyFunc(o)
			m.cache = append(m.cache, key)
			fmt.Fprintln(v, key)
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
			return m.NewDector(g, m.dected)
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

	if err := g.SetKeybinding(Core, gocui.KeyBackspace2, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, y, _ := viewutil.GetCursorPosition(g, v)
			if y >= len(m.cache) {
				log.Error("couldn't delete: the cursor is out of range.")
				return nil
			}

			log.Infof("switch to the remover: %s", m.cache[y])
			m.NewRemover(g, m.cache[y])
			return nil
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
			m.dected = dected
			log.Infof("detect and set the word: %s", dected)

			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Dector, gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			m.ReturnDector(g)
			m.dected = ""
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

	// Remover keybinding
	if err := g.SetKeybinding(Remover, gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if err := m.ReturnRemover(g, true); err != nil {
				log.Errorf("failed to delete: %s", err)
			}
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(Remover, gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			m.ReturnRemover(g, false)

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
	v, _ := g.View(Core)
	defer g.SetCurrentView(Core)
	defer v.SetOrigin(0, 0)
	defer v.SetCursor(0, 0)
	defer g.DeleteView(Dector)

	v, _ = g.View(Dector)
	s, _ := v.Line(0)
	s = strings.TrimSpace(s)

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
func (m *Manager) ReturnSwitcher(g *gocui.Gui) (runtime.UseCase, error) {
	v, _ := g.View(Core)
	defer g.SetCurrentView(Core)
	defer v.SetOrigin(0, 0)
	defer v.SetCursor(0, 0)
	defer g.DeleteView(Switcher)

	v, _ = g.View(Switcher)
	s, _ := v.Line(0)
	s = strings.TrimSpace(s)

	svc, ok := m.SvcEntries[s]
	if !ok {
		return nil, fmt.Errorf("there is no service: %s", s)
	}
	return svc, nil
}

// NewRemover switch to the remover and confirm to delete or not.
func (m *Manager) NewRemover(g *gocui.Gui, key string) error {
	m.removed = key

	w, h := g.Size()
	v, err := g.SetView(Remover, 0, h/4, w-1, h/4+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Delete(Y/n)?"
		v.FgColor = gocui.ColorRed
		v.Editable = true
		v.Editor = gocui.EditorFunc(inlineEditor)

		g.SetCurrentView(Remover)
	}

	return nil
}

// ReturnRemover switch to the Core.
func (m *Manager) ReturnRemover(g *gocui.Gui, delete bool) error {
	defer g.SetCurrentView(Core)
	defer g.DeleteView(Remover)

	v, _ := g.View(Remover)
	s, _ := v.Line(0)
	s = strings.TrimSpace(s)

	if !delete {
		return nil
	}
	if s != "Y" && s != "y" {
		return nil
	}

	if err := m.Svc.Delete(m.removed); err != nil {
		return err
	}
	return nil
}
