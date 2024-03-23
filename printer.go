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

type TableWriter struct {
	// segをフィールドに持つ必要ないのでは
	segs []TableSegment
	seg  TableSegment
	text string
}
type TableSegment struct {
	Start int
	End   int
}

func Write(m Model) {
	tw := TableWriter{}
	tw.render(m.table)
	if m.inplace {
		fp, err := os.Open(m.fpath)
		if err != nil {
			log.Fatal(err)
		}
		defer fp.Close()

		b := tw.replaceTable(fp)
		os.WriteFile(m.fpath, b, 0644)

	} else {
		fmt.Print(tw.text)
	}
}

func (t *TableWriter) render(m TableModel) {
	var sb strings.Builder
	var width int

	for _, c := range m.cols {
		sb.WriteString("| ")
		sb.WriteString(PadOrTruncate(c.Title.Value(), c.Width))
		width += c.Width
	}
	sb.WriteString(" |\n")

	for _, c := range m.cols {
		sb.WriteString("| ")
		sb.WriteString(strings.Repeat("-", c.Width))
	}
	sb.WriteString(" |\n")

	for _, row := range m.rows {
		for i, c := range row {
			sb.WriteString("| ")
			sb.WriteString(PadOrTruncate(c.Value(), m.cols[i].Width))
		}
		sb.WriteString(" |\n")
	}

	t.text = sb.String()
}

func (t *TableWriter) replaceTable(fp *os.File) []byte {
	t.findSegment(fp)
	fp.Seek(0, 0)
	b, _ := io.ReadAll(fp)
	b = append(b[:t.seg.Start-1],
		append([]byte(t.text), b[t.seg.End:]...)...)
	return b
}

// (?<=\|?\s*)-+
// ^\s*\|?\s*\-+

// TODO: delimeterの左寄せとか
// var tableDelimLeft = regexp.MustCompile(`^\s*\:\-+\s*$`)
// var tableDelimRight = regexp.MustCompile(`^\s*\-+\:\s*$`)
// var tableDelimCenter = regexp.MustCompile(`^\s*\:\-+\:\s*$`)
// var tableDelimNone = regexp.MustCompile(`^\s*\-+\s*$`)
var tableDelim = regexp.MustCompile(`^\s*\|?\s*\-+\s*\|?\s*`)

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
		if tableDelim.MatchString(l) {
			// header check
			log.Debug("line", l)
			if prevline == "" {
				continue
			}
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
	log.Debug(t)
	t.seg = TableSegment{Start: start, End: end}
	log.Debugf("start: %d, end: %d", start, end)
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
