package reader

import (
	"html/template"
	"time"

	bf "github.com/russross/blackfriday"
)

const (
	markdownExtensions = 0 |
		bf.EXTENSION_NO_INTRA_EMPHASIS |
		bf.EXTENSION_TABLES |
		bf.EXTENSION_FENCED_CODE |
		bf.EXTENSION_AUTO_HEADER_IDS |
		bf.EXTENSION_HEADER_IDS |
		bf.EXTENSION_AUTOLINK |
		bf.EXTENSION_STRIKETHROUGH |
		bf.EXTENSION_SPACE_HEADERS |
		bf.EXTENSION_FOOTNOTES |
		bf.EXTENSION_BACKSLASH_LINE_BREAK |
		bf.EXTENSION_DEFINITION_LISTS

	htmlRendererParams = 0 |
		bf.HTML_USE_XHTML |
		bf.HTML_FOOTNOTE_RETURN_LINKS |
		bf.HTML_USE_SMARTYPANTS |
		bf.HTML_SMARTYPANTS_FRACTIONS |
		bf.HTML_SMARTYPANTS_DASHES |
		bf.HTML_SMARTYPANTS_LATEX_DASHES
)

type Story struct {
	Name            string
	IsDraft         bool
	Title           string
	PublicationDate time.Time
	Content         template.HTML
	Tags            []string

	// Extras is a container for any additional data that's supposed to be
	// passed to the templates and doesn't fit into any other field.
	Extras interface{}
}

type metadata struct {
	Title      string      `json:"title"`
	IsDraft    bool        `json:"draft"`
	DateStr    string      `json:"date"`
	Categories []string    `json:"categories"`
	Extras     interface{} `json:"extras"`
}
