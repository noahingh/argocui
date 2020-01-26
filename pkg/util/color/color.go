package color

import (
	"fmt"
	"regexp"

	"github.com/jroimartin/gocui"
)

// ChangeColor change the color of the word.
func ChangeColor(word string, color gocui.Attribute) string {
	w := fmt.Sprintf("\033[3%d;1m%s\033[0m", color-1, word)
	return w
}

// HasColor validate whether the word has color or not.
func HasColor(word string) bool {
	re := regexp.MustCompile("\033\\[3[0-9]+;1m.*\033\\[0m")
	cnt := len(re.FindAll([]byte(word), -1))
	if cnt == 0 {
		return false
	}
	return true
}
