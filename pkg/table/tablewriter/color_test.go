package tablewriter

import (
	"strings"
	"testing"

	"github.com/hanjunlee/argocui/pkg/util/color"
	"github.com/jroimartin/gocui"
)

func Test_cntOfColor(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
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
