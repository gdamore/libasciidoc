package sgml

import (
	"bytes"
	"html/template"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/pkg/errors"
)

func (sr *sgmlRenderer) renderInlinePassthrough(ctx *Context, p types.InlinePassthrough) ([]byte, error) {
	renderedContent, err := sr.renderPassthroughContent(ctx, p)
	if err != nil {
		return nil, errors.Wrap(err, "unable to render passthrough")
	}
	switch p.Kind {
	case types.SinglePlusPassthrough:
		// rendered passthrough content is in an HTML-escaped form
		buf := &bytes.Buffer{}
		template.HTMLEscape(buf, renderedContent)
		return buf.Bytes(), nil
	default:
		return renderedContent, nil
	}
}

// renderPassthroughMacro renders the passthrough content in its raw from
func (sr *sgmlRenderer) renderPassthroughContent(ctx *Context, p types.InlinePassthrough) ([]byte, error) {
	buf := &bytes.Buffer{}
	for _, element := range p.Elements {
		switch element := element.(type) {
		case types.StringElement:
			// "string" elements must be rendered as-is, ie, without any HTML escaping.
			_, err := buf.WriteString(element.Content)
			if err != nil {
				return nil, err
			}
		default:
			renderedElement, err := sr.renderElement(ctx, element)
			if err != nil {
				return nil, err
			}
			_, err = buf.Write(renderedElement)
			if err != nil {
				return nil, err
			}

		}
	}
	return buf.Bytes(), nil
}
