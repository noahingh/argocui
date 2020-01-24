package view

import (
	"github.com/jroimartin/gocui"
)
// LineEditor make to edit only a single line, "enter" can't make a new line.
func LineEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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
		v.Clear()
		_, cy := v.Cursor()
		v.SetCursor(0, cy)
	}
}
