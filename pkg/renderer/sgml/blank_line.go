package sgml

import (
	"bytes"
	"github.com/bytesparadise/libasciidoc/pkg/types"

	log "github.com/sirupsen/logrus"
)

func (sr *sgmlRenderer) renderBlankLine(ctx *Context, _ types.BlankLine) ([]byte, error) {

	if ctx.IncludeBlankLine {
		buf := &bytes.Buffer{}
		if err := sr.blankLine.Execute(buf, nil); err != nil {
			return nil, err
		}
		log.Debug("rendering blank line")
		return buf.Bytes(), nil
	}
	return []byte{}, nil
}

func (sr *sgmlRenderer) renderLineBreak() ([]byte, error) {
	buf := &bytes.Buffer{}
	if err := sr.lineBreak.Execute(buf, nil); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
