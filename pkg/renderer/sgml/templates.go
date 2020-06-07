package sgml

import "sync"

type TemplateName string

// Templates represents all the templates we use.
type Templates struct {
	AdmonitionBlock         string
	AdmonitionParagraph     string
	Article                 string
	ArticleHeader           string
	BlankLine               string
	BlockImage              string
	BoldText                string
	CalloutList             string
	DelimitedBlockParagraph string
	DocumentDetails         string
	DocumentAuthorDetails   string
	ExampleBlock            string
	ExternalCrossReference  string
	FencedBlock             string
	Footnote                string
	FootnoteRef             string
	FootnoteRefPlain        string
	Footnotes               string
	InlineImage             string
	InternalCrossReference  string
	InvalidFootnote         string
	ItalicText              string
	LabeledList             string
	LabeledListHorizontal   string
	LineBreak               string
	Link                    string
	ListingBlock            string
	LiteralBlock            string
	ManpageHeader           string
	ManpageNameParagraph    string
	MonospaceText           string
	OrderedList             string
	Paragraph               string
	PassthroughBlock        string
	Preamble                string
	QAndAList               string
	QuoteBlock              string
	QuoteParagraph          string
	SectionContent          string
	SectionHeader           string
	SectionOne              string
	SidebarBlock            string
	SourceBlock             string
	SourceBlockContent      string
	SourceParagraph         string
	StringElement           string
	SubscriptText           string
	SuperscriptText         string
	Table                   string
	TocRoot                 string
	TocSection              string
	UnorderedList           string
	VerbatimLine            string
	VerseBlock              string
	VerseBlockParagraph     string
	VerseParagraph          string
}

type sgmlRenderer struct {
	functions funcMap
	templates *Templates
	prepared  sync.Once

	// Processed templates
	admonitionBlock         *textTemplate
	admonitionParagraph     *textTemplate
	article                 *textTemplate
	articleHeader           *textTemplate
	blankLine               *textTemplate
	blockImage              *textTemplate
	boldText                *textTemplate
	calloutList             *textTemplate
	delimitedBlockParagraph *textTemplate
	documentDetails         *textTemplate
	documentAuthorDetails   *textTemplate
	externalCrossReference  *textTemplate
	exampleBlock            *textTemplate
	fencedBlock             *textTemplate
	footnote                *textTemplate
	footnoteRef             *textTemplate
	footnoteRefPlain        *textTemplate
	footnotes               *textTemplate
	inlineImage             *textTemplate
	internalCrossReference  *textTemplate
	invalidFootnote         *textTemplate
	italicText              *textTemplate
	labeledList             *textTemplate
	labeledListHorizontal   *textTemplate
	lineBreak               *textTemplate
	link                    *textTemplate
	listingBlock            *textTemplate
	literalBlock            *textTemplate
	manpageHeader           *textTemplate
	manpageNameParagraph    *textTemplate
	monospaceText           *textTemplate
	orderedList             *textTemplate
	paragraph               *textTemplate
	passthroughBlock        *textTemplate
	preamble                *textTemplate
	qAndAList               *textTemplate
	quoteBlock              *textTemplate
	quoteParagraph          *textTemplate
	sectionContent          *textTemplate
	sectionHeader           *textTemplate
	sectionOne              *textTemplate
	sidebarBlock            *textTemplate
	sourceBlock             *textTemplate
	sourceBlockContent      *textTemplate
	sourceParagraph         *textTemplate
	stringElement           *textTemplate
	subscriptText           *textTemplate
	superscriptText         *textTemplate
	table                   *textTemplate
	tocRoot                 *textTemplate
	tocSection              *textTemplate
	unorderedList           *textTemplate
	verbatimLine            *textTemplate
	verseBlock              *textTemplate
	verseBlockParagraph     *textTemplate
	verseParagraph          *textTemplate
}

