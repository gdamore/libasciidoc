package sgml

import (
	"bytes"
	htmltemplate "html/template"
	"io"
	"strings"
	texttemplate "text/template"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Renderer implements the backend render interface by using sgml.
type Renderer interface {

	// Render renders a document to the given output stream.
	Render(ctx *Context, doc Document, output io.Writer) (Metadata, error)

	// SetFunction sets the named function.
	SetFunction(name string, fn interface{})

	// Templates returns the Templates used by this Renderer.
	// It cannot be altered on a given Renderer, since the old
	// templates may have already been parsed.
	Templates() Templates
}

func NewRenderer(t Templates) Renderer {
	sr := &sgmlRenderer{
		templates: t,
	}
	// Establish some default function handlers.
	sr.functions = funcMap{
		"render":         sr.renderElements,
		"renderElements": sr.renderElements,
		"renderInline":   sr.renderInlineElements,
		"renderList":     sr.renderListElements,
		"renderLines":    sr.renderLines,
		"escape":         EscapeString,
		"renderToC":      sr.renderTableOfContentsSections,
		"renderFootnote": sr.renderFootnote,
		"includeNewline": sr.includeNewline,
		"renderVerse":    sr.renderVerseBlockElement,
		"plainText":      sr.withPlainText,
		"trimRight":      sr.trimRight,
		"trimLeft":       sr.trimLeft,
		"trim":           sr.trimBoth,
	}

	return sr
}

func (sr *sgmlRenderer) trimLeft(s string) string {
	return strings.TrimLeft(s, " ")
}

func (sr *sgmlRenderer) trimRight(s string) string {
	return strings.TrimRight(s, " ")
}

func (sr *sgmlRenderer) trimBoth(s string) string {
	return strings.Trim(s, " ")
}

func (sr *sgmlRenderer) SetFunction(name string, fn interface{}) {
	sr.functions[name] = fn
}

// Templates returns the Templates being used by this renderer.
// A copy is made, as we cannot change the original Templates
// due to it already being used.
func (sr *sgmlRenderer) Templates() Templates {
	return sr.templates
}

func (sr *sgmlRenderer) newTemplate(name string, tmpl string, err error) (*textTemplate, error) {
	// NB: if the data is missing below, it will be an empty string.
	if err != nil {
		return nil, err
	}
	t := texttemplate.New(name)
	t.Funcs(sr.functions)
	t, err = t.Parse(tmpl)
	if err != nil {
		log.Errorf("failed to initialize '%s' template: %v", name, err)
		return nil, err
	}
	return t, nil
}

// Render renders the given document in HTML and writes the result in the given `writer`
func (sr *sgmlRenderer) Render(ctx *Context, doc Document, output io.Writer) (Metadata, error) {

	var md Metadata
	err := sr.prepareTemplates()
	if err != nil {
		return md, err
	}
	renderedTitle, err := sr.renderDocumentTitle(ctx, doc)
	if err != nil {
		return md, errors.Wrapf(err, "unable to render full document")
	}
	// needs to be set before rendering the content elements
	ctx.TableOfContents, err = sr.newTableOfContents(ctx, doc)
	if err != nil {
		return md, errors.Wrapf(err, "unable to render full document")
	}
	renderedHeader, renderedContent, err := sr.splitAndRender(ctx, doc)
	if err != nil {
		return md, errors.Wrapf(err, "unable to render full document")
	}

	if ctx.Config.IncludeHeaderFooter {
		log.Debugf("Rendering full document...")
		err = sr.article.Execute(output, struct {
			Generator     string
			Doctype       string
			Title         string
			Authors       string
			Header        string
			Role          string
			Content       sanitized
			RevNumber     string
			LastUpdated   string
			CSS           string
			IncludeHeader bool
			IncludeFooter bool
		}{
			Generator:     "libasciidoc", // TODO: externalize this value and include the lib version ?
			Doctype:       doc.Attributes.GetAsStringWithDefault(types.AttrDocType, "article"),
			Title:         string(renderedTitle),
			Authors:       sr.renderAuthors(doc),
			Header:        string(renderedHeader),
			Role:          documentRole(doc),
			Content:       sanitized(renderedContent), //nolint: gosec
			RevNumber:     doc.Attributes.GetAsStringWithDefault("revnumber", ""),
			LastUpdated:   ctx.Config.LastUpdated.Format(configuration.LastUpdatedFormat),
			CSS:           ctx.Config.CSS,
			IncludeHeader: !doc.Attributes.Has(types.AttrNoHeader),
			IncludeFooter: !doc.Attributes.Has(types.AttrNoFooter),
		})
		if err != nil {
			return md, errors.Wrapf(err, "unable to render full document")
		}
	} else {
		_, err = output.Write(renderedContent)
		if err != nil {
			return md, errors.Wrapf(err, "unable to render full document")
		}
	}
	// generate the metadata to be returned to the caller
	md.Title = string(renderedTitle)
	// arguably this should be a time.Time for use in Go
	md.LastUpdated = ctx.Config.LastUpdated.Format(configuration.LastUpdatedFormat)
	md.TableOfContents = ctx.TableOfContents
	return md, err
}

// splitAndRender the document with the header elements on one side
// and all other elements (table of contents, with preamble, content) on the other side,
// then renders the header and other elements
func (sr *sgmlRenderer) splitAndRender(ctx *Context, doc Document) ([]byte, []byte, error) {
	switch doc.Attributes.GetAsStringWithDefault(types.AttrDocType, "article") {
	case "manpage":
		return sr.splitAndRenderForManpage(ctx, doc)
	default:
		return sr.splitAndRenderForArticle(ctx, doc)
	}
}

// splits the document with the title of the section 0 (if available) on one side
// and all other elements (table of contents, with preamble, content) on the other side
func (sr *sgmlRenderer) splitAndRenderForArticle(ctx *Context, doc Document) ([]byte, []byte, error) {
	if ctx.Config.IncludeHeaderFooter {
		if header, found := doc.Header(); found {
			renderedHeader, err := sr.renderArticleHeader(ctx, header)
			if err != nil {
				return nil, nil, err
			}
			renderedContent, err := sr.renderDocumentElements(ctx, header.Elements, doc.Footnotes)
			if err != nil {
				return nil, nil, err
			}
			return renderedHeader, renderedContent, nil
		}
	}
	renderedContent, err := sr.renderDocumentElements(ctx, doc.Elements, doc.Footnotes)
	if err != nil {
		return nil, nil, err
	}
	return []byte{}, renderedContent, nil
}

// splits the document with the header elements on one side
// and the other elements (table of contents, with preamble, content) on the other side
func (sr *sgmlRenderer) splitAndRenderForManpage(ctx *Context, doc Document) ([]byte, []byte, error) {
	header, _ := doc.Header()
	nameSection := header.Elements[0].(types.Section)

	if ctx.Config.IncludeHeaderFooter {
		renderedHeader, err := sr.renderManpageHeader(ctx, header, nameSection)
		if err != nil {
			return nil, nil, err
		}
		renderedContent, err := sr.renderDocumentElements(ctx, header.Elements[1:], doc.Footnotes)
		if err != nil {
			return nil, nil, err
		}
		return renderedHeader, renderedContent, nil
	}
	// in that case, we still want to display the name section
	renderedHeader, err := sr.renderManpageHeader(ctx, types.Section{}, nameSection)
	if err != nil {
		return nil, nil, err
	}
	renderedContent, err := sr.renderDocumentElements(ctx, header.Elements[1:], doc.Footnotes)
	if err != nil {
		return nil, nil, err
	}
	result := &bytes.Buffer{}
	result.Write(renderedHeader)
	result.WriteString("\n")
	result.Write(renderedContent)
	return []byte{}, result.Bytes(), nil
}

func documentRole(doc Document) string {
	if header, found := doc.Header(); found {
		return header.Attributes.GetAsStringWithDefault(types.AttrRole, "")
	}
	return ""
}

func (sr *sgmlRenderer) renderAuthors(doc Document) string {
	authors, found := doc.Authors()
	if !found {
		return ""
	}
	authorStrs := make([]string, len(authors))
	for i, author := range authors {
		authorStrs[i] = author.FullName
	}
	return strings.Join(authorStrs, "; ")
}

func (sr *sgmlRenderer) renderDocumentTitle(ctx *Context, doc Document) ([]byte, error) {
	if header, found := doc.Header(); found {
		title, err := sr.renderPlainText(ctx, header.Title)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to render document title")
		}
		return title, nil
	}
	return nil, nil
}

