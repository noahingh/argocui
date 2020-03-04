package ui

import (
	"github.com/jroimartin/gocui"
)

func inlineEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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
		for true {
			line, _ := v.Line(0)
			if len(line) == 0 {
				break
			}

			v.EditDelete(true)
		}
	}
}
