package mdtt

import (
	"bytes"
	"strings"

	east "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	astext "github.com/yuin/goldmark/extension/ast"
)

type element struct {
	open  func(b *bytes.Buffer)
	close func(b *bytes.Buffer)
}

// NewElement returns the appropriate render Element for a given node.
func (tr *tableModelBuilder) NewElement(node ast.Node, source []byte) element {

	switch node.Kind() {

	case ast.KindLink:
		n := node.(*ast.Link)
		return element{
			open: func(b *bytes.Buffer) {
				b.WriteString("[")
			},
			close: func(b *bytes.Buffer) {
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

		return element{
			open: func(b *bytes.Buffer) {
				b.WriteString(u)
			},
			close: func(b *bytes.Buffer) {
				b.WriteString("")
			},
		}

	case ast.KindCodeSpan:
		return element{
			open: func(b *bytes.Buffer) {
				b.WriteString("`")
			},
			close: func(b *bytes.Buffer) {
				b.WriteString("`")
			},
		}

	case ast.KindEmphasis:
		n := node.(*ast.Emphasis)
		return element{
			open: func(b *bytes.Buffer) {
				if n.Level == 2 {
					b.WriteString("**")
					return
				}
				b.WriteString("_")
			},
			close: func(b *bytes.Buffer) {
				if n.Level == 2 {
					b.WriteString("**")
					return
				}
				b.WriteString("_")
			},
		}

	case astext.KindTableCell:
		return element{
			open: func(b *bytes.Buffer) {
			},
			close: func(b *bytes.Buffer) {
			},
		}

	case astext.KindTableRow:
		return element{
			open: func(b *bytes.Buffer) {
			},
			close: func(b *bytes.Buffer) {
			},
		}

	case ast.KindText:
		return element{
			open: func(b *bytes.Buffer) {
				b.WriteString(string(node.Text(source)))
			},
			close: func(b *bytes.Buffer) {
				b.WriteString("")
			},
		}

	case east.KindEmoji:
		n := node.(*east.Emoji)
		return element{
			open: func(b *bytes.Buffer) {
				b.WriteString(string(n.Value.Unicode))
			},
			close: func(b *bytes.Buffer) {
				b.WriteString("")
			},
		}

	default:
		return element{
			open: func(b *bytes.Buffer) {
			},
			close: func(b *bytes.Buffer) {
			},
		}
	}
}
