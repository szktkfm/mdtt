package mdtt

import (
	"bytes"

	"github.com/yuin/goldmark"
	east "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	astext "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

var highPriority = 100

// tableModelBuilder build TableModel from markdown
type tableModelBuilder struct {
	inTable bool
	buf     *bytes.Buffer
	// temporary storage of table rows
	rows []string
	// temporary storage of table columns
	cols []string
	// temporary storage of table alignments
	alignment []string
	tables    []table
}

type table struct {
	rows      []string
	cols      []string
	alignment []string
}

func parse(s []byte) []TableModel {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithExtensions(extension.Table),
	)

	builder := NewTableModelBuilder()

	md.SetRenderer(
		renderer.NewRenderer(
			renderer.WithNodeRenderers(
				util.Prioritized(builder, highPriority),
			),
		),
	)

	var _buf bytes.Buffer
	md.Convert(s, &_buf)

	return builder.build()

}

func (b *tableModelBuilder) build() []TableModel {

	var models []TableModel
	for _, t := range b.tables {
		var rows []naiveRow
		for i := range len(t.rows) / len(t.cols) {
			rows = append(rows, t.rows[i*len(t.cols):(i+1)*len(t.cols)])
		}

		var cols []column
		for i := range len(t.cols) {
			var maxWidth int
			for _, r := range rows {
				maxWidth = max(len(r[i]), maxWidth)
			}
			maxWidth = max(maxWidth, len(t.cols[i]))
			cols = append(cols, column{
				title:     NewCell(t.cols[i]),
				width:     maxWidth + 2,
				alignment: t.alignment[i],
			})
		}

		t := NewTableModel(
			WithColumns(cols),
			WithNaiveRows(rows),
			WithFocused(true),
			WithHeight(len(rows)+1),
		)

		style := defaultStyles()

		t.SetStyles(style)
		models = append(models, t)
	}
	return models
}

func NewTableModelBuilder() *tableModelBuilder {
	return &tableModelBuilder{buf: bytes.NewBuffer(nil)}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *tableModelBuilder) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// blocks
	reg.Register(ast.KindDocument, r.consumeNode)
	reg.Register(ast.KindHeading, r.consumeNode)
	reg.Register(ast.KindBlockquote, r.consumeNode)
	reg.Register(ast.KindCodeBlock, r.consumeNode)
	reg.Register(ast.KindFencedCodeBlock, r.consumeNode)
	reg.Register(ast.KindHTMLBlock, r.consumeNode)
	reg.Register(ast.KindList, r.consumeNode)
	reg.Register(ast.KindListItem, r.consumeNode)
	reg.Register(ast.KindParagraph, r.consumeNode)
	reg.Register(ast.KindTextBlock, r.consumeNode)
	reg.Register(ast.KindThematicBreak, r.consumeNode)

	// inlines
	reg.Register(ast.KindAutoLink, r.consumeNode)
	reg.Register(ast.KindCodeSpan, r.consumeNode)
	reg.Register(ast.KindEmphasis, r.consumeNode)
	reg.Register(ast.KindImage, r.consumeNode)
	reg.Register(ast.KindLink, r.consumeNode)
	reg.Register(ast.KindRawHTML, r.consumeNode)
	reg.Register(ast.KindText, r.consumeNode)
	reg.Register(ast.KindString, r.consumeNode)

	// tables
	reg.Register(astext.KindTable, r.consumeNode)
	reg.Register(astext.KindTableHeader, r.consumeNode)
	reg.Register(astext.KindTableRow, r.consumeNode)
	reg.Register(astext.KindTableCell, r.consumeNode)

	// definitions
	reg.Register(astext.KindDefinitionList, r.consumeNode)
	reg.Register(astext.KindDefinitionTerm, r.consumeNode)
	reg.Register(astext.KindDefinitionDescription, r.consumeNode)

	// footnotes
	reg.Register(astext.KindFootnote, r.consumeNode)
	reg.Register(astext.KindFootnoteList, r.consumeNode)
	reg.Register(astext.KindFootnoteLink, r.consumeNode)
	reg.Register(astext.KindFootnoteBacklink, r.consumeNode)

	// checkboxes
	reg.Register(astext.KindTaskCheckBox, r.consumeNode)

	// strikethrough
	reg.Register(astext.KindStrikethrough, r.consumeNode)

	// emoji
	reg.Register(east.KindEmoji, r.consumeNode)
}

func (tb *tableModelBuilder) consumeNode(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if node.Kind() == astext.KindTable {
			tb.inTable = true
		}

		if tb.inTable {
			e := tb.NewElement(node, source)
			e.open(tb.buf)
		}

	} else {
		if node.Kind() == astext.KindTable {
			tb.tables = append(tb.tables, table{tb.rows, tb.cols, tb.alignment})
			tb.rows = nil
			tb.cols = nil
			tb.alignment = nil
			tb.inTable = false
		}

		if tb.inTable {
			switch node.Kind() {
			case astext.KindTableCell:
				switch node.Parent().Kind() {
				case astext.KindTableHeader:
					n := node.(*astext.TableCell)
					tb.cols = append(tb.cols, tb.buf.String())
					tb.alignment = append(tb.alignment, n.Alignment.String())
					tb.buf.Reset()
				case astext.KindTableRow:
					tb.rows = append(tb.rows, tb.buf.String())
					tb.buf.Reset()
				}
			default:
				e := tb.NewElement(node, source)
				e.close(tb.buf)
			}
		}
	}
	return ast.WalkContinue, nil
}
