package sgml

import (
	"github.com/bytesparadise/libasciidoc/pkg/types"
)

func (sr *sgmlRenderer) renderIndexTerm(ctx *Context, t types.IndexTerm) ([]byte, error) {
	return sr.renderInlineElements(ctx, t.Term)
}

func (sr *sgmlRenderer) renderConcealedIndexTerm(_ types.ConcealedIndexTerm) ([]byte, error) {
	return []byte{}, nil // do not render
}
