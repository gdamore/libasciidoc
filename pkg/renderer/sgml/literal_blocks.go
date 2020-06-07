package sgml

import (
	"bytes"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"math"
	"strings"
)

func (sr *sgmlRenderer) renderLiteralBlock(ctx *Context, b types.LiteralBlock) ([]byte, error) {
	log.Debugf("rendering delimited block with content: %s", b.Lines)
	var lines []string
	if t, found := b.Attributes.GetAsString(types.AttrLiteralBlockType); found && t == types.LiteralBlockWithSpacesOnFirstLine {
		if len(b.Lines) == 1 {
			lines = []string{strings.TrimLeft(b.Lines[0], " ")}
		} else {
			lines = make([]string, len(b.Lines))
			// remove as many spaces as needed on each line
			spaceCount := float64(0)
			// first pass to determine the minimum number of spaces to remove
			for i, line := range b.Lines {
				l := strings.TrimLeft(line, " ")
				if i == 0 {
					spaceCount = float64(len(line) - len(l))
				} else {
					spaceCount = math.Min(spaceCount, float64(len(line)-len(l)))
				}
			}
			log.Debugf("trimming %d space(s) on each line", int(spaceCount))
			// then remove the same number of spaces on each line
			spaces := strings.Repeat(" ", int(spaceCount))
			for i, line := range b.Lines {
				lines[i] = strings.TrimPrefix(line, spaces)
			}
		}
	} else {
		lines = b.Lines
	}
	result := &bytes.Buffer{}
	err := sr.literalBlock.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			ID    string
			Title string
			Lines []string
		}{
			ID:    sr.renderElementID(b.Attributes),
			Title: sr.renderElementTitle(b.Attributes),
			Lines: lines,
		}})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to render delimited block")
	}
	return result.Bytes(), nil
}
