package mdtt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
)

func (t *TableWriter) render(m TableModel) {
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

	t.text = sb.String()
}

func Write(m Model) {
	tw := TableWriter{}
	tw.write(m)
}

func (t *TableWriter) write(m Model) {

	t.render(m.table)

	if m.inplace {
		fp, err := os.Open(m.fpath)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()
		t.findSegment(fp)
		fp.Seek(0, 0)
		b, _ := io.ReadAll(fp)
		b = append(b[:t.seg.start-1],
			append([]byte(t.text), b[t.seg.end:]...)...)

		os.WriteFile(m.fpath, b, 0644)
	} else {
		fmt.Print(t.text)
	}
}

// (?<=\|?\s*)-+
// ^\s*\|?\s*\-+

// TODO: delimeterの左寄せとか
var tableDelimLeft = regexp.MustCompile(`^\s*\:\-+\s*$`)
var tableDelimRight = regexp.MustCompile(`^\s*\-+\:\s*$`)
var tableDelimCenter = regexp.MustCompile(`^\s*\:\-+\:\s*$`)
var tableDelimNone = regexp.MustCompile(`^\s*\-+\s*$`)
var tableDelim = regexp.MustCompile(`^\s*\|?\s*\-+`)

func (t *TableWriter) findSegment(fp io.Reader) {
	// fmt.Println([]byte(fmt.Sprint(fp)))
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
		l := scanner.Text()
		if inTable {
			if l == "" {
				inTable = false
				end = pos
			}
		}

		pos += len(l) + 1

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
		prevline = l
		prevlen = len(l) + 1
	}
	if inTable {
		end = pos
	}

	// TODO: listで返す
	t.seg = TableSegment{start: start, end: end}
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

type TableWriter struct {
	segs []TableSegment
	seg  TableSegment
	text string
}
type TableSegment struct {
	start int
	end   int
}
