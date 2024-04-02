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

		b := tw.replaceTable(fp, m.choose)
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
		sb.WriteString(PadOrTruncate(c.Title.Value(), c.Width-1))
		width += c.Width
	}
	sb.WriteString("|\n")

	for _, c := range m.cols {
		sb.WriteString("| ")
		sb.WriteString(strings.Repeat("-", c.Width-2))
		sb.WriteString(" ")
	}
	sb.WriteString("|\n")

	for _, row := range m.rows {
		for i, c := range row {
			sb.WriteString("| ")
			sb.WriteString(PadOrTruncate(c.Value(), m.cols[i].Width-1))
		}
		sb.WriteString("|\n")
	}

	t.text = sb.String()
}

func (t *TableWriter) replaceTable(fp *os.File, idx int) []byte {
	t.findSegment(fp)
	fp.Seek(0, 0)
	b, _ := io.ReadAll(fp)
	b = append(b[:t.segs[idx].Start-1],
		append([]byte(t.text), b[t.segs[idx].End:]...)...)
	return b
}

var (
	tableDelimLeft   = regexp.MustCompile(`^\s*\:\-+\s*$`)
	tableDelimRight  = regexp.MustCompile(`^\s*\-+\:\s*$`)
	tableDelimCenter = regexp.MustCompile(`^\s*\:\-+\:\s*$`)
	tableDelimNone   = regexp.MustCompile(`^\s*\-+\s*$`)
	thematicBreak    = regexp.MustCompile(`^\s{0,3}((-\s*){3,}|(\*\s*){3,}|(_\s*){3,})\s*$`)
	prefixSpace      = regexp.MustCompile(`^\s{0,3}`)
	// prefixSpace      = regexp.MustCompile(`^\s{0,3}`)
	fencedCodeBlock = regexp.MustCompile("^```|~~~.*$")
)

func (t *TableWriter) findSegment(fp io.Reader) {
	scanner := bufio.NewScanner(fp)

	var (
		prevlen  int
		prevline string
		pos      int
		start    int
		end      int
	)

	var (
		inTable     bool
		inCodeBlock bool
	)

	var segs []TableSegment
	for scanner.Scan() {
		l := scanner.Text()

		// TODO: ```で始まって ~~~で閉じている場合
		if isCodeFence(l) {
			inCodeBlock = !inCodeBlock
		}

		if inTable {
			if isNewLine(l) || isThematicBreak(l) {
				inTable = false
				end = pos
				segs = append(segs, TableSegment{Start: start, End: pos})
			}
		}

		if inCodeBlock {
			pos += len(l) + 1
			prevline = l
			prevlen = len(l) + 1
			continue
		}

		pos += len(l) + 1
		if isTableDelimiter(l) {
			// header check
			if isNewLine(prevline) {
				continue
			}

			// spaceがprefixに四つ以上ある場合はfalse
			if isTableHeader(prevline, l) {
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
		segs = append(segs, TableSegment{Start: start, End: pos})
		end = pos
	}

	// TODO: listで返す
	log.Debug(t)
	t.seg = TableSegment{Start: start, End: end}
	t.segs = segs
	log.Debugf("start: %d, end: %d", start, end)
}

func isTableHeader(header string, delim string) bool {
	return len(strings.Split(trimPipe(header), "|")) <= len(strings.Split(trimPipe(delim), "|"))
}

func isCodeFence(s string) bool {
	return fencedCodeBlock.MatchString(s)
}

func isThematicBreak(s string) bool {
	return thematicBreak.MatchString(s)
}

func isNewLine(s string) bool {
	return s == ""
}
func isTableDelimiter(s string) bool {
	// TODO: prefixのスペース問題
	// スペースが4つ以上の場合は捨てる
	delim, _, _ := strings.Cut(
		trimPipe(prefixSpace.ReplaceAllString(s, "")), "|")

	if tableDelimLeft.MatchString(delim) ||
		tableDelimRight.MatchString(delim) ||
		tableDelimCenter.MatchString(delim) {
		return true
	}
	if tableDelimNone.MatchString(delim) && strings.Contains(s, "|") {
		return true
	}

	return false
}

func trimPipe(l string) string {
	if len(l) == 0 {
		return l
	}
	// spaceをtrimする
	if l[0] == '|' {
		l = l[1:]
	}
	if l[len(l)-1] == '|' {
		l = l[:len(l)-1]
	}
	return l
}
