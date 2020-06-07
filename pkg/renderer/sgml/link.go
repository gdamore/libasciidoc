package sgml

import (
	"bytes"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (sr *sgmlRenderer) renderLink(ctx *Context, l types.InlineLink) ([]byte, error) { //nolint: unparam
	result := &bytes.Buffer{}
	location := l.Location.String()
	var text []byte
	class := ""
	var err error
	// TODO; support `mailto:` positional attributes
	positionals := l.Attributes.Positionals()
	if len(positionals) > 0 {
		buf := &bytes.Buffer{}
		for i, arg := range positionals {
			t, err := sr.renderInlineElements(ctx, arg)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to render external link")
			}
			buf.Write(t)
			if i < len(positionals)-1 {
				buf.WriteString(",")
			}
		}
		text = buf.Bytes()
	} else {
		class = "bare"
		text = []byte(location)
	}
	err = sr.link.Execute(result, struct {
		URL   string
		Text  string
		Class string
	}{
		URL:   location,
		Text:  string(text),
		Class: class,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to render external link")
	}
	log.Debugf("rendered external link: %s", result.Bytes())
	return result.Bytes(), nil
}
