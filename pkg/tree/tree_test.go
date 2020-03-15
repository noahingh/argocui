package tree

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	wf "github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/davecgh/go-spew/spew"
	"github.com/hanjunlee/tree"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace = "argo"
)

func TestClient_GetTreeRoot(t *testing.T) {
	type args struct {
		w *wf.Workflow
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		{
			name: "hello world",
			args: args{
				w: &wf.Workflow{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "hello-world",
						Namespace: namespace,
					},
					Status: wf.WorkflowStatus{
						Nodes: map[string]wf.NodeStatus{
							"hello-world": wf.NodeStatus{
								ID:          "hello-world",
								DisplayName: "hello-world",
								Type:        wf.NodeTypePod,
								Phase:       wf.NodeSucceeded,
								Message:     "hello world",
							},
						},
					},
				},
			},
			want: [][]string{
				[]string{fmt.Sprintf("%s %s", icons[wf.NodeSucceeded], "hello-world"), "hello-world", "0s", "hello world"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTreeRoot(tt.args.w)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetTreeNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetTreeNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initTree(t *testing.T) {
	type args struct {
		items    map[string]item
		rootName string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			// https://github.com/argoproj/argo/blob/master/examples/hello-world.yaml
			name: "hello world",
			args: args{
				items: map[string]item{
					"hello-world": item{
						DisplayName: "hello-world",
						Type:        wf.NodeTypePod,
						Phase:       wf.NodeSucceeded,
						StartedAt:   metav1.NewTime(time.Now()),
					},
				},
				rootName: "hello-world",
			},
			want: []string{
				fmt.Sprintf("%s %s", icons[wf.NodeSucceeded], "hello-world"),
			},
			wantErr: false,
		},
		{
			// https://github.com/argoproj/argo/blob/master/examples/retry-container.yaml
			name: "retry",
			args: args{
				items: map[string]item{
					"retry": item{
						DisplayName: "retry",
						Type:        wf.NodeTypeRetry,
						Phase:       wf.NodeSucceeded,
						StartedAt:   metav1.NewTime(time.Now()),
						Children: []string{
							"retry-0",
							"retry-1",
							"retry-2",
						},
					},
					"retry-2": item{
						DisplayName: "retry(2)",
						Type:        wf.NodeTypePod,
						Phase:       wf.NodeSucceeded,
						StartedAt:   metav1.NewTime(time.Now().Add(3 * time.Second)),
					},
					"retry-1": item{
						DisplayName: "retry(1)",
						Type:        wf.NodeTypePod,
						Phase:       wf.NodeFailed,
						StartedAt:   metav1.NewTime(time.Now().Add(2 * time.Second)),
					},
					"retry-0": item{
						DisplayName: "retry(0)",
						Type:        wf.NodeTypePod,
						Phase:       wf.NodeFailed,
						StartedAt:   metav1.NewTime(time.Now().Add(1 * time.Second)),
					},
				},
				rootName: "retry",
			},
			want: []string{
				fmt.Sprintf("%s %s", icons[wf.NodeSucceeded], "retry"),
				fmt.Sprintf("├─%s %s", icons[wf.NodeFailed], "retry(0)"),
				fmt.Sprintf("├─%s %s", icons[wf.NodeFailed], "retry(1)"),
				fmt.Sprintf("└─%s %s", icons[wf.NodeSucceeded], "retry(2)"),
			},
			wantErr: false,
		},
		{
			// https://github.com/argoproj/argo/blob/master/examples/coinflip.yaml
			name: "coinflip",
			args: args{
				items: map[string]item{
					"coinflip": item{
						DisplayName: "coinflip",
						Type:        wf.NodeTypeSteps,
						Phase:       wf.NodeSucceeded,
						StartedAt:   metav1.NewTime(time.Now()),
						Children: []string{
							"coinflip-1",
						},
					},
					"coinflip-1": item{
						DisplayName: "[0]",
						Type:        wf.NodeTypeStepGroup,
						Phase:       wf.NodeSucceeded,
						StartedAt:   metav1.NewTime(time.Now()),
						BoundaryID:  "coinflip",
						Children: []string{
							"coinflip-2",
						},
					},
					"coinflip-2": item{
						DisplayName: "flip-coin",
						Type:        wf.NodeTypePod,
						Phase:       wf.NodeSucceeded,
						StartedAt:   metav1.NewTime(time.Now()),
						BoundaryID:  "coinflip",
						Children: []string{
							"coinflip-3",
						},
					},
					"coinflip-3": item{
						DisplayName: "[1]",
						Type:        wf.NodeTypeStepGroup,
						Phase:       wf.NodeSucceeded,
						StartedAt:   metav1.NewTime(time.Now().Add(1 * time.Second)),
						BoundaryID:  "coinflip",
						Children: []string{
							"coinflip-4",
							"coinflip-5",
						},
					},
					"coinflip-4": item{
						DisplayName: "heads",
						Type:        wf.NodeTypePod,
						Phase:       wf.NodeSucceeded,
						StartedAt:   metav1.NewTime(time.Now().Add(1 * time.Second)),
						BoundaryID:  "coinflip",
					},
					"coinflip-5": item{
						DisplayName: "tails",
						Type:        wf.NodeTypePod,
						Phase:       wf.NodeSucceeded,
						StartedAt:   metav1.NewTime(time.Now().Add(1 * time.Second)),
						BoundaryID:  "coinflip",
					},
				},
				rootName: "coinflip",
			},
			want: []string{
				fmt.Sprintf("%s %s", icons[wf.NodeSucceeded], "coinflip"),
				fmt.Sprintf("├───%s %s", icons[wf.NodeSucceeded], "flip-coin"),
				fmt.Sprintf("└─%s", icons[wf.NodeSucceeded]),
				fmt.Sprintf("  ├─%s %s", icons[wf.NodeSucceeded], "heads"),
				fmt.Sprintf("  └─%s %s", icons[wf.NodeSucceeded], "tails"),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := tt.args.items[tt.args.rootName]
			tree := tree.NewTree(root)
			err := initTree(tree, tt.args.items, root)

			if (err != nil) != tt.wantErr {
				t.Errorf("getTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			got := tree.Render()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getTree() got = %s, want %s", spew.Sdump(got), spew.Sdump(tt.want))
			}
		})
	}
}
