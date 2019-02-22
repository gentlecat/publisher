package story

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"
)

const (
	dateFormat         = "2006-Jan-02" // UTC
	markdownFileFormat = "md"

	metadataSeparator = "\n+++\n"

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

type storiesSlice []Story

func (a storiesSlice) Len() int           { return len(a) }
func (a storiesSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a storiesSlice) Less(i, j int) bool { return a[i].PublicationDate.After(a[j].PublicationDate) }

// ReadAll reads all stories in a specified directory and returns a list of
// them. The list is sorted by publication date.
func ReadAll(storiesDir string, skipDrafts bool) (*[]Story, error) {
	stories := make([]Story, 0)

	files, err := ioutil.ReadDir(storiesDir)
	if err != nil {
		return &stories, err
	}
	for _, f := range files {
		if !isStoryFile(f) {
			continue
		}
		log.Printf("Reading file: %s\n", f.Name())
		s, err := read(path.Join(storiesDir, f.Name()))
		// TODO: Consider not failing the whole process, but skipping bad file instead
		if err != nil {
			return &stories, err
		}
		if skipDrafts && s.IsDraft {
			continue
		}
		stories = append(stories, s)
	}
	sort.Sort(storiesSlice(stories))
	return &stories, nil
}

func isStoryFile(f os.FileInfo) bool {
	if f.IsDir() || !strings.HasSuffix(strings.ToLower(f.Name()), "."+markdownFileFormat) {
		return false
	}
	return true
}

func read(storyFilePath string) (s Story, err error) {
	data, err := ioutil.ReadFile(storyFilePath)
	if err != nil {
		return s, err
	}

	parts := strings.SplitN(string(data), metadataSeparator, 2)
	if len(parts) > 2 {
		return s, errors.New("story file hasn't been split up correctly")
	}

	m, err := parseMetadata(parts[0])
	if err != nil {
		return s, errors.New(fmt.Sprint("failed to parse metadata JSON: ", err))
	}
	s.IsDraft = m.IsDraft
	s.Name = clearPath(storyFilePath)
	s.Title = m.Title
	s.Tags = lowerAll(m.Categories)
	s.PublicationDate, err = time.Parse(dateFormat, m.DateStr)
	if err != nil {
		return s, err
	}
	if len(parts) == 2 {
		s.Content = parseContent(parts[1])
	}
	return s, err
}

func parseMetadata(metadataJSON string) (metadata, error) {
	var m metadata
	err := json.Unmarshal([]byte(metadataJSON), &m)
	if err != nil {
		return metadata{}, err
	}
	return m, nil
}

// clearPath removes path and format parts from the story path leaving only its name.
func clearPath(filePath string) string {
	_, file := path.Split(filePath)
	const fileFormatSeparator = "."
	formatParts := strings.Split(file, fileFormatSeparator)
	return strings.Join(formatParts[:len(formatParts)-1], fileFormatSeparator)
}

func lowerAll(strs []string) []string {
	out := make([]string, len(strs))
	for i, v := range strs {
		out[i] = strings.ToLower(v)
	}
	return out
}

// parseContent parses a story in markdown format and converts it to HTML.
func parseContent(content string) template.HTML {
	renderer := blackfriday.HtmlRenderer(commonHtmlFlags, "", "")
	return template.HTML(blackfriday.Markdown([]byte(content), renderer, markdownExtensions))
}
