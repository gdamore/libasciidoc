package sgml

import (
	"bytes"
	"strconv"
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (sr *sgmlRenderer) renderPreamble(ctx *Context, p types.Preamble) ([]byte, error) {
	log.Debugf("rendering preamble...")
	result := &bytes.Buffer{}
	// the <div id="preamble"> wrapper is only necessary
	// if the document has a section 0
	err := sr.preamble.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			Wrapper  bool
			Elements []interface{}
		}{
			Wrapper:  ctx.HasHeader,
			Elements: p.Elements,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error while rendering preamble")
	}
	// log.Debugf("rendered preamble: %s", result.Bytes())
	return result.Bytes(), nil
}

func (sr *sgmlRenderer) renderSection(ctx *Context, s types.Section) ([]byte, error) {
	log.Debugf("rendering section level %d", s.Level)
	renderedSectionTitle, err := sr.renderSectionTitle(ctx, s)
	if err != nil {
		return nil, errors.Wrapf(err, "error while rendering section")
	}
	result := &bytes.Buffer{}
	// select the appropriate template for the section
	var tmpl *textTemplate
	if s.Level == 1 {
		tmpl = sr.sectionOne
	} else {
		tmpl = sr.sectionContent
	}
	err = tmpl.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			Class        string
			SectionTitle string
			Elements     []interface{}
		}{
			Class:        "sect" + strconv.Itoa(s.Level),
			SectionTitle: renderedSectionTitle,
			Elements:     s.Elements,
		}})
	if err != nil {
		return nil, errors.Wrapf(err, "error while rendering section")
	}
	// log.Debugf("rendered section: %s", result.Bytes())
	return result.Bytes(), nil
}

func (sr *sgmlRenderer) renderSectionTitle(ctx *Context, s types.Section) (string, error) {
	result := &bytes.Buffer{}
	renderedContent, err := sr.renderInlineElements(ctx, s.Title)
	if err != nil {
		return "", errors.Wrapf(err, "error while rendering sectionTitle content")
	}
	renderedContentStr := strings.TrimSpace(string(renderedContent))
	id := sr.renderElementID(s.Attributes)
	err = sr.sectionHeader.Execute(result, struct {
		Level   int
		ID      string
		Content string
	}{
		Level:   s.Level + 1,
		ID:      id,
		Content: renderedContentStr,
	})
	if err != nil {
		return "", errors.Wrapf(err, "error while rendering sectionTitle")
	}
	// log.Debugf("rendered sectionTitle: %s", result.Bytes())
	return result.String(), nil
}
