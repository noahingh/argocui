package tablewriter

import (
	"regexp"
	"strings"
)

const (
	// colorRegexp = "\033\\[3[0-9]+;1m([0-9]|[a-z]|[A-Z]|◷|●|✔|○|✖|⚠)*\033\\[0m"
	colorRegexp = "\033\\[3[0-9]+;1m.+\033\\[0m"
)

type coloredWord struct {
	prefix  string
	word    string
	postfix string
}

func splitByColoredWord(s string) []string {
	ret := make([]string, 0)

	re := regexp.MustCompile(colorRegexp)
	indice := re.FindAllStringIndex(s, -1)
	if len(indice) == 0 {
		ret = append(ret, s)
		return ret
	}

	prev := 0
	for _, i := range indice {
		begin, end := i[0], i[1]
		if prev != begin {
			ret = append(ret, s[prev:begin])
		}

		ret = append(ret, s[begin:end])
		prev = end
	}
	if prev != len(s) {
		ret = append(ret, s[prev:len(s)])
	}

	return ret
}

func splitToColoredWords(s string) []coloredWord {
	const (
		postfix = "\033[0m"
	)
	var (
		ret = make([]coloredWord, 0)
	)

	words := splitByColoredWord(s)
	for _, w := range words {
		if w == "" {
			continue
		}

		re := regexp.MustCompile("\033\\[3[0-9]+;1m")

		// colored word.
		if prefix := re.FindString(w); prefix != "" {
			w = strings.TrimLeft(w, prefix)
			w = strings.TrimRight(w, postfix)

			ret = append(ret, coloredWord{
				prefix:  prefix,
				word:    w,
				postfix: postfix,
			})
		} else {
			ret = append(ret, coloredWord{
				word: w,
			})
		}
	}
	return ret
}

func cutWord(s string, limit int) string {
	ws := splitToColoredWords(s)

	size := 0
	for _, w := range ws {
		size = size + len(w.word)
	}

	removed := size - limit
	for i := len(ws) - 1; i >= 0; i-- {
		if removed <= 0 {
			break
		}

		w := ws[i]

		if len(w.word) <= removed {
			ws = ws[:i]
			removed = removed - len(w.word)
		} else {
			w.word = w.word[:len(w.word)-removed]
			w.word = dotdotdot(w.word)
			ws[i] = w
			removed = 0
		}
	}

	var ret string
	for _, w := range ws {
		ret = strings.Join([]string{ret, w.prefix, w.word, w.postfix}, "")
	}
	return ret
}

func dotdotdot(s string) string {
	if len(s) > 3 {
		s = s[:len(s)-3] + "..."
	}
	return s
}

// cntOfColor return how many words has color.
func cntOfColor(word string) int {
	re := regexp.MustCompile(colorRegexp)
	return len(re.FindAllString(word, -1))
}
