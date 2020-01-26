package tablewriter

import (
	"fmt"
	"io"
	"strings"

	"github.com/hanjunlee/argocui/pkg/util/color"

	// "github.com/willf/pad"
	padUtf8 "github.com/willf/pad/utf8"
)

// TableWriter render the data as table-like format.
type TableWriter struct {
	writer io.Writer

	// columns of the header
	columns      []string
	columnWidths []int
	// the border of the header
	headerBorder bool

	data [][]string
}

// NewTableWriter return a new table writer.
func NewTableWriter(w io.Writer) *TableWriter {
	return &TableWriter{
		writer: w,
		data:   [][]string{},
	}
}

// SetColumns set columns of table at the header.
func (t *TableWriter) SetColumns(cols []string) {
	t.columns = cols
}

// SetColumnWidths set the width of columns.
// If the size of value is over the width, it mosaic with "...".
func (t *TableWriter) SetColumnWidths(widths []int) {
	t.columnWidths = widths
}

// SetHeaderBorder set the border of the header.
func (t *TableWriter) SetHeaderBorder(ok bool) {
	t.headerBorder = ok
}

// Append append the row.
func (t *TableWriter) Append(row []string) {
	t.data = append(t.data, row)
}

// AppendBulk append the rows.
func (t *TableWriter) AppendBulk(rows [][]string) {
	for _, r := range rows {
		t.data = append(t.data, r)
	}
}

// Render write the table.
func (t *TableWriter) Render() error {
	widths := t.getColumnWidths()

	// header.
	if t.columns != nil {
		l := getLine(t.columns, widths)
		fmt.Fprintln(t.writer, l)
	}

	// border.
	if t.headerBorder {
		fmt.Fprintln(t.writer, strings.Repeat("─", sum(widths...)))
	}

	for _, r := range t.data {
		l := getLine(r, widths)
		fmt.Fprintln(t.writer, l)
	}

	return nil
}

// it return the set of max widths of column, but it is overrided if the column widths attribute exist.
func (t *TableWriter) getColumnWidths() []int {
	cc := t.getColumnCnt()
	maxWidths := make([]int, cc)

	for _, row := range t.data {
		for i, word := range row {
			if mw := maxWidths[i]; len(word) > mw {
				maxWidths[i] = len(word)
			}
		}
	}

	if t.columnWidths == nil {
		return maxWidths
	}

	// override the width of column from the attribute.
	for i := range maxWidths {
		if i >= len(t.columnWidths) {
			break
		}

		maxWidths[i] = t.columnWidths[i]
	}
	return maxWidths
}

// it return the max count of column.
func (t *TableWriter) getColumnCnt() int {
	ret := 0

	if t.columns != nil {
		ret = len(t.columns)
	}

	for _, row := range t.data {
		if len(row) > ret {
			ret = len(row)
		}
	}
	return ret
}

func getLine(row []string, widths []int) string {
	const (
		widthColor = 11
	)
	var (
		ret = ""
	)

	for i, word := range row {
		width := widths[i]
		// complement the width for a color.
		if color.HasColor(word) {
			width = width + widthColor
		}

		word = dotdotdot(word, width)
		ret = ret + padUtf8.Right(word, width, " ")
	}
	return ret
}

func dotdotdot(word string, size int) string {
	if len(word) > size {
		word = word[:size-3] + "..."
	}
	return word
}

func sum(ints ...int) int {
	ret := 0
	for _, i := range ints {
		ret = ret + i
	}
	return ret
}
