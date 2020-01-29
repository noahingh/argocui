package tablewriter

import (
	"regexp"
)

// cntOfColor return how many words has color.
func cntOfColor(word string) int {
	re := regexp.MustCompile("\033\\[3[0-9]+;1m([0-9]|[a-z]|[A-Z])*\033\\[0m")
	return len(re.FindAll([]byte(word), -1))
}
