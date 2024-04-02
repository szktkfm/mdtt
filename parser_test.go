package mdtt

import (
	"io"
	"os"
	"testing"
)

func TestParse(t *testing.T) {

	cols1 := []Column{
		{Title: NewCell("Key binding"), Width: 4},
		{Title: NewCell("Description"), Width: 20},
	}
	rows1 := []NaiveRow{
		{"`Arrows`, `hjkl`", "Move"},
	}
	want1 := NewTable(
		WithColumns(cols1),
		WithNaiveRows(rows1),
	)

	cols2 := []Column{
		{Title: NewCell("Key binding"), Width: 4},
		{Title: NewCell("Description"), Width: 20},
	}
	rows2 := []NaiveRow{
		{"**Esc**, _q_", "Exit"},
	}
	want2 := NewTable(
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
		if c.Title.Value() != y.cols[i].Title.Value() {
			ret = false
		}
	}

	for i, r := range x.rows {
		for j, c := range r {
			if c.Value() != y.rows[i][j].Value() {
				ret = false
			}
		}
	}
	return ret
}
