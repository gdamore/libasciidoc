package parser

import (
	"io"

	"github.com/bytesparadise/libasciidoc/pkg/configuration"
	"github.com/bytesparadise/libasciidoc/pkg/types"
	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

// ParseDocument parses the content of the reader identitied by the filename
func ParseDocument(r io.Reader, config configuration.Configuration) (types.Document, error) {
	draftDoc, err := ParseDraftDocument(r, config)
	if err != nil {
		return types.Document{}, err
	}
	attrs := types.AttributesWithOverrides{
		Content:   types.Attributes{},
		Overrides: config.AttributeOverrides,
	}
	// also, add all front-matter key/values
	attrs.Add(draftDoc.FrontMatter.Content)
	// also, add all AttributeDeclaration at the top of the document
	attrs.Add(draftDoc.Attributes())

	// apply document attribute substitutions and re-parse paragraphs that were affected
	blocks, _, err := applyAttributeSubstitutions(draftDoc.Blocks, attrs)
	if err != nil {
		return types.Document{}, err
	}

	// now, merge list items into proper lists
	blocks, err = rearrangeListItems(blocks.([]interface{}), false)
	if err != nil {
		return types.Document{}, err
	}
	// filter out blocks not needed in the final doc
	blocks = filter(blocks.([]interface{}), allMatchers...)

	blocks, footnotes := processFootnotes(blocks.([]interface{}))
	// now, rearrange elements in a hierarchical manner
	doc := rearrangeSections(blocks.([]interface{}))
	// also, set the footnotes
	doc.Footnotes = footnotes
	// insert the preamble at the right location
	doc = includePreamble(doc)
	// and add all remaining attributes, too
	extraAttrs := attrs.All()
	if doc.Attributes == nil && len(extraAttrs) > 0 {
		doc.Attributes = types.Attributes{}
	}
	doc.Attributes.Add(extraAttrs)
	// also insert the table of contents
	doc = includeTableOfContentsPlaceHolder(doc)
	// finally
	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("final document:")
		spew.Dump(doc)
	}
	return doc, nil
}
