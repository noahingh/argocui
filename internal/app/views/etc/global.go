package etc 

import (
	"github.com/jroimartin/gocui"
)

// GlobalKeybinding is keybinding for gui.
func GlobalKeybinding(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, 
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit	
		}); err != nil {
		return err
	}
	return nil
}
