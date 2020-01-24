package list

import (
	"fmt"

	"github.com/asaskevich/EventBus"
	"github.com/hanjunlee/argocui/pkg/argo"
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"
	"github.com/jroimartin/gocui"
	"k8s.io/client-go/tools/cache"
)

func (c *Config) keybinding(g *gocui.Gui, s argo.UseCase, bus EventBus.Bus) error {
	if err := g.SetKeybinding(viewName, '/', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			const (
				eventSetSearchView = "search:set-view"
			)
			c.log.Debugf("publish the event: search: %s", eventSetSearchView)
			bus.Publish(eventSetSearchView)
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewName, 'k', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorUp(g, v, upperBoundOfCursor)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewName, 'j', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorDown(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewName, 'H', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorTop(g, v, upperBoundOfCursor)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(viewName, 'L', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorBottom(g, v)
		}); err != nil {
		return err
	}

	// if err := g.SetKeybinding(viewName, gocui.KeyCtrlL, gocui.ModNone,
	// 	func(g *gocui.Gui, v *gocui.View) error {
	// 		_, py, _ := viewutil.GetCursorPosition(g, v)
	// 		key, err := uiClientset.List.GetKeyAtLine(py)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		log.Infof("set the client of logs '%s'", key)
	// 		uiClientset.Logs.SetKey(key)

	// 		lv, err := newLogsView(g)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		g.SetCurrentView(lv.Name())
	// 		return nil
	// 	}); err != nil {
	// 	return err
	// }

	if err := g.SetKeybinding(viewName, gocui.KeyCtrlZ, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, py, _ := viewutil.GetCursorPosition(g, v)
			key, err := c.getKeyAtCursor(py)
			if err != nil {
				c.log.Errorf("fail to get key: %s", err)
				return nil
			}

			c.log.Infof("delete the workflow: %s.", key)
			err = s.Delete(key)
			if err != nil {
				c.log.Errorf("fail to delete the workflow: %s", err)
				return nil
			}

			return nil
		}); err != nil {
		return err
	}
	return nil
}

func (c *Config) getKeyAtCursor(cursor int) (string, error) {
	idx := cursor - upperBoundOfCursor
	if idx < 0 || idx > len(c.cache) {
		return "", fmt.Errorf("cursor out of range: %d", cursor)
	}

	w := c.cache[idx]
	key, _ := cache.MetaNamespaceKeyFunc(w)
	return key, nil
}
