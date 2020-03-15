package ui

import (
	"fmt"

	"github.com/hanjunlee/argocui/internal/config"
	"github.com/hanjunlee/argocui/pkg/util/argo"
	"github.com/hanjunlee/argocui/pkg/util/color"

	"github.com/jroimartin/gocui"
)

func (m *Manager) layoutHelper(v *gocui.View) {
	const (
		argoRev    = config.ArgoVersion
		argocuiRev = config.Version
		homePage   = config.HomePage
	)
	var (
		context string
	)

	context, _ = argo.GetCurrentContext()

	fmt.Fprintln(v, "")
	fmt.Fprintf(v, "Context:      %s\n", color.ChangeColor(context, gocui.ColorYellow))
	fmt.Fprintf(v, "Namespace:    %s\n", color.ChangeColor(m.namespace, gocui.ColorYellow))
	fmt.Fprintf(v, "Argo Rev:     %s\n", color.ChangeColor(argoRev, gocui.ColorYellow))
	fmt.Fprintf(v, "Argocui Rev:  %s\n", color.ChangeColor(argocuiRev, gocui.ColorYellow))
	fmt.Fprintf(v, "Homepage:     %s\n", color.ChangeColor(homePage, gocui.ColorYellow))
}
