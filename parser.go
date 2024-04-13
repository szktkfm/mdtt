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

// TableModelBuilder build TableModel from markdown
type TableModelBuilder struct {
	inTable bool
	buf     *bytes.Buffer
	// temprary storage of table rows
	rows []string
	// temprary storage of table columns
	cols []string
	// temprary storage of table alignments
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

func (b *TableModelBuilder) build() []TableModel {

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

func NewTableModelBuilder() *TableModelBuilder {
	var buf []byte
	return &TableModelBuilder{
		buf: bytes.NewBuffer(buf),
	}
}

// RegisterFuncs implements NodeRenderer.RegisterFuncs.
func (r *TableModelBuilder) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
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

func (tb *TableModelBuilder) renderNode(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		if node.Kind() == astext.KindTable {
			tb.inTable = true
		}

		if tb.inTable {
			e := tb.NewElement(node, source)
			e.Renderer(tb.buf)
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
				e.Finisher(tb.buf)
			}
		}
	}
	return ast.WalkContinue, nil
}
