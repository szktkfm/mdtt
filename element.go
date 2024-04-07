package mdtt

import (
	"bytes"
	"strings"

	"github.com/charmbracelet/log"
	east "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	astext "github.com/yuin/goldmark/extension/ast"
)

type Element struct {
	Entering string
	Renderer func(b *bytes.Buffer)
	Finisher func(b *bytes.Buffer)
}

// NewElement returns the appropriate render Element for a given node.
func (tr *ModelBuilder) NewElement(node ast.Node, source []byte) Element {

	switch node.Kind() {

	case ast.KindLink:
		n := node.(*ast.Link)
		log.Debug(n)
		return Element{
			Renderer: func(b *bytes.Buffer) {
				b.WriteString("[")
			},
			Finisher: func(b *bytes.Buffer) {
				b.WriteString("](")
				b.WriteString(string(n.Destination))
				b.WriteString(")")
			},
		}
	case ast.KindAutoLink:
		n := node.(*ast.AutoLink)
		u := string(n.URL(source))
		if n.AutoLinkType == ast.AutoLinkEmail && !strings.HasPrefix(strings.ToLower(u), "mailto:") {
			u = "mailto:" + u
		}

		return Element{
			Renderer: func(b *bytes.Buffer) {
				b.WriteString(u)
			},
			Finisher: func(b *bytes.Buffer) {
				b.WriteString("")
			},
		}

	case ast.KindCodeSpan:
		return Element{
			Renderer: func(b *bytes.Buffer) {
				b.WriteString("`")
			},
			Finisher: func(b *bytes.Buffer) {
				b.WriteString("`")
			},
		}

	case ast.KindEmphasis:
		n := node.(*ast.Emphasis)
		return Element{
			Renderer: func(b *bytes.Buffer) {
				if n.Level == 2 {
					b.WriteString("**")
					return
				}
				b.WriteString("_")
			},
			Finisher: func(b *bytes.Buffer) {
				if n.Level == 2 {
					b.WriteString("**")
					return
				}
				b.WriteString("_")
			},
		}

	case astext.KindTableCell:
		return Element{
			Renderer: func(b *bytes.Buffer) {
			},
			Finisher: func(b *bytes.Buffer) {
			},
		}

	case astext.KindTableRow:
		return Element{
			Renderer: func(b *bytes.Buffer) {
			},
			Finisher: func(b *bytes.Buffer) {
			},
		}

	case ast.KindText:
		return Element{
			Renderer: func(b *bytes.Buffer) {
				b.WriteString(string(node.Text(source)))
			},
			Finisher: func(b *bytes.Buffer) {
				b.WriteString("")
			},
		}

	case east.KindEmoji:
		n := node.(*east.Emoji)
		return Element{
			Renderer: func(b *bytes.Buffer) {
				b.WriteString(string(n.Value.Unicode))
			},
			Finisher: func(b *bytes.Buffer) {
				b.WriteString("")
			},
		}

	default:
		return Element{
			Renderer: func(b *bytes.Buffer) {
			},
			Finisher: func(b *bytes.Buffer) {
			},
		}
	}
}
