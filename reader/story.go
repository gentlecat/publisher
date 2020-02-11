package reader

import (
	"html/template"
	"time"

	bf "github.com/russross/blackfriday"
)

const (
	markdownExtensions = 0 |
		bf.NoIntraEmphasis |
		bf.Tables |
		bf.FencedCode |
		bf.AutoHeadingIDs |
		bf.Autolink |
		bf.Strikethrough |
		bf.SpaceHeadings |
		bf.HeadingIDs |
		bf.Footnotes |
		bf.BackslashLineBreak |
		bf.DefinitionLists

	htmlRendererParams = 0 |
		bf.UseXHTML |
		bf.FootnoteReturnLinks |
		bf.Smartypants |
		bf.SmartypantsFractions |
		bf.SmartypantsDashes |
		bf.SmartypantsLatexDashes
)

type Story struct {
	Name            string
	IsDraft         bool
	Title           string
	PublicationDate time.Time
	Content         template.HTML
	Tags            []string
}

type metadata struct {
	Title      string   `json:"title"`
	IsDraft    bool     `json:"draft"`
	DateStr    string   `json:"date"`
	Categories []string `json:"categories"`
}
