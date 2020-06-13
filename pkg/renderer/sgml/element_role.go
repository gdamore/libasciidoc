package sgml

import (
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
)

func (r *sgmlRenderer) renderElementRole(attrs types.Attributes) string {
	a := attrs[types.AttrRole]
	roles := a.([]string)
	return strings.Join(roles, " ")
}
