package list

import (
	"fmt"

	"github.com/asaskevich/EventBus"
	"github.com/hanjunlee/argocui/pkg/argo"
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"
	"github.com/jroimartin/gocui"
	"k8s.io/client-go/tools/cache"
)

// Keybinding bind keys on the list view.
func Keybinding(g *gocui.Gui, s argo.UseCase, bus *EventBus.EventBus) error {
	if err := g.SetKeybinding(viewName, '/', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			log.Debug("publish the event: search:set-view.")
			bus.Publish("search:set-view")
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
			key, err := getKeyAtCursor(py)
			if err != nil {
				log.Errorf("fail to get key: %s", err)
				return nil
			}

			log.Infof("delete the workflow: %s.", key)
			err = s.Delete(key)
			if err != nil {
				log.Errorf("fail to delete the workflow: %s", err)
				return nil
			}

			return nil
		}); err != nil {
		return err
	}
	return nil
}

func getKeyAtCursor(cursor int) (string, error) {
	idx := cursor - upperBoundOfCursor
	if idx < 0 || idx > len(conf.cache) {
		return "", fmt.Errorf("cursor out of range: %d", cursor)
	}

	w := conf.cache[idx]
	key, _ := cache.MetaNamespaceKeyFunc(w)
	return key, nil
}
