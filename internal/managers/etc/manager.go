package etc

import (
	"fmt"
	"time"

	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	"github.com/hanjunlee/argocui/pkg/util/color"
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"
	"github.com/hanjunlee/argocui/internal/config"

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
	if err := m.layoutInfo(g, 0, 1, 2*grid-1, y/4-1); err != nil {
		return err
	}

	if err := m.layoutHelp(g, 3*grid, 1, 6*grid-1, y/4-1); err != nil {
		return err
	}

	if err := m.layoutBrand(g, 6*grid, 1, 11*grid-1, y/4-1); err != nil {
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
				argoRev    = config.ArgoVersion
				argocuiRev = config.Version
				homePage   = config.HomePage
			)
			var (
				context string
				err     error
			)

			if context, err = argoutil.GetCurrentContext(); err != nil {
				m.log.Errorf("couldn't get the context: %s.", err)
				context = ""
			}

			v.Clear()
			fmt.Fprintln(v, "")
			fmt.Fprintf(v, "Context:      %s\n", color.ChangeColor(context, gocui.ColorYellow))
			fmt.Fprintf(v, "Argo Rev:     %s\n", color.ChangeColor(argoRev, gocui.ColorYellow))
			fmt.Fprintf(v, "Argocui Rev:  %s\n", color.ChangeColor(argocuiRev, gocui.ColorYellow))
			fmt.Fprintf(v, "Homepage:     %s\n", color.ChangeColor(homePage, gocui.ColorYellow))
			return nil
		})
	}

	return nil
}

func (m *Manager) layoutHelp(g *gocui.Gui, x0, x1, y0, y1 int) error {
	v, err := g.SetView("help", x0, x1, y0, y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Frame = false
		help := `
<H>: move to the top
<L>: move to the bottom
<k>: move up
<j>: move down
<ctrl+l>: follow log
<ctrl+g>: tree
</>: search 
<ctrl+del>: delete
`
		fmt.Fprintf(v, help)
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
		fmt.Fprintf(v, color.ChangeColor(config.Logo, gocui.ColorYellow))
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
