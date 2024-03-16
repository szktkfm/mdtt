package mdtt

import (
	"testing"
)

func TestParse(t *testing.T) {

	columns := []Column{
		{Title: NewCell("Key binding"), Width: 4},
		{Title: NewCell("Description"), Width: 20},
	}
	rows := []NaiveRow{
		{"`Arrows`, `hjkl`", "Move"},
		{"**Esc**, _q_", "Exit"},
	}
	want := New(
		WithColumns(columns),
		WithNaiveRows(rows),
	)

	got := parse("testdata/01.md")

	if !isEqualTables(want, got.table) {
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
