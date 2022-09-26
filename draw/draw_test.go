package draw

import (
	"drawgraph"
	"testing"
)

func TestDraw(t *testing.T) {
	type args struct {
		name   string
		outDir string
		nodes  []draw_graph.Node
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "example",
			args: args{
				name:   "example",
				outDir: "/tmp",
				nodes: []draw_graph.Node{
					{
						ID:   1,
						Text: "root",
					},
					{
						ID:     2,
						Text:   "a",
						Parent: 1,
					},
					{
						ID:     3,
						Text:   "b",
						Parent: 1,
					},
					{
						ID:     4,
						Text:   "c",
						Parent: 2,
					},
					{
						ID:   5,
						Text: "d",
					},
					{
						ID:     6,
						Text:   "e",
						Parent: 5,
					},
					{
						ID:     7,
						Text:   "f",
						Parent: 5,
					},
				},
			},
			wantErr: false,
		}, {
			name: "example-cn",
			args: args{
				name:   "example-cn",
				outDir: "/tmp",
				nodes: []draw_graph.Node{
					{
						ID:   1,
						Text: "根节点",
					},
					{
						ID:     2,
						Text:   "测试测试测试测试测试测试测试",
						Parent: 1,
					},
					{
						ID:     3,
						Text:   "测试",
						Parent: 1,
					},
					{
						ID:     4,
						Text:   "测试测试测试测试测试测试测试",
						Parent: 2,
					},
					{
						ID:   5,
						Text: "测试测试测试测试测试测试测试",
					},
					{
						ID:     6,
						Text:   "测试测试测试测试测试测试测试",
						Parent: 5,
					},
					{
						ID:     7,
						Text:   "测试测试测试测试测试测试测试",
						Parent: 5,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Draw(tt.args.name, tt.args.outDir, tt.args.nodes); (err != nil) != tt.wantErr {
				t.Errorf("Draw() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
