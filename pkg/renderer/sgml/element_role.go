package sgml

import (
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
)

func (r *sgmlRenderer) renderElementRole(attrs types.Attributes) string {
	a, _ := attrs[types.AttrRole]
	roles, _ := a.([]string)
	return strings.Join(roles, " ")
}