func (sr *sgmlRenderer) renderArticleHeader(ctx *Context, header types.Section) ([]byte, error) {
	renderedHeader, err := sr.renderInlineElements(ctx, header.Title)
	if err != nil {
		return nil, err
	}
	documentDetails, err := sr.renderDocumentDetails(ctx)
	if err != nil {
		return nil, err
	}

	output := &bytes.Buffer{}
	err = sr.articleHeader.Execute(output, struct {
		Header  string
		Details *htmltemplate.HTML // TODO: convert to sanitized (no need to be a pointer)
	}{
		Header:  string(renderedHeader),
		Details: documentDetails,
	})
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func (sr *sgmlRenderer) renderManpageHeader(ctx *Context, header types.Section, nameSection types.Section) ([]byte, error) {
	renderedHeader, err := sr.renderInlineElements(ctx, header.Title)
	if err != nil {
		return nil, err
	}
	renderedName, err := sr.renderInlineElements(ctx, nameSection.Title)
	if err != nil {
		return nil, err
	}
	description := nameSection.Elements[0].(types.Paragraph) // TODO: type check
	if description.Attributes == nil {
		description.Attributes = types.Attributes{}
	}
	description.Attributes.AddNonEmpty(types.AttrKind, "manpage")
	renderedContent, err := sr.renderParagraph(ctx, description)
	if err != nil {
		return nil, err
	}
	output := &bytes.Buffer{}
	err = sr.manpageHeader.Execute(output, struct {
		Header    string
		Name      string
		Content   sanitized
		IncludeH1 bool
	}{
		Header:    string(renderedHeader),
		Name:      string(renderedName),
		Content:   sanitized(renderedContent), //nolint: gosec
		IncludeH1: len(renderedHeader) > 0,
	})
	if err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

// renderDocumentElements renders all document elements, including the footnotes,
// but not the HEAD and BODY containers
func (sr *sgmlRenderer) renderDocumentElements(ctx *Context, source []interface{}, footnotes []types.Footnote) ([]byte, error) {
	elements := []interface{}{}
	for i, e := range source {
		switch e := e.(type) {
		case types.Preamble:
			if !e.HasContent() { // why !HasContent ???
				// retain the preamble
				elements = append(elements, e)
				continue
			}
			// retain everything "as-is"
			elements = source
		case types.Section:
			if e.Level == 0 {
				// retain the section's elements...
				elements = append(elements, e.Elements)
				// ... and add the other elements (in case there's another section 0...)
				elements = append(elements, source[i+1:]...)
				continue
			}
			// retain everything "as-is"
			elements = source
		default:
			// retain everything "as-is"
			elements = source
		}
	}
	buff := &bytes.Buffer{}
	renderedElements, err := sr.renderElements(ctx, elements)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to render document elements")
	}
	buff.Write(renderedElements)
	renderedFootnotes, err := sr.renderFootnotes(ctx, footnotes)
	if err != nil {
		return []byte{}, errors.Wrapf(err, "failed to render document elements")
	}
	buff.Write(renderedFootnotes)
	return buff.Bytes(), nil
}
