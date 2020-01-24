package etc

import (
	"fmt"
	"time"

	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	"github.com/hanjunlee/argocui/pkg/util/color"
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"
	"github.com/jroimartin/gocui"
	log "github.com/sirupsen/logrus"
)

// Manager is the manager of etc.
type Manager struct {
	log *log.Entry
}

// NewManager create a new manager.
func NewManager() *Manager {
	return &Manager{
		log: log.WithFields(log.Fields{
			"pkg": "etc",
		}),
	}
}

// Layout lay out the info and the brand.
func (m *Manager) Layout(g *gocui.Gui) error {
	x, y := g.Size()
	grid := x / 12
	padding := grid
	if err := m.layoutInfo(g, padding, 0, 5*grid-1, y/4-1); err != nil {
		return err
	}

	if err := m.layoutBrand(g, 6*grid, 0, x-1, y/4-1); err != nil {
		return err
	}

	if err := m.keybinding(g); err != nil {
		return err
	}
	return nil
}

func (m *Manager) layoutInfo(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView("info", x0, y0, x1, y1)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false

		go viewutil.RefreshViewPeriodic(g, v, 1*time.Second, func() error {
			const (
				argoRev    = "v2.4.1"
				argocuiRev = "v0.0.1"
				homePage   = "github.com/hanjunlee/argocui"
			)
			var (
				context   string
				namespace string
				err       error
			)

			if context, err = argoutil.GetCurrentContext(); err != nil {
				m.log.Errorf("couldn't get the context: %s.", err)
				context = ""
			}

			if namespace, err = argoutil.GetNamespace(); err != nil {
				m.log.Errorf("couldn't get the namespace: %s.", err)
				namespace = ""
			}

			v.Clear()
			fmt.Fprintln(v, "")
			fmt.Fprintf(v, "Context:      %s\n", color.ChangeColor(context, gocui.ColorYellow))
			fmt.Fprintf(v, "Namespace:    %s\n", color.ChangeColor(namespace, gocui.ColorYellow))
			fmt.Fprintf(v, "Argo Rev:     %s\n", color.ChangeColor(argoRev, gocui.ColorYellow))
			fmt.Fprintf(v, "Argocui Rev:  %s\n", color.ChangeColor(argocuiRev, gocui.ColorYellow))
			fmt.Fprintf(v, "Homepage:     %s\n", color.ChangeColor(homePage, gocui.ColorYellow))
			return nil
		})
	}

	return nil
}

func (m *Manager) layoutBrand(g *gocui.Gui, x0, x1, y0, y1 int) error {
	v, err := g.SetView("brand", x0, x1, y0, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		brand := `
     _____                            _________       .__ 
    /  _  \_______  ____   ____       \_   ___ \ __ __|__|
   /  /_\  \_  __ \/ ___\ /  _ \ --   /    \  \/|  |  \  |
  /    |    \  | \/ /_/  >  <_> ) --- \     \___|  |  /  |
  \____|__  /__|  \___  / \____/ --    \______  /____/|__|
       	  \/     /_____/                      \/          
`
		fmt.Fprintf(v, color.ChangeColor(brand, gocui.ColorYellow))
	}

	return nil
}

func (m *Manager) keybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
	}); err != nil {
		return err
	}
	return nil
}