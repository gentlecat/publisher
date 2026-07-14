package reader

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	bf "github.com/russross/blackfriday/v2"
)

// ContentFileName is the name of the Markdown file holding a story's content
// inside its bundle directory.
const ContentFileName = "index.md"

// MetadataFileName is the name of the JSON file holding a story's metadata
// inside its bundle directory.
const MetadataFileName = "metadata.json"

type Configuration struct {
	// DateFormat is the format used to parse creation date of a story.
	DateFormat string

	// SkipDrafts indicates whether stories with a "draft" status should be
	// ignored.
	//
	// Useful for separation between prod and local environments.
	SkipDrafts bool
}

type StoryReader interface {
	// ReadAll reads (parses) all stories in a particular directory.
	ReadAll(storiesDir string) (*[]Story, error)
}

func NewReader() *Configuration {
	return &Configuration{
		DateFormat: "2006-Jan-02",
	}
}

// ReadAll reads all stories in a specified directory and returns a list of
// them. Each story is a bundle: a subdirectory containing a metadata file, a
// content file, and any static assets. The list is sorted by publication date.
func (r *Configuration) ReadAll(storiesDir string) (*[]Story, error) {
	stories := make([]Story, 0)

	entries, err := os.ReadDir(storiesDir)
	if err != nil {
		return &stories, err
	}
	for _, e := range entries {
		storyDir := path.Join(storiesDir, e.Name())
		if !r.isStoryDir(e, storyDir) {
			continue
		}
		fmt.Printf("  - %s\n", e.Name())
		s, err := r.read(storyDir)
		if err != nil {
			log.Printf("Failed to read story: %s. Error: %s\n", e.Name(), err)
			continue
		}
		if r.SkipDrafts && s.Status == StatusDraft {
			continue
		}
		stories = append(stories, s)
	}
	sort.Sort(storiesSlice(stories))
	return &stories, nil
}

// isStoryDir reports whether the given directory entry is a story bundle. A
// bundle is a directory containing a metadata file. Anything else (loose files
// or directories used purely for grouping, such as an archive) is ignored.
func (r *Configuration) isStoryDir(e os.DirEntry, storyDir string) bool {
	if !e.IsDir() {
		return false
	}
	_, err := os.Stat(path.Join(storyDir, MetadataFileName))
	return err == nil
}

func (r *Configuration) read(storyDir string) (s Story, err error) {
	metadataData, err := os.ReadFile(path.Join(storyDir, MetadataFileName))
	if err != nil {
		return s, err
	}

	m, err := parseMetadata(string(metadataData))
	if err != nil {
		return s, errors.New(fmt.Sprint("failed to parse metadata JSON: ", err))
	}

	s.Status, err = m.resolveStatus()
	if err != nil {
		return s, err
	}
	s.Name = path.Base(storyDir)
	s.SourceDir = storyDir
	s.Title = m.Title
	s.Category = m.Category
	s.Tags = lowerAll(m.Tags)
	s.PublicationDate, err = time.Parse(r.DateFormat, m.DateStr)
	s.Extras = m.Extras

	if err != nil {
		return s, err
	}

	content, err := os.ReadFile(path.Join(storyDir, ContentFileName))
	if err != nil {
		return s, fmt.Errorf("failed to read content file: %w", err)
	}
	s.Content = parseContent(string(content))

	return s, nil
}

// parseMetadata parses the metadata block of the file.
func parseMetadata(metadataJSON string) (metadata, error) {
	var m metadata
	err := json.Unmarshal([]byte(metadataJSON), &m)
	if err != nil {
		return metadata{}, err
	}
	return m, nil
}

// parseContent parses a story in Markdown format and converts it to HTML.
func parseContent(content string) template.HTML {
	return template.HTML(bf.Run([]byte(content), bf.WithRenderer(renderer), bf.WithExtensions(markdownExtensions)))
}

func lowerAll(strs []string) []string {
	out := make([]string, len(strs))
	for i, v := range strs {
		out[i] = strings.ToLower(v)
	}
	return out
}

// storiesSlice is a wrapper type for a slice of stories, which provides sorting
// capability (by publication date).
type storiesSlice []Story

func (s storiesSlice) Len() int           { return len(s) }
func (s storiesSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s storiesSlice) Less(i, j int) bool { return s[i].PublicationDate.After(s[j].PublicationDate) }