func (sr *sgmlRenderer) prepareTemplates() error {
	t := sr.templates
	var err error

	sr.prepared.Do(func() {

		sr.admonitionBlock, err = sr.newTemplate("admonition-block", t.AdmonitionBlock, err)
		sr.admonitionParagraph, err = sr.newTemplate("admonition-paragraph", t.AdmonitionParagraph, err)
		sr.article, err = sr.newTemplate("article", t.Article, err)
		sr.articleHeader, err = sr.newTemplate("article-header", t.ArticleHeader, err)
		sr.blankLine, err = sr.newTemplate("blank-line", t.BlankLine, err)
		sr.blockImage, err = sr.newTemplate("block-image", t.BlockImage, err)
		sr.boldText, err = sr.newTemplate("bold-text", t.BoldText, err)
		sr.calloutList, err = sr.newTemplate("callout-list", t.CalloutList, err)
		sr.delimitedBlockParagraph, err = sr.newTemplate("delimited-block-paragraph", t.DelimitedBlockParagraph, err)
		sr.documentDetails, err = sr.newTemplate("document-details", t.DocumentDetails, err)
		sr.documentAuthorDetails, err = sr.newTemplate("document-author-details", t.DocumentAuthorDetails, err)
		sr.exampleBlock, err = sr.newTemplate("example-block", t.ExampleBlock, err)
		sr.externalCrossReference, err = sr.newTemplate("external-xref", t.ExternalCrossReference, err)
		sr.fencedBlock, err = sr.newTemplate("fenced-block", t.FencedBlock, err)
		sr.footnote, err = sr.newTemplate("footnote", t.Footnote, err)
		sr.footnoteRef, err = sr.newTemplate("footnote-ref", t.FootnoteRef, err)
		sr.footnoteRefPlain, err = sr.newTemplate("footnote-ref-plain", t.FootnoteRefPlain, err)
		sr.footnotes, err = sr.newTemplate("footnotes", t.Footnotes, err)
		sr.inlineImage, err = sr.newTemplate("inline-image", t.InlineImage, err)
		sr.internalCrossReference, err = sr.newTemplate("internal-xref", t.InternalCrossReference, err)
		sr.invalidFootnote, err = sr.newTemplate("invalid-footnote", t.InvalidFootnote, err)
		sr.italicText, err = sr.newTemplate("italic-text", t.ItalicText, err)
		sr.labeledList, err = sr.newTemplate("labeled-list", t.LabeledList, err)
		sr.labeledListHorizontal, err = sr.newTemplate("labeled-list-horizontal", t.LabeledListHorizontal, err)
		sr.lineBreak, err = sr.newTemplate("line-break", t.LineBreak, err)
		sr.link, err = sr.newTemplate("link", t.Link, err)
		sr.listingBlock, err = sr.newTemplate("listing", t.ListingBlock, err)
		sr.literalBlock, err = sr.newTemplate("literal-block", t.LiteralBlock, err)
		sr.manpageHeader, err = sr.newTemplate("manpage-header", t.ManpageHeader, err)
		sr.manpageNameParagraph, err = sr.newTemplate("manpage-name-paragraph", t.ManpageNameParagraph, err)
		sr.monospaceText, err = sr.newTemplate("monospace-text", t.MonospaceText, err)
		sr.orderedList, err = sr.newTemplate("ordered-list", t.OrderedList, err)
		sr.paragraph, err = sr.newTemplate("paragraph", t.Paragraph, err)
		sr.passthroughBlock, err = sr.newTemplate("passthrough", t.PassthroughBlock, err)
		sr.preamble, err = sr.newTemplate("preamble", t.Preamble, err)
		sr.qAndAList, err = sr.newTemplate("qanda-block", t.QAndAList, err)
		sr.quoteBlock, err = sr.newTemplate("quote-block", t.QuoteBlock, err)
		sr.quoteParagraph, err = sr.newTemplate("quote-paragraph", t.QuoteParagraph, err)
		sr.sectionContent, err = sr.newTemplate("section-content", t.SectionContent, err)
		sr.sectionHeader, err = sr.newTemplate("section-header", t.SectionHeader, err)
		sr.sectionOne, err = sr.newTemplate("section-one", t.SectionOne, err)
		sr.stringElement, err = sr.newTemplate("string-element", t.StringElement, err)
		sr.sidebarBlock, err = sr.newTemplate("sidebar-block", t.SidebarBlock, err)
		sr.sourceBlock, err = sr.newTemplate("source-block", t.SourceBlock, err)
		sr.sourceBlockContent, err = sr.newTemplate("source-block-content", t.SourceBlockContent, err)
		sr.sourceParagraph, err = sr.newTemplate("source-paragraph", t.SourceParagraph, err)
		sr.subscriptText, err = sr.newTemplate("subscript", t.SubscriptText, err)
		sr.superscriptText, err = sr.newTemplate("superscript", t.SuperscriptText, err)
		sr.table, err = sr.newTemplate("table", t.Table, err)
		sr.tocRoot, err = sr.newTemplate("toc-root", t.TocRoot, err)
		sr.tocSection, err = sr.newTemplate("toc-section", t.TocSection, err)
		sr.unorderedList, err = sr.newTemplate("unordered-list", t.UnorderedList, err)
		sr.verbatimLine, err = sr.newTemplate("verbatim-line", t.VerbatimLine, err)
		sr.verseBlock, err = sr.newTemplate("verse", t.VerseBlock, err)
		sr.verseBlockParagraph, err = sr.newTemplate("verse-block-paragraph", t.VerseBlockParagraph, err)
		sr.verseParagraph, err = sr.newTemplate("verse-paragraph", t.VerseParagraph, err)
	})

	return err
}
