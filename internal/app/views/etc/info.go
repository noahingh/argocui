package etc

import (
	"fmt"
	"time"

	argoutil "github.com/hanjunlee/argocui/pkg/util/argo"
	"github.com/hanjunlee/argocui/pkg/util/color"
	viewutil "github.com/hanjunlee/argocui/pkg/util/view"
	"github.com/jroimartin/gocui"
	"github.com/sirupsen/logrus"
)

const (
	infoViewName  = "info"
)

// InfoConfig is the configuration of a info view.
type InfoConfig struct {
	log *logrus.Entry
}

// NewInfoConfig create a new config.
func NewInfoConfig() *InfoConfig {
	return &InfoConfig{
		log: logrus.WithFields(logrus.Fields{
			"pkg": "etc",
			"view": "info",
		}),
	}
}

// Layout lay out the view of informations of Argo CUI.
func (c *InfoConfig) Layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(infoViewName, x0, y0, x1, y1)
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
				c.log.Errorf("couldn't get the context: %s.", err)
				context = ""
			}

			if namespace, err = argoutil.GetNamespace(); err != nil {
				c.log.Errorf("couldn't get the namespace: %s.", err)
				namespace = ""
			}

			v.Clear()
			fmt.Fprintln(v, "")
			fmt.Fprintf(v, "Context:     %s\n", color.ChangeColor(context, gocui.ColorYellow))
			fmt.Fprintf(v, "Namespace:   %s\n", color.ChangeColor(namespace, gocui.ColorYellow))
			fmt.Fprintf(v, "Argo Rev:    %s\n", color.ChangeColor(argoRev, gocui.ColorYellow))
			fmt.Fprintf(v, "Argocui Rev: %s\n", color.ChangeColor(argocuiRev, gocui.ColorYellow))
			fmt.Fprintf(v, "Homepage:    %s\n", color.ChangeColor(homePage, gocui.ColorYellow))
			return nil
		})
	}

	return nil
}
