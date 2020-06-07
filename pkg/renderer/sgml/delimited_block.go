package sgml

import (
	"bytes"
	"strconv"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (sr *sgmlRenderer) renderDelimitedBlock(ctx *Context, b types.DelimitedBlock) ([]byte, error) {
	log.Debugf("rendering delimited block of kind '%v'", b.Attributes[types.AttrKind])
	var err error
	kind := b.Kind
	switch kind {
	case types.Fenced:
		return sr.renderFencedBlock(ctx, b)
	case types.Listing:
		return sr.renderListingBlock(ctx, b)
	case types.Source:
		return sr.renderSourceBlock(ctx, b)
	case types.Example:
		return sr.renderExampleBlock(ctx, b)
	case types.Quote, types.MarkdownQuote:
		return sr.renderQuoteBlock(ctx, b)
	case types.Verse:
		return sr.renderVerseBlock(ctx, b)
	case types.Sidebar:
		return sr.renderSidebarBlock(ctx, b)
	case types.Passthrough:
		return sr.renderPassthrough(ctx, b)
	default:
		return nil, errors.Wrapf(err, "unable to render delimited block")
	}
}

func (sr *sgmlRenderer) renderFencedBlock(ctx *Context, b types.DelimitedBlock) ([]byte, error) {
	previousWithinDelimitedBlock := ctx.WithinDelimitedBlock
	previousIncludeBlankLine := ctx.IncludeBlankLine
	defer func() {
		ctx.WithinDelimitedBlock = previousWithinDelimitedBlock
		ctx.IncludeBlankLine = previousIncludeBlankLine
	}()
	ctx.WithinDelimitedBlock = true
	ctx.IncludeBlankLine = true
	result := &bytes.Buffer{}
	err := sr.fencedBlock.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			ID       string
			Title    string
			Elements []interface{}
		}{
			ID:       sr.renderElementID(b.Attributes),
			Title:    sr.renderElementTitle(b.Attributes),
			Elements: discardTrailingBlankLines(b.Elements),
		},
	})
	return result.Bytes(), err
}

func (sr *sgmlRenderer) renderListingBlock(ctx *Context, b types.DelimitedBlock) ([]byte, error) {
	previousWithinDelimitedBlock := ctx.WithinDelimitedBlock
	previousIncludeBlankLine := ctx.IncludeBlankLine
	defer func() {
		ctx.WithinDelimitedBlock = previousWithinDelimitedBlock
		ctx.IncludeBlankLine = previousIncludeBlankLine
	}()
	ctx.WithinDelimitedBlock = true
	ctx.IncludeBlankLine = true
	result := &bytes.Buffer{}
	err := sr.listingBlock.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			ID       string
			Title    string
			Elements []interface{}
		}{
			ID:       sr.renderElementID(b.Attributes),
			Title:    sr.renderElementTitle(b.Attributes),
			Elements: discardTrailingBlankLines(b.Elements),
		},
	})
	return result.Bytes(), err
}

func (sr *sgmlRenderer) renderSourceBlock(ctx *Context, b types.DelimitedBlock) ([]byte, error) {
	previousWithinDelimitedBlock := ctx.WithinDelimitedBlock
	previousIncludeBlankLine := ctx.IncludeBlankLine
	defer func() {
		ctx.WithinDelimitedBlock = previousWithinDelimitedBlock
		ctx.IncludeBlankLine = previousIncludeBlankLine
	}()
	ctx.WithinDelimitedBlock = true
	ctx.IncludeBlankLine = true
	// first, render the content
	contentBuf := &bytes.Buffer{}
	err := sr.sourceBlockContent.Execute(contentBuf, ContextualPipeline{
		Context: ctx,
		Data: struct {
			Elements []interface{}
		}{
			Elements: discardTrailingBlankLines(b.Elements),
		}})
	if err != nil {
		return []byte{}, err
	}
	content := contentBuf.String()

	highlighter, _ := ctx.Attributes.GetAsString(types.AttrSyntaxHighlighter)
	language, found := b.Attributes.GetAsString(types.AttrLanguage)
	if found && highlighter == "pygments" {
		// using github.com/alecthomas/chroma to highlight the content
		contentBuf = &bytes.Buffer{}
		lexer := lexers.Get(language)
		lexer = chroma.Coalesce(lexer)
		style := styles.Fallback
		if s, found := ctx.Attributes.GetAsString("pygments-style"); found {
			style = styles.Get(s)
		}
		iterator, err := lexer.Tokenise(nil, content)
		if err != nil {
			return []byte{}, err
		}
		options := []html.Option{
			html.ClassPrefix("tok-"),
			html.PreventSurroundingPre(true),
		}
		// extra option: inline CSS instead of classes
		if ctx.Attributes.GetAsStringWithDefault("pygments-css", "classes") == "style" {
			options = append(options, html.WithClasses(false))
		} else {
			options = append(options, html.WithClasses(true))
		}
		// extra option: line numbers
		if b.Attributes.Has(types.AttrLineNums) {
			options = append(options, html.WithLineNumbers(true))
		}
		err = html.New(options...).Format(contentBuf, style, iterator)
		if err != nil {
			return []byte{}, err
		}
		content = contentBuf.String()
	}

	result := &bytes.Buffer{}
	err = sr.sourceBlock.Execute(result, struct {
		ID                string
		Title             string
		Language          string
		SyntaxHighlighter string
		Content           string
	}{
		ID:                sr.renderElementID(b.Attributes),
		Title:             sr.renderElementTitle(b.Attributes),
		SyntaxHighlighter: highlighter,
		Language:          language,
		Content:           content,
	})
	return result.Bytes(), err
}

