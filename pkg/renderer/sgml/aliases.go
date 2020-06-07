package sgml

import (
	html "html/template"
	text "text/template"

	"github.com/bytesparadise/libasciidoc/pkg/renderer"
	"github.com/bytesparadise/libasciidoc/pkg/types"
)

// These type aliases provide local names for names in other packages,
// thereby helping minimize collisions based on conflicting package
// names.  It also reduces the imports we have to use everywhere else.

type Context = renderer.Context
type Document = types.Document
type Metadata = types.Metadata
type textTemplate = text.Template
type funcMap = text.FuncMap

// sanitized is for post-render output, which is already sanitized
// and should be considered safe.
type sanitized = html.HTML
