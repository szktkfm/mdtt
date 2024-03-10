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
	defer f.Close()
	s, _ := io.ReadAll(f)

	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithExtensions(extension.Table),
	)

	builder := NewModelBuilder(Options{})

	md.SetRenderer(
		renderer.NewRenderer(
			renderer.WithNodeRenderers(
				util.Prioritized(builder, highPriority),
			),
		),
	)

	// Convert markdown to HTML
	var buf bytes.Buffer
	md.Convert(s, &buf)

	// fmt.Println(string(s))
	// log.Debug(buf.String())
	log.Debug("rows", "rows", builder.rows)
	log.Debug("cols", "cols", builder.cols)

	// table.WithHeight(7),
	m := builder.build()

	return m
}

func (b *ModelBuilder) build() Model {
	var rows []NaiveRow
	for i := range len(b.rows) / len(b.cols) {
		rows = append(rows, b.rows[i*len(b.cols):(i+1)*len(b.cols)])
	}

	var cols []Column
	for i := range len(b.cols) {
		cols = append(cols, Column{Title: NewCell(b.cols[i]), Width: 20})
	}

	t := New(
		WithColumns(cols),
		WithNaiveRows(rows),
		WithFocused(true),
	)

	style := DefaultStyles()

	t.SetStyles(style)
	m := Model{t}
	return m
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
	buf     *bytes.Buffer
	rows    []string
	cols    []string
}

// NewModelBuilder returns a new ANSIRenderer with style and options set.
func NewModelBuilder(options Options) *ModelBuilder {
	var buf []byte
	return &ModelBuilder{
		buf: bytes.NewBuffer(buf),
	}
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
	if entering {
		log.Debugf(">Start %v = %v", node.Kind().String(), string(node.Text(source)))

		if node.Kind() == astext.KindTable {
			r.inTable = true
		}

		if r.inTable {
			e := r.NewElement(node, source)
			e.Renderer(r.buf)
		}

	} else {
		log.Debugf("<End %v", node.Kind().String())

		if node.Kind() == astext.KindTable {
			r.inTable = false
		}

		if r.inTable {
			switch node.Kind() {
			case astext.KindTableCell:
				switch node.Parent().Kind() {
				case astext.KindTableHeader:
					r.cols = append(r.cols, r.buf.String())
					r.buf.Reset()
				case astext.KindTableRow:
					r.rows = append(r.rows, r.buf.String())
					r.buf.Reset()
				}
			default:
				e := r.NewElement(node, source)
				e.Finisher(r.buf)
			}
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
