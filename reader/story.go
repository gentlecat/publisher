package reader

import (
	"fmt"
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

// Status represents the publication state of a story.
type Status string

const (
	// StatusPublished stories are included everywhere: listings (such as the
	// index page), the RSS feed, the sitemap, and their own page. This is the
	// default when no status is specified.
	StatusPublished Status = "published"

	// StatusUnlisted stories get their own page but are kept out of listings,
	// the RSS feed, and the sitemap. They're only reachable via a direct link.
	StatusUnlisted Status = "unlisted"

	// StatusDraft stories are only built when the generator is running in draft
	// mode. In production builds they're skipped entirely.
	StatusDraft Status = "draft"
)

type Story struct {
	Name            string
	Status          Status
	Title           string
	PublicationDate time.Time
	Content         template.HTML
	Category        string
	Tags            []string

	// SourceDir is the path to the story's source directory (its bundle). It
	// holds the content and metadata files alongside any static assets (such as
	// images) that should be published next to the story.
	SourceDir string

	// Extras is a container for any additional data that's supposed to be
	// passed to the templates and doesn't fit into any other field.
	Extras interface{}
}

// IsListed reports whether the story should appear in listings such as the
// index page, the RSS feed, and the sitemap.
func (s Story) IsListed() bool {
	return s.Status != StatusUnlisted
}

type metadata struct {
	Title    string   `json:"title"`
	Status   Status   `json:"status"`
	DateStr  string   `json:"date"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`

	// IsDraft is the legacy way of marking a story as a draft.
	//
	// Deprecated: use Status ("draft") instead. It's only kept around so that
	// content that hasn't been migrated yet keeps working.
	IsDraft bool `json:"draft"`

	Extras interface{} `json:"extras"`
}

// resolveStatus returns the story status described by the metadata, defaulting
// to published and falling back to the legacy draft flag when no status is set.
func (m metadata) resolveStatus() (Status, error) {
	switch m.Status {
	case "":
		if m.IsDraft {
			return StatusDraft, nil
		}
		return StatusPublished, nil
	case StatusPublished, StatusUnlisted, StatusDraft:
		return m.Status, nil
	default:
		return "", fmt.Errorf("unknown status %q", m.Status)
	}
}
