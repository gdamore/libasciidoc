package sgml

import (
	"bytes"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/pkg/errors"
)

func (sr *sgmlRenderer) renderCalloutList(ctx *Context, l types.CalloutList) ([]byte, error) {
	result := &bytes.Buffer{}
	err := sr.calloutList.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			ID    string
			Title string
			Role  string
			Items []types.CalloutListItem
		}{
			ID:    sr.renderElementID(l.Attributes),
			Title: l.Attributes.GetAsStringWithDefault(types.AttrTitle, ""),
			Role:  l.Attributes.GetAsStringWithDefault(types.AttrRole, ""),
			Items: l.Items,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to render callout list")
	}
	return result.Bytes(), nil
}
