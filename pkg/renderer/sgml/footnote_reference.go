package sgml

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/pkg/errors"
)

func (sr *sgmlRenderer) renderFootnote(ctx *Context, elements []interface{}) (string, error) {
	result, err := sr.renderInlineElements(ctx, elements)
	if err != nil {
		return "", errors.Wrapf(err, "unable to render foot note content")
	}
	return strings.TrimSpace(string(result)), nil
}

func (sr *sgmlRenderer) renderFootnoteReference(note types.FootnoteReference) ([]byte, error) {
	result := &bytes.Buffer{}
	if note.ID != types.InvalidFootnoteReference && !note.Duplicate {
		// valid case for a footnote with content, with our without an explicit reference
		err := sr.footnote.Execute(result, struct {
			ID  int
			Ref string
		}{
			ID:  note.ID,
			Ref: note.Ref,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "unable to render footnote")
		}
	} else if note.Duplicate {
		// valid case for a footnote with content, with our without an explicit reference
		err := sr.footnoteRef.Execute(result, struct {
			ID  int
			Ref string
		}{
			ID:  note.ID,
			Ref: note.Ref,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "unable to render footnote")
		}
	} else {
		// invalid footnote
		err := sr.invalidFootnote.Execute(result, struct {
			Ref string
		}{
			Ref: note.Ref,
		})
		if err != nil {
			return nil, errors.Wrapf(err, "unable to render missing footnote")
		}
	}
	return result.Bytes(), nil
}

func (sr *sgmlRenderer) renderFootnoteReferencePlainText(note types.FootnoteReference) ([]byte, error) {
	result := &bytes.Buffer{}
	if note.ID != types.InvalidFootnoteReference {
		// valid case for a footnote with content, with our without an explicit reference
		err := sr.footnoteRefPlain.Execute(result, struct {
			ID    int
			Class string
		}{
			ID:    note.ID,
			Class: "footnote",
		})
		if err != nil {
			return nil, errors.Wrapf(err, "unable to render footnote")
		}
	} else {
		return nil, fmt.Errorf("unable to render missing footnote")
	}
	return result.Bytes(), nil
}

func (sr *sgmlRenderer) renderFootnotes(ctx *Context, notes []types.Footnote) ([]byte, error) {
	// skip if there's no foot note in the doc
	if len(notes) == 0 {
		return []byte{}, nil
	}
	result := &bytes.Buffer{}
	err := sr.footnotes.Execute(result,
		ContextualPipeline{
			Context: ctx,
			Data: struct {
				Footnotes []types.Footnote
			}{
				Footnotes: notes,
			},
		})
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to render footnotes")
	}
	return result.Bytes(), nil
}
