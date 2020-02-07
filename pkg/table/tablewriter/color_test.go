package tablewriter

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/hanjunlee/argocui/pkg/util/color"
	"github.com/jroimartin/gocui"
)

func getColoerdWord(word string, c gocui.Attribute) coloredWord {
	return coloredWord{
		prefix:  fmt.Sprintf("\033[3%d;1m", c-1),
		word:    word,
		postfix: "\033[0m",
	}
}

func Test_splitByColoredWord(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no colored words",
			args: args{
				s: "foo bar",
			},
			want: []string{
				"foo bar",
			},
		},
		{
			name: "mixed 1",
			args: args{
				s: color.ChangeColor("foo", gocui.ColorYellow) + "bar",
			},
			want: []string{
				color.ChangeColor("foo", gocui.ColorYellow), 
				"bar",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitByColoredWord(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitByColoredWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitToColoredWords(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []coloredWord
	}{
		{
			name: "no colored words",
			args: args{
				s: "foo bar",
			},
			want: []coloredWord{
				coloredWord{word: "foo bar"},
			},
		},
		{
			name: "mixed 1",
			args: args{
				s: color.ChangeColor("foo", gocui.ColorYellow) + "bar",
			},
			want: []coloredWord{
				getColoerdWord("foo", gocui.ColorYellow),
				coloredWord{word: "bar"},
			},
		},
		{
			name: "mixed 2",
			args: args{
				s: "foo" + color.ChangeColor("bar", gocui.ColorYellow),
			},
			want: []coloredWord{
				coloredWord{word: "foo"},
				getColoerdWord("bar", gocui.ColorYellow),
			},
		},
		{
			name: "colored words",
			args: args{
				s: color.ChangeColor("foo", gocui.ColorYellow) + color.ChangeColor("bar", gocui.ColorGreen),
			},
			want: []coloredWord{
				getColoerdWord("foo", gocui.ColorYellow),
				getColoerdWord("bar", gocui.ColorGreen),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitToColoredWords(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitToColoredWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cutWord(t *testing.T) {
	type args struct {
		s     string
		limit int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no color",
			args: args{
				s:     "foo bar",
				limit: 5,
			},
			want: "fo...",
		},
		{
			name: "mixed 1",
			args: args{
				s:     color.ChangeColor("foo", gocui.ColorYellow) + " bar",
				limit: 3,
			},
			want: color.ChangeColor("foo", gocui.ColorYellow),
		},
		{
			name: "mixed 2",
			args: args{
				s:     color.ChangeColor("foo", gocui.ColorYellow) + "barbaz",
				limit: 5,
			},
			want: color.ChangeColor("foo", gocui.ColorYellow) + "ba",
		},
		{
			name: "mixed 3",
			args: args{
				s:     color.ChangeColor("foo", gocui.ColorYellow) + "barbaz",
				limit: 7,
			},
			want: color.ChangeColor("foo", gocui.ColorYellow) + "b...",
		},
		{
			name: "color",
			args: args{
				s:     color.ChangeColor("foobarbaz", gocui.ColorYellow),
				limit: 7,
			},
			want: color.ChangeColor("foob...", gocui.ColorYellow),
		},
		{
			name: "two color",
			args: args{
				s:     color.ChangeColor("foobarbaz", gocui.ColorYellow) + color.ChangeColor("foobarbaz", gocui.ColorRed),
				limit: 7,
			},
			want: color.ChangeColor("foob...", gocui.ColorYellow),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cutWord(tt.args.s, tt.args.limit); got != tt.want {
				t.Errorf("dotdotdot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cntOfColor(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "no color",
			args: args{
				word: "nothing",
			},
			want: 0,
		},
		{
			name: "a single color",
			args: args{
				word: color.ChangeColor("yellow", gocui.ColorYellow),
			},
			want: 1,
		},
		{
			name: "many colors",
			args: args{
				word: strings.Join([]string{
					color.ChangeColor("red", gocui.ColorRed),
					color.ChangeColor("blue", gocui.ColorBlue),
					color.ChangeColor("gree", gocui.ColorGreen),
				}, " "),
			},
			want: 3,
		},
		{
			name: "special char",
			args: args{
				word: strings.Join([]string{
					color.ChangeColor("◷", gocui.ColorRed),
					color.ChangeColor("○", gocui.ColorRed),
				}, " "),
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cntOfColor(tt.args.word); got != tt.want {
				t.Logf("word = %v", tt.args.word)
				t.Errorf("getCntOfColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
