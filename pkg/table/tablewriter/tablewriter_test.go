package tablewriter

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

type twriter struct {
	buffer []string
}

func (w *twriter) Write(bs []byte) (int, error) {
	w.buffer = append(w.buffer, string(bs))
	return len(bs), nil
}

func TestTableWriter_Render(t *testing.T) {
	type fields struct {
		writer       io.Writer
		columns      []string
		columnWidths []int
		headerBorder bool
		data         [][]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "without columns",
			fields: fields{
				writer: &twriter{},
				data: [][]string{
					[]string{"aaaaa", "bbbbb", "ccc"},
				},
			},
			want: []string{
				fmt.Sprintln("aaaaa" + "bbbbb" + "ccc"),
			},
			wantErr: false,
		},
		{
			name: "with columns",
			fields: fields{
				writer: &twriter{},
				columns: []string{"aaaaa", "bbbbb", "ccccc"},
				headerBorder: true,
				data: [][]string{
					[]string{"aaaaa", "bbbbb", "ccccc"},
				},
			},
			want: []string{
				fmt.Sprintln("aaaaa" + "bbbbb" + "ccccc"),
				fmt.Sprintln("─────" + "─────" + "─────"),
				fmt.Sprintln("aaaaa" + "bbbbb" + "ccccc"),
			},
			wantErr: false,
		},
		{
			name: "override width",
			fields: fields{
				writer: &twriter{},
				columns: []string{"aaaa", "bbbb", "cccc"},
				columnWidths: []int{4, 4, 4},
				headerBorder: true,
				data: [][]string{
					[]string{"aaaaa", "bbbbb", "ccccc"},
				},
			},
			want: []string{
				fmt.Sprintln("aaaa" + "bbbb" + "cccc"),
				fmt.Sprintln("────" + "────" + "────"),
				fmt.Sprintln("a..." + "b..." + "c..."),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TableWriter{
				writer:       tt.fields.writer,
				columns:      tt.fields.columns,
				columnWidths: tt.fields.columnWidths,
				headerBorder: tt.fields.headerBorder,
				data:         tt.fields.data,
			}
			if err := tr.Render(); (err != nil) != tt.wantErr {
				t.Errorf("TableWriter.Render() error = %v, wantErr %v", err, tt.wantErr)
			}
			w := tr.writer.(*twriter)
			if got := w.buffer; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableWriter.Render() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTableWriter_getColumnWidths(t *testing.T) {
	type fields struct {
		writer       io.Writer
		columns      []string
		columnWidths []int
		headerBorder bool
		data         [][]string
	}
	tests := []struct {
		name   string
		fields fields
		want   []int
	}{
		{
			name: "without override",
			fields: fields{
				writer: &twriter{},
				data: [][]string{
					[]string{"aaaaa", "bbbbb", "ccccc"},
					[]string{"aaaaa", "bbbb", "ccccc"},
				},
			},
			want: []int{5, 5, 5},
		},
		{
			name: "with override",
			fields: fields{
				writer:       &twriter{},
				columnWidths: []int{3, 3, 3},
				data: [][]string{
					[]string{"aaaaa", "bbbbb", "ccccc"},
					[]string{"aaaaa", "bbbb", "ccccc"},
				},
			},
			want: []int{3, 3, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TableWriter{
				writer:       tt.fields.writer,
				columns:      tt.fields.columns,
				columnWidths: tt.fields.columnWidths,
				headerBorder: tt.fields.headerBorder,
				data:         tt.fields.data,
			}
			if got := tr.getColumnWidths(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TableWriter.getColumnWidths() = %v, want %v", got, tt.want)
			}
		})
	}
}
