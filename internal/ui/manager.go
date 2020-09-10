package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hanjunlee/argocui/internal/config"
	"github.com/hanjunlee/argocui/internal/runtime"
	"github.com/hanjunlee/argocui/internal/ui/mock"
	"github.com/hanjunlee/argocui/internal/ui/namespace"
	"github.com/hanjunlee/argocui/internal/ui/workflow"
	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	colorutil "github.com/hanjunlee/argocui/pkg/util/color"
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"

	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/cache"
)

const (
	// Core is the core view.
	viewCore string = "core"
	// Dector is the dector view.
	Dector string = "dector"
	// Switcher is the switcher view.
	Switcher string = "switcher"
	// Informer is the informer view.
	Informer string = "informer"
	// Follower is th follower view.
	Follower string = "follower"
	// Remover is the remover view.
	Remover string = "remover"
	// Messenger is the messenger view.
	Messenger string = "messenger"

	// headerSize is the size of table header.
	headerSize = 1
)

// Manager is the manager of UI.
type Manager struct {
	svc          runtime.UseCase
	mockSvc      runtime.UseCase
	namespaceSvc runtime.UseCase
	workflowSvc  runtime.UseCase

	// search
	// namespace is the context of the manager.
	namespace string
	// cache is keys of runtime object after search query.
	cache []string

	// dector
	// dected is the string dected by the Dector.
	dected string

	// follower
	logs   []runtime.Log
	cancel context.CancelFunc

	// remover
	// removed is the key which is removed.
	removed string
}

// NewManager create a new UI manager. The namespace of the manager is depends on the configuration of the user.
func NewManager(mock runtime.UseCase, namespace runtime.UseCase, workflow runtime.UseCase) *Manager {
	ns, _ := argoutil.GetNamespace()

	return &Manager{
		svc:          workflow,
		mockSvc:      mock,
		namespaceSvc: namespace,
		workflowSvc:  workflow,
		namespace:    ns,
	}
}

