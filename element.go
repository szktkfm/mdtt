package mdtt

import (
	"bytes"

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

	// case ast.KindLink:
	// 	return Element{}

	// case ast.KindAutoLink:
	// 	return Element{}

	case ast.KindImage:
		// n := node.(*ast.Image)
		// text := string(n.Text(source))
		// return Element{
		// 	Renderer: &ImageElement{
		// 		Text:    text,
		// 		BaseURL: ctx.options.BaseURL,
		// 		URL:     string(n.Destination),
		// 	},
		// }
		return Element{}

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
		// n := node.(*ast.CodeSpan)
		return Element{
			Renderer: func(b *bytes.Buffer) {
				b.WriteString(string(node.Text(source)))
			},
			Finisher: func(b *bytes.Buffer) {
				b.WriteString("")
			},
		}

		//
	// case astext.KindTaskCheckBox:
	// case ast.KindTextBlock:

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
