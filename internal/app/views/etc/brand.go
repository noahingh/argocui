package etc

import (
	"fmt"

	"github.com/hanjunlee/argocui/pkg/util/color"
	"github.com/jroimartin/gocui"
	"github.com/sirupsen/logrus"
)

const (
	brandViewName = "brand"
)

// BrandConfig is the configuration of a info view.
type BrandConfig struct {
}

// NewBrandConfig create a new config.
func NewBrandConfig() *BrandConfig {
	return &BrandConfig{}
}

// Layout lay out the view of brand.
func (b *BrandConfig) Layout(g *gocui.Gui, x0, y0, x1, y1 int) error {
	v, err := g.SetView(brandViewName, x0, y0, x1, y1)
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
