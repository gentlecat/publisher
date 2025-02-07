package reader

import (
	"html/template"
	"time"

	bf "github.com/russross/blackfriday/v2"
)

const (
	markdownExtensions = bf.AutoHeadingIDs |
		bf.Autolink |
		bf.BackslashLineBreak |
		bf.DefinitionLists |
		bf.FencedCode |
		bf.Footnotes |
		bf.HeadingIDs |
		bf.NoIntraEmphasis |
		bf.SpaceHeadings |
		bf.Strikethrough |
		bf.Tables

	htmlRendererFlags = bf.FootnoteReturnLinks |
		bf.NofollowLinks |
		bf.NoopenerLinks |
		bf.NoreferrerLinks |
		bf.Smartypants |
		bf.SmartypantsDashes |
		bf.SmartypantsFractions |
		bf.SmartypantsLatexDashes |
		bf.UseXHTML
)

var (
	renderer = bf.NewHTMLRenderer(bf.HTMLRendererParameters{
		Flags: htmlRendererFlags,
	})
)

type Story struct {
	Name            string
	IsDraft         bool
	Title           string
	PublicationDate time.Time
	Content         template.HTML
	Category        string
	Tags            []string

	// Extras is a container for any additional data that's supposed to be
	// passed to the templates and doesn't fit into any other field.
	Extras interface{}
}

type metadata struct {
	Title    string      `json:"title"`
	IsDraft  bool        `json:"draft"`
	DateStr  string      `json:"date"`
	Category string      `json:"category"`
	Tags     []string    `json:"tags"`
	Extras   interface{} `json:"extras"`
}