// Layout lay out the resource of service.
func (m *Manager) Layout(g *gocui.Gui) error {
	w, h := g.Size()

	v, err := g.SetView("context", 0, 1, w/3-1, h/5-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false

	}
	v.Clear()
	m.layoutContext(v)

	v, err = g.SetView("brand", w/2, 1, w-1, h/5-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		fmt.Fprintf(v, colorutil.ChangeColor(config.Logo, gocui.ColorYellow))
	}

	// messenger
	v, err = g.SetView(Messenger, 0, h-2, w-1, h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
	}

	// core
	v, err = g.SetView(viewCore, 0, h/5+3, w-1, h-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.Frame = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.SetCursor(0, headerSize)

		g.SetCurrentView(viewCore)
	}

	objs := m.svc.Search(m.namespace, m.dected)

	// cache first.
	m.cache = make([]string, 0)
	for _, o := range objs {
		key, _ := cache.MetaNamespaceKeyFunc(o)
		m.cache = append(m.cache, key)
	}

	// presentor
	var p Presentor

	if m.svc == m.mockSvc {
		p = mock.NewPresentor()
	} else if m.svc == m.namespaceSvc {
		p = namespace.NewPresentor()
	} else if m.svc == m.workflowSvc {
		p = workflow.NewPresentor()
	}

	v.Clear()
	p.PresentCore(v, objs)

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
	if err := g.SetKeybinding(viewCore, gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, y, _ := viewutil.GetCursorPosition(g, v)
			y = y - headerSize
			if y >= len(m.cache) {
				log.Error("the cursor is out of range.")
				return nil
			}

			o, err := m.svc.Get(m.cache[y])
			if err != nil {
				log.Errorf("failed to get the object: %s", err)
				return nil
			}

			if m.svc != m.namespaceSvc {
				return nil
			}
			m.namespace, _ = cache.MetaNamespaceKeyFunc(o)
			m.svc = m.workflowSvc
			log.Infof("switch namespace: %s", m.namespace)
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewCore, 'k', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorUp(g, v, headerSize)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewCore, 'j', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorDown(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewCore, 'H', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorTop(g, v, headerSize)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewCore, 'L', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorBottom(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewCore, '/', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			log.Infof("new dector")
			return m.NewDector(g, m.dected)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewCore, ':', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			log.Infof("new switcher")
			return m.NewSwitcher(g)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewCore, gocui.KeyBackspace2, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, y, _ := viewutil.GetCursorPosition(g, v)
			y = y - headerSize
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

	if err := g.SetKeybinding(viewCore, gocui.KeyCtrlG, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, y, _ := viewutil.GetCursorPosition(g, v)
			y = y - headerSize
			if y >= len(m.cache) {
				log.Error("couldn't delete: the cursor is out of range.")
				return nil
			}

			key := m.cache[y]
			if m.svc == m.mockSvc {
				m.Warn(g, "sorry, animal is not implemented yet.")
			} else if m.svc == m.namespaceSvc {
				m.Warn(g, "sorry, namespace is not implemented yet.")
			} else if m.svc == m.workflowSvc {
				log.Infof("switch to the informer: %s", key)
				m.NewInformer(g, key)
			}

			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewCore, gocui.KeyCtrlL, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, y, _ := viewutil.GetCursorPosition(g, v)
			y = y - headerSize
			if y >= len(m.cache) {
				log.Error("couldn't delete: the cursor is out of range.")
				return nil
			}

			key := m.cache[y]
			o, err := m.svc.Get(key)
			if err != nil {
				log.Errorf("failed to get the object: %s", err)
				return nil
			}

			gvk, _, _ := objectKind(o)
			switch gvk.Kind {
			case "Animal":
				m.Warn(g, "sorry, animal is not implemented yet.")
			case "Namespace":
				m.Warn(g, "sorry, namespace couldn't follow up.")
			case "Workflow":
				log.Infof("switch to the follower: %s", key)
				m.NewFollower(g, key)
			default:
				return nil
			}

			return nil
		}); err != nil {
		return err
	}

	// Dector keybinding
	if err := g.SetKeybinding(Dector, gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			dected, err := m.ReturnDector(g)
			if err != nil {
				log.Errorf("failed to search: %s", err)
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
				m.Error(g, fmt.Sprintf("failed to switch: %s", err))
				return nil
			}
			m.svc = svc
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
				m.Error(g, fmt.Sprintf("failed to switch: %s", err))
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
	v, err := g.SetView(Dector, 0, h/5, w-1, h/5+2)
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

// ReturnDector return the result from the dector and back to the viewCore.
func (m *Manager) ReturnDector(g *gocui.Gui) (string, error) {
	v, _ := g.View(viewCore)
	defer g.SetCurrentView(viewCore)
	defer v.SetOrigin(0, 0)
	defer v.SetCursor(0, headerSize)
	defer g.DeleteView(Dector)

	v, _ = g.View(Dector)
	s, _ := v.Line(0)
	s = strings.TrimSpace(s)

	return s, nil
}

// NewSwitcher create and switch to the Switcher
func (m *Manager) NewSwitcher(g *gocui.Gui) error {
	w, h := g.Size()
	v, err := g.SetView(Switcher, 0, h/5, w-1, h/5+2)
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

// ReturnSwitcher return the service from the switcher and back to the viewCore.
func (m *Manager) ReturnSwitcher(g *gocui.Gui) (runtime.UseCase, error) {
	v, _ := g.View(viewCore)
	defer g.SetCurrentView(viewCore)
	defer v.SetOrigin(0, 0)
	defer v.SetCursor(0, headerSize)
	defer g.DeleteView(Switcher)

	v, _ = g.View(Switcher)
	s, _ := v.Line(0)
	s = strings.TrimSpace(s)

	var (
		svc runtime.UseCase
		err error
	)
	switch s {
	case "mock":
		svc = m.mockSvc
	case "ns":
		svc = m.namespaceSvc
	case "wf":
		svc = m.workflowSvc
	default:
		svc = nil
		err = fmt.Errorf("there is no service: %s", s)
	}
	return svc, err
}

// NewRemover switch to the remover and confirm to delete or not.
func (m *Manager) NewRemover(g *gocui.Gui, key string) error {
	m.removed = key

	w, h := g.Size()
	v, err := g.SetView(Remover, 0, h/5, w-1, h/5+2)
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

// ReturnRemover switch to the viewCore.
func (m *Manager) ReturnRemover(g *gocui.Gui, delete bool) error {
	defer g.SetCurrentView(viewCore)
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

	if err := m.svc.Delete(m.removed); err != nil {
		return err
	}
	return nil
}

// Warn show up the message on the Messenger.
// It's recommended to use in GUI level such as keybinding and laytout.
func (m *Manager) Warn(g *gocui.Gui, message string) {
	v, _ := g.View(Messenger)
	v.Clear()

	message = colorutil.ChangeColor(message, gocui.ColorYellow)
	v.Write([]byte(message))
	go func() {
		time.Sleep(2 * time.Second)
		v.Clear()
	}()
}

// Error show up the message on the Messenger.
// It's recommended to use in GUI level such as keybinding and laytout.
func (m *Manager) Error(g *gocui.Gui, message string) {
	v, _ := g.View(Messenger)
	v.Clear()

	message = colorutil.ChangeColor(message, gocui.ColorRed)
	v.Write([]byte(message))
	go func() {
		time.Sleep(2 * time.Second)
		v.Clear()
	}()
	return
}
