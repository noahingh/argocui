package argo

import (
	"strings"

	"github.com/asaskevich/EventBus"
	"github.com/hanjunlee/argocui/pkg/argo"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

type subManager struct {
	uc  argo.UseCase
	bus EventBus.Bus

	log *log.Entry
}

func newSubManager(uc argo.UseCase, bus EventBus.Bus) *subManager {
	return &subManager{
		uc:  uc,
		bus: bus,
		log: log.WithFields(log.Fields{
			"pkg":  "argo-manager",
			"file": "sub.go",
		}),
	}
}

const (
	subViewName = "sub"
)

func (s *subManager) layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(subViewName, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Title = "Search"
		v.FgColor = gocui.ColorYellow
		v.Editable = true
		v.Editor = gocui.EditorFunc(subEditor)

		s.keybinding(g)
		s.subscribe(g)
	}

	return nil
}

func subEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyCtrlU:
		for true {
			line, _ := v.Line(0)
			if len(line) == 0 {
				break
			}

			v.EditDelete(true)
		}
	}
}

func trimLine(l string) string {
	l = strings.TrimSpace(l)
	return l
}

func (s *subManager) keybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding(subViewName, gocui.KeyEnter, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			pattern, _ := v.Line(0)
			pattern = trimLine(pattern)

			s.log.Debugf("publish the event: list: %s.", eventCollectionSetNamePattern)
			s.bus.Publish(eventCollectionSetNamePattern, pattern)

			s.log.Debugf("publish the event: list: %s.", eventCollectionSetView)
			s.bus.Publish(eventCollectionSetView)
			return nil
		}); err != nil {
		return err
	}
	return nil
}

const (
	eventSubSetView = "sub:set-view"
)

// subscribe set events to be triggered in other views.
func (s *subManager) subscribe(g *gocui.Gui) error {
	if err := s.bus.Subscribe(eventSubSetView, func() {
		s.log.Info("set the current view search.")
		g.SetCurrentView(subViewName)
	}); err != nil {
		return err
	}
	return nil
}
