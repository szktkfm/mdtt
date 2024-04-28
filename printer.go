package mdtt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"

	"github.com/charmbracelet/log"
)

type tableWriter struct {
	// rendered table text
	text string
}

type tableLocator struct {
	ranges      []tableRange
	inTable     bool
	codeFence   string
	inCodeBlock bool
}

type tableRange struct {
	Start int
	End   int
}

func Write(m Model) {
	tw := tableWriter{}
	tw.render(m.table)
	if m.inplace {
		if err := tw.writeFile(m); err != nil {
			log.Error("Error writing file: ", err)
			return
		}
	} else {
		fmt.Print(tw.text)
	}
}
func (tw *tableWriter) writeFile(m Model) error {
	fp, err := os.Open(m.fpath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer fp.Close()
	b := tw.replaceTable(fp, m.choose)

	if err := os.WriteFile(m.fpath, b, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

func (t *tableWriter) render(m TableModel) {
	var sb strings.Builder
	var width int

	newline := "\n"
	if runtime.GOOS == "windows" {
		newline = "\r\n"
	}

	// render header
	for _, c := range m.cols {
		sb.WriteString("| ")
		sb.WriteString(padOrTruncate(c.title.value(), max(c.width-1, 2)))
		width += c.width
	}
	sb.WriteString("|" + newline)

	// render delimiter
	for _, c := range m.cols {

		if c.alignment == "left" {
			sb.WriteString("|:")
			sb.WriteString(strings.Repeat("-", max(c.width-2, 1)))
			sb.WriteString(" ")
			continue
		} else if c.alignment == "center" {
			sb.WriteString("|:")
			sb.WriteString(strings.Repeat("-", max(c.width-2, 1)))
			sb.WriteString(":")
			continue
		} else if c.alignment == "right" {
			sb.WriteString("| ")
			sb.WriteString(strings.Repeat("-", max(c.width-2, 1)))
			sb.WriteString(":")
			continue
		}

		sb.WriteString("| ")
		sb.WriteString(strings.Repeat("-", max(c.width-2, 1)))
		sb.WriteString(" ")
	}
	sb.WriteString("|" + newline)

	// render rows
	for _, row := range m.rows {
		for i, c := range row {
			sb.WriteString("| ")
			sb.WriteString(padOrTruncate(c.value(), max(m.cols[i].width-1, 2)))
		}
		sb.WriteString("|" + newline)
	}

	t.text = sb.String()
}

func (t *tableWriter) replaceTable(fp *os.File, idx int) []byte {

	tl := tableLocator{}
	tl.findLocations(fp)

	fp.Seek(0, 0)
	b, _ := io.ReadAll(fp)
	b = append(b[:tl.ranges[idx].Start-1],
		append([]byte(t.text), b[min(len(b), tl.ranges[idx].End):]...)...)
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
		prevLine string
		pos      int
		start    int
	)

	for scanner.Scan() {
		line := scanner.Text()

		if tl.isCodeFence(line) {
			tl.inCodeBlock = !tl.inCodeBlock
			tl.codeFence = trimSpace(line)
		}

		if tl.inTable {
			if isBlankLine(line) || isThematicBreak(line) {
				tl.inTable = false
				tl.ranges = append(tl.ranges, tableRange{Start: start, End: pos})
			}
		}

		pos += len(line) + 1

		if tl.inCodeBlock {
			prevLine = line
			prevlen = len(line) + 1
			continue
		}

		if isTableDelimiter(line) && isTableHeader(prevLine, line) {
			tl.inTable = true
			start = pos - len(line) - prevlen
		}

		prevLine = line
		prevlen = len(line) + 1
	}
	if tl.inTable {
		tl.ranges = append(tl.ranges, tableRange{Start: start, End: pos})
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
