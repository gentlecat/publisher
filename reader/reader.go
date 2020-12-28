package reader

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	bf "github.com/russross/blackfriday"
)

type Configuration struct {
	// DateFormat is the format used to parse creation date of a story.
	DateFormat string

	// StoryFileFormat specifies expected file format for the files which are
	// going to be parsed. No dot necessary. For example, "md".
	StoryFileFormat string

	// MetadataSeparator specifies the string that's used to separate metadata
	// part of the story file and the rest of the content. Metadata is expected
	// to be first, followed by the specified separator, followed by content.
	MetadataSeparator string

	// SkipDrafts indicates whether draft stories should be ignored. This flag
	// will be set in the metadata of each story.
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
		DateFormat:        "2006-Jan-02",
		StoryFileFormat:   "md",
		MetadataSeparator: "\n+++\n",
	}
}

// ReadAll reads all stories in a specified directory and returns a list of
// them. The list is sorted by publication date.
func (r *Configuration) ReadAll(storiesDir string) (*[]Story, error) {
	stories := make([]Story, 0)

	files, err := ioutil.ReadDir(storiesDir)
	if err != nil {
		return &stories, err
	}
	for _, f := range files {
		if !r.isStoryFile(f) {
			continue
		}
		fmt.Printf("  - %s\n", f.Name())
		s, err := r.read(path.Join(storiesDir, f.Name()))
		if err != nil {
			log.Printf("Failed to read file: %s. Error: %s\n", f.Name(), err)
			continue
		}
		if r.SkipDrafts && s.IsDraft {
			continue
		}
		stories = append(stories, s)
	}
	sort.Sort(storiesSlice(stories))
	return &stories, nil
}

func (r *Configuration) isStoryFile(f os.FileInfo) bool {
	if f.IsDir() || !strings.HasSuffix(strings.ToLower(f.Name()), "."+r.StoryFileFormat) {
		return false
	}
	return true
}

func (r *Configuration) read(storyFilePath string) (s Story, err error) {
	data, err := ioutil.ReadFile(storyFilePath)
	if err != nil {
		return s, err
	}

	parts := strings.SplitN(string(data), r.MetadataSeparator, 2)
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
	s.PublicationDate, err = time.Parse(r.DateFormat, m.DateStr)
	s.Extras = m.Extras

	if err != nil {
		return s, err
	}
	if len(parts) == 2 {
		s.Content = parseContent(parts[1])
	}
	return s, err
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

// parseContent parses a story in markdown format and converts it to HTML.
func parseContent(content string) template.HTML {
	return template.HTML(
		bf.Run(
			[]byte(content),
			bf.WithExtensions(markdownExtensions),
			bf.WithRenderer(
				bf.NewHTMLRenderer(
					bf.HTMLRendererParameters{Flags: htmlRendererParams}))))
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

// storiesSlice is wrapper type for a slice of stories, which provides sorting
// capability (by publication date).
type storiesSlice []Story

func (s storiesSlice) Len() int           { return len(s) }
func (s storiesSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s storiesSlice) Less(i, j int) bool { return s[i].PublicationDate.After(s[j].PublicationDate) }
