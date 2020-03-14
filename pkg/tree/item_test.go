package tree

import (
	"fmt"
	"testing"
	"time"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/hanjunlee/tree"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_item_String(t *testing.T) {
	tests := []struct {
		name string
		i    item
		want string
	}{
		{
			name: "retry with only one child",
			i: item{
				Phase:       wf.NodeSucceeded,
				DisplayName: "──retry(0)",
			},
			want: fmt.Sprintf("%s%s %s", tabRetrySingleChild, icons[wf.NodeSucceeded], "retry(0)"),
		},
		{
			name: "step group with only one child",
			i: item{
				Phase:       wf.NodeSucceeded,
				DisplayName: "·─coinflip",
			},
			want: fmt.Sprintf("%s%s %s", tabStepGroupMultiChild, icons[wf.NodeSucceeded], "coinflip"),
		},
		{
			name: "template name",
			i: item{
				Phase:        wf.NodeSucceeded,
				DisplayName:  "coinflip",
				TemplateName: "coinflip",
			},
			want: fmt.Sprintf("%s %s (%s)", icons[wf.NodeSucceeded], "coinflip", "coinflip"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.String(); got != tt.want {
				t.Errorf("item.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_item_Less(t *testing.T) {
	type args struct {
		comp tree.Item
	}
	tests := []struct {
		name string
		i    item
		args args
		want bool
	}{
		{
			name: "comparing by the start time",
			i: item{
				DisplayName: "A",
				StartedAt:   metav1.NewTime(time.Now().Add(-time.Second)),
				Type:        wf.NodeTypePod,
			},
			args: args{
				comp: item{
					DisplayName: "B",
					StartedAt:   metav1.NewTime(time.Now()),
					Type:        wf.NodeTypePod,
				},
			},
			want: true,
		},
		{
			name: "comparing by the display name",
			i: item{
				DisplayName: "A",
				StartedAt:   metav1.NewTime(time.Now()),
				Type:        wf.NodeTypePod,
			},
			args: args{
				comp: item{
					DisplayName: "B",
					StartedAt:   metav1.NewTime(time.Now()),
					Type:        wf.NodeTypePod,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.i.Less(tt.args.comp); got != tt.want {
				t.Errorf("item.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}
