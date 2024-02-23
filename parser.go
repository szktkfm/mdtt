package mdtt

import (
	"bytes"
	"io"
	"os"

	"github.com/charmbracelet/log"
	"github.com/muesli/termenv"
	"github.com/yuin/goldmark"
	east "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	astext "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

var highPriority = 100

func parse(file string) Model {
	f, _ := os.Open(file)
	s, _ := io.ReadAll(f)

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithExtensions(extension.Table),
	)

	tr := NewRenderer(Options{})

	md.SetRenderer(
		renderer.NewRenderer(
			renderer.WithNodeRenderers(
				util.Prioritized(tr, highPriority),
			),
		),
	)

	// Convert markdown to HTML
	var buf bytes.Buffer
	md.Convert(s, &buf)

	// fmt.Println(string(s))
	// log.Debug(buf.String())
	log.Debug("rows", "rows", tr.rows)
	log.Debug("cols", "cols", tr.cols)

	var rows []NaiveRow
	for i := range len(tr.rows) / len(tr.cols) {
		rows = append(rows, tr.rows[i*len(tr.cols):(i+1)*len(tr.cols)])
	}

	var cols []Column
	for i := range len(tr.cols) {
		cols = append(cols, Column{Title: NewCell(tr.cols[i]), Width: 20})
	}

	t := New(
		WithColumns(cols),
		WithNaiveRows(rows),
		WithFocused(true),
		// table.WithHeight(7),
	)

	style := DefaultStyles()

	t.SetStyles(style)

	return Model{t}
}

// Options is used to configure an ANSIRenderer.
type Options struct {
	BaseURL          string
	WordWrap         int
	PreserveNewLines bool
	ColorProfile     termenv.Profile
}

// ModelBuilder build tea.Model from  markdown
type ModelBuilder struct {
	inTable bool
	rows    []string
	cols    []string
}

// NewRenderer returns a new ANSIRenderer with style and options set.
func NewRenderer(options Options) *ModelBuilder {
	return &ModelBuilder{}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *ModelBuilder) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks
	reg.Register(ast.KindDocument, r.renderNode)
	reg.Register(ast.KindHeading, r.renderNode)
	reg.Register(ast.KindBlockquote, r.renderNode)
	reg.Register(ast.KindCodeBlock, r.renderNode)
	reg.Register(ast.KindFencedCodeBlock, r.renderNode)
	reg.Register(ast.KindHTMLBlock, r.renderNode)
	reg.Register(ast.KindList, r.renderNode)
	reg.Register(ast.KindListItem, r.renderNode)
	reg.Register(ast.KindParagraph, r.renderNode)
	reg.Register(ast.KindTextBlock, r.renderNode)
	reg.Register(ast.KindThematicBreak, r.renderNode)

	// inlines
	reg.Register(ast.KindAutoLink, r.renderNode)
	reg.Register(ast.KindCodeSpan, r.renderNode)
	reg.Register(ast.KindEmphasis, r.renderNode)
	reg.Register(ast.KindImage, r.renderNode)
	reg.Register(ast.KindLink, r.renderNode)
	reg.Register(ast.KindRawHTML, r.renderNode)
	reg.Register(ast.KindText, r.renderNode)
	reg.Register(ast.KindString, r.renderNode)

	// tables
	reg.Register(astext.KindTable, r.renderNode)
	reg.Register(astext.KindTableHeader, r.renderNode)
	reg.Register(astext.KindTableRow, r.renderNode)
	reg.Register(astext.KindTableCell, r.renderNode)

	// definitions
	reg.Register(astext.KindDefinitionList, r.renderNode)
	reg.Register(astext.KindDefinitionTerm, r.renderNode)
	reg.Register(astext.KindDefinitionDescription, r.renderNode)

	// footnotes
	reg.Register(astext.KindFootnote, r.renderNode)
	reg.Register(astext.KindFootnoteList, r.renderNode)
	reg.Register(astext.KindFootnoteLink, r.renderNode)
	reg.Register(astext.KindFootnoteBacklink, r.renderNode)

	// checkboxes
	reg.Register(astext.KindTaskCheckBox, r.renderNode)

	// strikethrough
	reg.Register(astext.KindStrikethrough, r.renderNode)

	// emoji
	reg.Register(east.KindEmoji, r.renderNode)
}

func (r *ModelBuilder) renderNode(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	// children get rendered by their parent
	// if isChild(node) {
	// 	return ast.WalkContinue, nil
	// }

	if entering {
		log.Debugf(">Start %v = %v", node.Kind().String(), string(node.Text(source)))
		log.Debugf("\n")
		log.Debugf("\n")

		switch node.Kind() {

		case ast.KindDocument:
		}

		if node.Kind() == astext.KindTable {
			r.inTable = true
		}

		if r.inTable {
			if node.Kind() == astext.KindTableCell {
				switch node.Parent().Kind() {
				case astext.KindTableHeader:
					r.cols = append(r.cols, string(node.Text(source)))
				case astext.KindTableRow:
					r.rows = append(r.rows, string(node.Text(source)))
				}
			}
		}

	} else {
		log.Debugf("<End %v", node.Kind().String())
		log.Debugf("\n")

		if node.Kind() == astext.KindTable {
			r.inTable = false
		}

		switch node.Kind() {

		case ast.KindDocument:
		}

	}

	return ast.WalkContinue, nil
}

func isChild(node ast.Node) bool {
	if node.Parent() != nil && node.Parent().Kind() == ast.KindBlockquote {
		// skip paragraph within blockquote to avoid reflowing text
		return true
	}
	for n := node.Parent(); n != nil; n = n.Parent() {
		// These types are already rendered by their parent
		switch n.Kind() {
		case ast.KindLink, ast.KindImage, ast.KindEmphasis, astext.KindStrikethrough, astext.KindTableCell:
			return true
		}
	}

	return false
}
