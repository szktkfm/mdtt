package mdtt

import (
	"io"
	"os"
	"testing"
)

func TestParse(t *testing.T) {

	cols1 := []column{
		{title: NewCell("Key binding"), width: 4},
		{title: NewCell("Description"), width: 20},
	}
	rows1 := []naiveRow{
		{"`Arrows`, `hjkl`", "Move"},
	}
	want1 := NewTableModel(
		WithColumns(cols1),
		WithNaiveRows(rows1),
	)

	cols2 := []column{
		{title: NewCell("Key binding"), width: 4},
		{title: NewCell("Description"), width: 20},
	}
	rows2 := []naiveRow{
		{"**Esc**, _q_", "Exit"},
	}
	want2 := NewTableModel(
		WithColumns(cols2),
		WithNaiveRows(rows2),
	)

	f, _ := os.Open("testdata/01.md")
	defer f.Close()
	s, _ := io.ReadAll(f)
	got := parse(s)

	if !isEqualTables(want1, got[0]) {
		t.Error("Table value is mismatch")
	}

	if !isEqualTables(want2, got[1]) {
		t.Error("Table value is mismatch")
	}
}

func isEqualTables(x, y TableModel) bool {
	ret := true
	for i, c := range x.cols {
		if c.title.value() != y.cols[i].title.value() {
			ret = false
		}
	}

	for i, r := range x.rows {
		for j, c := range r {
			if c.value() != y.rows[i][j].value() {
				ret = false
			}
		}
	}
	return ret
}
