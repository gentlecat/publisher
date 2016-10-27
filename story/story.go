package story

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

const (
	dateFormat         = "2006-Jan-02" // UTC
	markdownFileFormat = "md"

	markdownExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTO_HEADER_IDS |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_FOOTNOTES |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS
	commonHtmlFlags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_FOOTNOTE_RETURN_LINKS |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
)

type Story struct {
	Title   string
	Date    time.Time
	Content template.HTML
}

type metadata struct {
	Stories []storyMetadata `json:"stories"`
}

type storyMetadata struct {
	Name      string `json:"name"`
	Title     string `json:"title"`
	DateStr   string `json:"date"`
	Published bool   `json:"published"`
}

func ReadStories(metadataFile, storiesDir string) map[string]Story {
	stories := make(map[string]Story)

	metadataJSON, err := ioutil.ReadFile(metadataFile)
	check(err)
	var metadata metadata
	err = json.Unmarshal(metadataJSON, &metadata)
	check(err)

	for _, story := range metadata.Stories {
		storyContentPath := filepath.Join(storiesDir, story.Name+"."+markdownFileFormat)
		if _, err := os.Stat(storyContentPath); os.IsNotExist(err) {
			log.Fatalf("Can't find story: %s", story.Name)
		}
		if story.Published {
			date, err := time.Parse(dateFormat, story.DateStr)
			check(err)
			stories[story.Name] = Story{
				Title:   story.Title,
				Date:    date,
				Content: parseStoryContent(storyContentPath),
			}
		}
	}

	return stories
}

// ReadStory parses a story in markdown format and converts it to HTML.
func parseStoryContent(filePath string) template.HTML {
	data, err := ioutil.ReadFile(filePath)
	check(err)

	renderer := blackfriday.HtmlRenderer(commonHtmlFlags, "", "")
	unsafe := blackfriday.Markdown(data, renderer, markdownExtensions)
	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	return template.HTML(policy.SanitizeBytes(unsafe))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
