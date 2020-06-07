package sgml

import "github.com/bytesparadise/libasciidoc/pkg/types"

func (sr *sgmlRenderer) renderElementID(attrs types.Attributes) string {
	if id, ok := attrs[types.AttrID].(string); ok {
		return id
	}
	return ""
}
