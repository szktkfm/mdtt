package mdtt

import (
	"fmt"

	east "github.com/yuin/goldmark-emoji/ast"
	"github.com/yuin/goldmark/ast"
	astext "github.com/yuin/goldmark/extension/ast"
)

type Element struct {
	Entering string
	Exiting  string
	Renderer func()
	Finisher func()
}

// NewElement returns the appropriate render Element for a given node.
func (tr *ModelBuilder) NewElement(node ast.Node, source []byte) Element {
	// ctx := tr.context
	// fmt.Print(strings.Repeat("  ", ctx.blockStack.Len()), node.Type(), node.Kind())
	// defer fmt.Println()

	switch node.Kind() {
	// // Document
	// case ast.KindDocument:
	// // Heading
	// case ast.KindHeading:
	// // Paragraph
	// case ast.KindParagraph:
	// // Blockquote
	// case ast.KindBlockquote:
	// // Lists
	// case ast.KindList:

	// case ast.KindListItem:
	// // Text Elements
	// case ast.KindText:

	// case ast.KindEmphasis:

	// case astext.KindStrikethrough:
	case ast.KindThematicBreak:
		// return Element{
		// 	Entering: "",
		// 	Exiting:  "",
		// 	Renderer: &BaseElement{
		// 		Style: ctx.options.Styles.HorizontalRule,
		// 	},
		// }

	// Links
	case ast.KindLink:
		//TODO
		// n := node.(*ast.Link)
		// return Element{
		// 	Renderer: &LinkElement{
		// 		Text:    textFromChildren(node, source),
		// 		BaseURL: ctx.options.BaseURL,
		// 		URL:     string(n.Destination),
		// 	},
		// }

	case ast.KindAutoLink:
		// n := node.(*ast.AutoLink)
		// u := string(n.URL(source))
		// label := string(n.Label(source))
		// if n.AutoLinkType == ast.AutoLinkEmail && !strings.HasPrefix(strings.ToLower(u), "mailto:") {
		// 	u = "mailto:" + u
		// }

		// return Element{
		// 	Renderer: &LinkElement{
		// 		Text:    label,
		// 		BaseURL: ctx.options.BaseURL,
		// 		URL:     u,
		// 	},
		// }

	// Images
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

	case ast.KindCodeSpan:
		// // n := node.(*ast.CodeSpan)
		// e := &BlockElement{
		// 	Block: &bytes.Buffer{},
		// 	Style: cascadeStyle(ctx.blockStack.Current().Style, ctx.options.Styles.Code, false),
		// }
		// return Element{
		// 	Renderer: e,
		// 	Finisher: e,
		// }

	// Tables
	case astext.KindTable:
		// te := &TableElement{}
		// return Element{
		// 	Entering: "\n",
		// 	Renderer: te,
		// 	Finisher: te,
		// }

	case astext.KindTableCell:
		// s := ""
		// n := node.FirstChild()
		// for n != nil {
		// 	switch t := n.(type) {
		// 	case *ast.AutoLink:
		// 		s += string(t.Label(source))
		// 	default:
		// 		s += string(n.Text(source))
		// 	}

		// 	n = n.NextSibling()
		// }

		// return Element{
		// 	Renderer: &TableCellElement{
		// 		Text: s,
		// 		Head: node.Parent().Kind() == astext.KindTableHeader,
		// 	},
		// }

	case astext.KindTableHeader:
		// return Element{
		// 	Finisher: &TableHeadElement{},
		// }
	case astext.KindTableRow:
		// return Element{
		// 	Finisher: &TableRowElement{},
		// }

	// HTML Elements
	case ast.KindHTMLBlock:
	// 	n := node.(*ast.HTMLBlock)
	// 	return Element{
	// 		Renderer: &BaseElement{
	// 			Token: ctx.SanitizeHTML(string(n.Text(source)), true),
	// 			Style: ctx.options.Styles.HTMLBlock.StylePrimitive,
	// 		},
	// 	}
	// case ast.KindRawHTML:
	// 	n := node.(*ast.RawHTML)
	// 	return Element{
	// 		Renderer: &BaseElement{
	// 			Token: ctx.SanitizeHTML(string(n.Text(source)), true),
	// 			Style: ctx.options.Styles.HTMLSpan.StylePrimitive,
	// 		},
	// 	}

	// Definition Lists

	// Handled by parents
	case astext.KindTaskCheckBox:
		// // handled by KindListItem
		// return Element{}
	case ast.KindTextBlock:
		// return Element{}

	case east.KindEmoji:
		// n := node.(*east.Emoji)
		// return Element{
		// 	Renderer: &BaseElement{
		// 		Token: string(n.Value.Unicode),
		// 	},
		// }

	// Unknown case
	default:
		fmt.Println("Warning: unhandled element", node.Kind().String())
		return Element{}
	}
	return Element{}
}
