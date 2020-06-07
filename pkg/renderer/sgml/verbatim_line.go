package sgml

import (
	"bytes"
	"github.com/bytesparadise/libasciidoc/pkg/types"
)

func (sr *sgmlRenderer) renderVerbatimLine(l types.VerbatimLine) ([]byte, error) {
	result := &bytes.Buffer{}
	if err := sr.verbatimLine.Execute(result, l); err != nil {
		return nil, err
	}
	return result.Bytes(), nil
}
