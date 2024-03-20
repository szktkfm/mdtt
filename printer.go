package mdtt

import (
	"fmt"
	"strings"
)

func print(m TableModel) {
	var sb strings.Builder
	var width int

	for _, c := range m.cols {
		sb.WriteString("|")
		sb.WriteString(PadOrTruncate(c.Title.Value(), c.Width))
		width += c.Width
	}
	sb.WriteString("|\n")

	for _, c := range m.cols {
		sb.WriteString("|")
		sb.WriteString(strings.Repeat("-", c.Width))
	}
	sb.WriteString("|\n")

	for _, row := range m.rows {
		for i, c := range row {
			sb.WriteString("|")
			sb.WriteString(PadOrTruncate(c.Value(), m.cols[i].Width))
		}
		sb.WriteString("|\n")
	}
	sb.WriteString("\n")

	fmt.Print(sb.String())
}
