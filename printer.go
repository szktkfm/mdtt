package mdtt

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
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

// (?<=\|?\s*)-+
// ^\s*\|?\s*\-+

var tableDelimLeft = regexp.MustCompile(`^\s*\:\-+\s*$`)
var tableDelimRight = regexp.MustCompile(`^\s*\-+\:\s*$`)
var tableDelimCenter = regexp.MustCompile(`^\s*\:\-+\:\s*$`)
var tableDelimNone = regexp.MustCompile(`^\s*\-+\s*$`)
var tableDelim = regexp.MustCompile(`^\s*\|?\s*\-+`)

func (t *TableWriter) findSegment(fp io.Reader) {
	fmt.Println([]byte(fmt.Sprint(fp)))
	scanner := bufio.NewScanner(fp)

	var (
		prevlen  int
		prevline string
		pos      int
		inTable  bool
		start    int
		end      int
	)

	for scanner.Scan() {
		fmt.Println("byte: ", scanner.Bytes())
		l := scanner.Text()
		if inTable {
			if l == "" {
				inTable = false
				end = pos
			}
		}

		pos += len(l) + 1
		fmt.Println("pos: ", pos)
		fmt.Println("prevlen: ", prevlen)

		if tableDelimLeft.MatchString(l) ||
			tableDelimRight.MatchString(l) ||
			tableDelimCenter.MatchString(l) ||
			tableDelimNone.MatchString(l) ||
			tableDelim.MatchString(l) {
			// header check
			if len(strings.Split(trimPipe(prevline), "|")) <= len(strings.Split(trimPipe(l), "|")) {
				inTable = true
				start = pos - len(l) - prevlen
			} else {
				continue
			}

		}
		fmt.Println("start: ", start)

		prevline = l
		prevlen = len(l) + 1
	}
	if inTable {
		end = pos
	}

	// TODO: listで返す
	t.seg = TableSegment{Start: start - 1, Length: end - start}
}

func trimPipe(l string) string {
	if l[0] == '|' {
		l = l[1:]
	}
	if l[len(l)-1] == '|' {
		l = l[:len(l)-1]
	}
	return l
}

func (t *TableWriter) write(s []byte) {

}

type TableWriter struct {
	segs []TableSegment
	seg  TableSegment
}
type TableSegment struct {
	Start  int
	Length int
}
