package color

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

// ChangeColor change the color of the word.
func ChangeColor(word string, color gocui.Attribute) string {
	w := fmt.Sprintf("\033[3%d;1m%s\033[0m", color-1, word)
	return w
}
