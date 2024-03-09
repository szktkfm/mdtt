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
	text string
}

type tableLocator struct {
	locs        []TableLocation
	inTable     bool
	codeFence   string
	inCodeBlock bool
}

type TableLocation struct {
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

	// render header
	for _, c := range m.cols {
		sb.WriteString("| ")
		sb.WriteString(padOrTruncate(c.title.value(), c.width-1))
		width += c.width
	}
	sb.WriteString("|\n")

	// render delimiter
	for _, c := range m.cols {
		log.Debug("column", c.title.value(), c.alignment)

		if c.alignment == "left" {
			sb.WriteString("|:")
			sb.WriteString(strings.Repeat("-", c.width-2))
			sb.WriteString(" ")
			continue
		} else if c.alignment == "center" {
			sb.WriteString("|:")
			sb.WriteString(strings.Repeat("-", c.width-2))
			sb.WriteString(":")
			continue
		} else if c.alignment == "right" {
			sb.WriteString("| ")
			sb.WriteString(strings.Repeat("-", c.width-2))
			sb.WriteString(":")
			continue
		}

		sb.WriteString("| ")
		sb.WriteString(strings.Repeat("-", c.width-2))
		sb.WriteString(" ")
	}
	sb.WriteString("|\n")

	// render rows
	for _, row := range m.rows {
		for i, c := range row {
			sb.WriteString("| ")
			sb.WriteString(padOrTruncate(c.value(), m.cols[i].width-1))
		}
		sb.WriteString("|\n")
	}

	t.text = sb.String()
}

func (t *TableWriter) replaceTable(fp *os.File, idx int) []byte {

	tl := tableLocator{}
	tl.findLocations(fp)

	fp.Seek(0, 0)
	b, _ := io.ReadAll(fp)
	b = append(b[:tl.locs[idx].Start-1],
		append([]byte(t.text), b[tl.locs[idx].End:]...)...)
	return b
}

var (
	tableDelimLeft    = regexp.MustCompile(`^\s*\:\-+\s*$`)
	tableDelimRight   = regexp.MustCompile(`^\s*\-+\:\s*$`)
	tableDelimCenter  = regexp.MustCompile(`^\s*\:\-+\:\s*$`)
	tableDelimNone    = regexp.MustCompile(`^\s*\-+\s*$`)
	thematicBreak     = regexp.MustCompile(`^\s{0,3}((-\s*){3,}|(\*\s*){3,}|(_\s*){3,})\s*$`)
	prefixIgnoreSpace = regexp.MustCompile(`^\s{0,3}`)
	fencedCodeBlock   = regexp.MustCompile("^```|~~~.*$")
	codeIndent        = regexp.MustCompile(`^\s{4,}`)
)

func (tl *tableLocator) findLocations(fp io.Reader) {
	scanner := bufio.NewScanner(fp)

	var (
		prevlen  int
		prevline string
		pos      int
		start    int
	)

	for scanner.Scan() {
		l := scanner.Text()

		if tl.isCodeFence(l) {
			tl.inCodeBlock = !tl.inCodeBlock
			tl.codeFence = trimSpace(l)
		}

		if tl.inTable {
			if isBlankLine(l) || isThematicBreak(l) {
				tl.inTable = false
				tl.locs = append(tl.locs, TableLocation{Start: start, End: pos})
			}
		}

		pos += len(l) + 1

		if tl.inCodeBlock {
			prevline = l
			prevlen = len(l) + 1
			continue
		}

		if isTableDelimiter(l) && isTableHeader(prevline, l) {
			tl.inTable = true
			start = pos - len(l) - prevlen
		}

		prevline = l
		prevlen = len(l) + 1
	}
	if tl.inTable {
		tl.locs = append(tl.locs, TableLocation{Start: start, End: pos})
	}
}

func isTableHeader(header string, delim string) bool {
	if isBlankLine(header) {
		return false
	}
	if isCodeIndent(header) {
		return false
	}
	return len(strings.Split(
		trimPipe(trimLeadingSpace(header)), "|")) <=
		len(strings.Split(
			trimPipe(trimLeadingSpace(delim)), "|"))
}

func isTableDelimiter(s string) bool {

	if isCodeIndent(s) {
		return false
	}

	delim, _, _ := strings.Cut(
		trimPipe(trimLeadingSpace(s)), "|")

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

func (tl *tableLocator) isCodeFence(s string) bool {
	if !tl.inCodeBlock {
		return fencedCodeBlock.MatchString(s)
	} else {
		return strings.HasPrefix(tl.codeFence, trimSpace(s))
	}
}

func isCodeIndent(s string) bool {
	return codeIndent.MatchString(s)
}

func isThematicBreak(s string) bool {
	return thematicBreak.MatchString(s)
}

func isBlankLine(s string) bool {
	return s == ""
}

func trimLeadingSpace(s string) string {
	return prefixIgnoreSpace.ReplaceAllString(s, "")
}

func trimPipe(s string) string {
	return strings.Trim(s, "|")
}

func trimSpace(s string) string {
	return strings.TrimSpace(s)
}