func (sr *sgmlRenderer) renderExampleBlock(ctx *Context, b types.DelimitedBlock) ([]byte, error) {
	result := &bytes.Buffer{}
	if k, ok := b.Attributes[types.AttrAdmonitionKind].(types.AdmonitionKind); ok {
		err := sr.admonitionBlock.Execute(result, ContextualPipeline{
			Context: ctx,
			Data: struct {
				ID        string
				Class     string
				IconClass string
				IconTitle string
				Title     string
				Elements  []interface{}
			}{
				ID:        sr.renderElementID(b.Attributes),
				Class:     renderClass(k),
				IconClass: renderIconClass(ctx, k),
				IconTitle: renderIconTitle(k),
				Title:     sr.renderElementTitle(b.Attributes),
				Elements:  discardTrailingBlankLines(b.Elements),
			},
		})
		return result.Bytes(), err
	}
	// default, example block
	var title string
	if b.Attributes.Has(types.AttrTitle) {
		title = "Example " + strconv.Itoa(ctx.GetAndIncrementExampleBlockCounter()) + ". " + sr.renderElementTitle(b.Attributes)
	}
	err := sr.exampleBlock.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			ID       string
			Title    string
			Elements []interface{}
		}{
			ID:       sr.renderElementID(b.Attributes),
			Title:    title,
			Elements: discardTrailingBlankLines(b.Elements),
		},
	})
	return result.Bytes(), err
}

func (sr *sgmlRenderer) renderQuoteBlock(ctx *Context, b types.DelimitedBlock) ([]byte, error) {
	result := &bytes.Buffer{}
	err := sr.quoteBlock.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			ID          string
			Title       string
			Attribution Attribution
			Elements    []interface{}
		}{
			ID:          sr.renderElementID(b.Attributes),
			Title:       sr.renderElementTitle(b.Attributes),
			Attribution: newDelimitedBlockAttribution(b),
			Elements:    b.Elements,
		},
	})
	return result.Bytes(), err
}

func (sr *sgmlRenderer) renderVerseBlock(ctx *Context, b types.DelimitedBlock) ([]byte, error) {
	result := &bytes.Buffer{}
	err := sr.verseBlock.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			ID          string
			Title       string
			Attribution Attribution
			Elements    []interface{}
		}{
			ID:          sr.renderElementID(b.Attributes),
			Title:       sr.renderElementTitle(b.Attributes),
			Attribution: newDelimitedBlockAttribution(b),
			Elements:    discardTrailingBlankLines(b.Elements),
		},
	})
	return result.Bytes(), err
}

func (sr *sgmlRenderer) renderVerseBlockElement(ctx *Context, element interface{}) ([]byte, error) {
	previousIncludeBlankLine := ctx.IncludeBlankLine
	defer func() {
		ctx.IncludeBlankLine = previousIncludeBlankLine
	}()
	ctx.IncludeBlankLine = true
	switch e := element.(type) {
	case types.Paragraph:
		return sr.renderVerseBlockParagraph(ctx, e)
	case types.BlankLine:
		return sr.renderBlankLine(ctx, e)
	default:
		return nil, errors.Errorf("unexpected type of element to include in verse block: %T", element)
	}
}

func (sr *sgmlRenderer) renderVerseBlockParagraph(ctx *Context, p types.Paragraph) ([]byte, error) {
	log.Debugf("rendering paragraph with %d line(s) within a delimited block or a list", len(p.Lines))
	result := &bytes.Buffer{}
	err := sr.verseBlockParagraph.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			Lines [][]interface{}
		}{
			Lines: p.Lines,
		},
	})
	return result.Bytes(), err
}

func (sr *sgmlRenderer) renderSidebarBlock(ctx *Context, b types.DelimitedBlock) ([]byte, error) {
	result := &bytes.Buffer{}
	err := sr.sidebarBlock.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			ID       string
			Title    string
			Elements []interface{}
		}{
			ID:       sr.renderElementID(b.Attributes),
			Title:    sr.renderElementTitle(b.Attributes),
			Elements: discardTrailingBlankLines(b.Elements),
		},
	})
	return result.Bytes(), err
}

func (sr *sgmlRenderer) renderPassthrough(ctx *Context, b types.DelimitedBlock) ([]byte, error) {
	result := &bytes.Buffer{}
	err := sr.passthroughBlock.Execute(result, ContextualPipeline{
		Context: ctx,
		Data: struct {
			ID       string
			Elements []interface{}
		}{
			ID:       sr.renderElementID(b.Attributes),
			Elements: discardTrailingBlankLines(b.Elements),
		},
	})
	return result.Bytes(), err
}

func discardTrailingBlankLines(elements []interface{}) []interface{} {
	// discard blank elements at the end
	log.Debugf("discarding trailing blank lines on %d elements...", len(elements))
	filteredElements := make([]interface{}, len(elements))
	copy(filteredElements, elements)

	for {
		if len(filteredElements) == 0 {
			break
		}
		if l, ok := filteredElements[len(filteredElements)-1].(types.VerbatimLine); ok && l.IsEmpty() {
			log.Debugf("element of type '%T' at position %d is a blank line, discarding it", filteredElements[len(filteredElements)-1], len(filteredElements)-1)
			// remove last element of the slice since it's a blank line
			filteredElements = filteredElements[:len(filteredElements)-1]
		} else if _, ok := filteredElements[len(filteredElements)-1].(types.BlankLine); ok {
			log.Debugf("element of type '%T' at position %d is a blank line, discarding it", filteredElements[len(filteredElements)-1], len(filteredElements)-1)
			// remove last element of the slice since it's a blank line
			filteredElements = filteredElements[:len(filteredElements)-1]
		} else {
			break
		}
	}
	log.Debugf("returning %d elements", len(filteredElements))
	return filteredElements
}
