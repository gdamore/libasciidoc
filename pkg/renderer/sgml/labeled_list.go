package sgml

import (
	"bytes"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/pkg/errors"
)

func (sr *sgmlRenderer) renderLabeledList(ctx *Context, l types.LabeledList) ([]byte, error) {
	tmpl, err := sr.getLabeledListTmpl(l)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to render labeled list")
	}

	result := &bytes.Buffer{}
	// here we must preserve the HTML tags
	err = tmpl.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			ID    string
			Title string
			Role  string
			Items []types.LabeledListItem
		}{
			ID:    sr.renderElementID(l.Attributes),
			Title: sr.renderElementTitle(l.Attributes),
			Role:  l.Attributes.GetAsStringWithDefault(types.AttrRole, ""),
			Items: l.Items,
		},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to render labeled list")
	}
	// log.Debugf("rendered labeled list: %s", result.Bytes())
	return result.Bytes(), nil
}

func (sr *sgmlRenderer) getLabeledListTmpl(l types.LabeledList) (*textTemplate, error) {
	if layout, ok := l.Attributes["layout"]; ok {
		switch layout {
		case "horizontal":
			return sr.labeledListHorizontal, nil
		default:
			return nil, errors.Errorf("unsupported labeled list layout: %s", layout)
		}
	}
	if l.Attributes.Has(types.AttrQandA) {
		return sr.qAndAList, nil
	}
	return sr.labeledList, nil
}
