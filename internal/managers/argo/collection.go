package argo

import (
	"fmt"
	"time"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	h "github.com/argoproj/pkg/humanize"
	"github.com/asaskevich/EventBus"
	"github.com/hanjunlee/argocui/pkg/argo"
	tw "github.com/hanjunlee/argocui/pkg/table/tablewriter"
	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/cache"
)

const (
	collectionViewName   = "collection"
	collectionUpperBound = 2
)

type collectionManager struct {
	namespace   string
	namePattern string
	cache       []*wf.Workflow

	uc  argo.UseCase
	bus EventBus.Bus

	log *log.Entry
}

func newCollectionManager(uc argo.UseCase, bus EventBus.Bus) *collectionManager {
	return &collectionManager{
		namespace:   "*",
		namePattern: "",
		cache:       []*wf.Workflow{},
		uc:          uc,
		bus:         bus,
		log: log.WithFields(log.Fields{
			"pkg": "argo manager",
		}),
	}
}

func (cm *collectionManager) pattern() string {
	return fmt.Sprintf("%s/*%s*", cm.namespace, cm.namePattern)
}

// lay out the collection view.
func (cm *collectionManager) layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	var (
		period = 1 * time.Second
	)

	v, err := g.SetView(collectionViewName, x0, y0, x1, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// settings of list view
		v.Highlight = true
		v.SelBgColor = gocui.ColorYellow
		v.SelFgColor = gocui.ColorBlack
		v.Frame = false
		v.SetCursor(0, collectionUpperBound)

		// focus on
		g.SetCurrentView(collectionViewName)

		// set refresh
		go viewutil.RefreshViewPeriodic(g, v, period, func() error {
			v.Clear()

			wfs := cm.uc.Search(cm.pattern())
			cm.cache = wfs

			err := render(v, toRows(wfs))
			if err != nil {
				fmt.Fprintf(v, "error occurs: %s", err)
			}
			return nil
		})

		cm.keybinding(g)
		cm.subscribe(g)
	}

	return nil
}

// render present workflows as table-like format.
func render(v *gocui.View, datas [][]string) error {
	var width, nameWidth int

	// set widths for each column.
	width, _ = v.Size()
	nameWidth = width - 70
	if nameWidth < 45 {
		nameWidth = 45
	}

	t := tw.NewTableWriter(v)

	t.SetColumns([]string{"NAMESPACE", "NAME", "STATUS", "AGE", "DURATION"})
	t.SetColumnWidths([]int{30, nameWidth, 20, 10, 10})
	t.SetHeaderBorder(true)
	t.AppendBulk(datas)
	return t.Render()
}

func toRows(wfs []*wf.Workflow) [][]string {
	var (
		rows = make([][]string, 0)
	)

	for _, w := range wfs {
		var (
			namespace = w.GetNamespace()
			name      = w.GetName()
			status    = argoutil.WorkflowStatus(w)
			age       = h.RelativeDurationShort(w.ObjectMeta.CreationTimestamp.Time, time.Now())
			duration  = h.RelativeDurationShort(w.Status.StartedAt.Time, w.Status.FinishedAt.Time)
		)
		rows = append(rows, []string{namespace, name, status, age, duration})
	}
	return rows
}

// keybinding of the colleciton view
func (cm *collectionManager) keybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding(collectionViewName, '/', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			const (
				eventSetSearchView = "search:set-view"
			)
			cm.log.Debugf("publish the event: search: %s", eventSetSearchView)
			cm.bus.Publish(eventSetSearchView)
			return nil
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(collectionViewName, 'k', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorUp(g, v, collectionUpperBound)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(collectionViewName, 'j', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorDown(g, v)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(collectionViewName, 'H', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorTop(g, v, collectionUpperBound)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding(collectionViewName, 'L', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return viewutil.MoveCursorBottom(g, v)
		}); err != nil {
		return err
	}

	// if err := g.SetKeybinding(collectionViewName, gocui.KeyCtrlL, gocui.ModNone,
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

	if err := g.SetKeybinding(collectionViewName, gocui.KeyBackspace2, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			_, py, _ := viewutil.GetCursorPosition(g, v)
			key, err := cm.getKeyAtCursor(py)
			if err != nil {
				cm.log.Errorf("fail to get key: %s", err)
				return nil
			}

			cm.log.Infof("delete the workflow: %s.", key)
			err = cm.uc.Delete(key)
			if err != nil {
				cm.log.Errorf("fail to delete the workflow: %s", err)
				return nil
			}

			return nil
		}); err != nil {
		return err
	}
	return nil
}

func (cm *collectionManager) getKeyAtCursor(cursor int) (string, error) {
	const (
		collectionUpperBound = 2
	)

	idx := cursor - collectionUpperBound
	if idx < 0 || idx > len(cm.cache) {
		return "", fmt.Errorf("cursor out of range: %d", cursor)
	}

	w := cm.cache[idx]
	key, _ := cache.MetaNamespaceKeyFunc(w)
	return key, nil
}

// subscribes of the collection view.
const (
	eventCollectionSet         = "collection:set-view"
	eventCollectionNamePattern = "collection:set-name-pattern"
)

func (cm *collectionManager) subscribe(g *gocui.Gui) error {
	const (
		viewName = "collection"
	)

	if err := cm.bus.Subscribe(eventCollectionSet, func() {
		cm.log.Info("set the current view list.")
		g.SetCurrentView(collectionViewName)
	}); err != nil {
		return err
	}

	if err := cm.bus.Subscribe(eventCollectionNamePattern, func(pattern string) {
		if pattern == cm.namePattern {
			return
		}

		cm.log.Infof("set the name of pattern %s.", pattern)
		cm.namePattern = pattern

		cm.log.Info("init cursor of the view.")
		v, _ := g.View(collectionViewName)

		v.SetOrigin(0, 0)
		v.SetCursor(0, collectionUpperBound)
	}); err != nil {
		return err
	}
	return nil
}
