package color

import (
	"testing"

	"github.com/jroimartin/gocui"
)

func TestHasColor(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "no color",
			args: args{
				word: "no color",
			},
			want: false,
		},
		{
			name: "hello world",
			args: args{
				word: ChangeColor("hello world", gocui.ColorYellow),
			},
			want: true,
		},
		{
			name: "digit",
			args: args{
				word: ChangeColor("01234", gocui.ColorYellow),
			},
			want: true,
		},
		{
			name: "uni-code",
			args: args{
				word: ChangeColor("한국어", gocui.ColorYellow),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Logf("word: %s", tt.args.word)
			if got := HasColor(tt.args.word); got != tt.want {
				t.Errorf("HasColor() = %v, want %v", got, tt.want)
			}
		})
	}
}
