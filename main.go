package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"os"
	"regexp"
)

const (
	contentLoc    = "content/"
	metadatalFile = contentLoc + "metadata.json"
	storiesLoc    = contentLoc + "stories/"
	templFileLoc  = contentLoc + "templates/"
	staticLoc     = contentLoc + "static/"

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

var (
	templates = make(map[string]*template.Template)
	stories   = make(map[StoryURLPath]Story)
)

type Page struct {
	Title   string
	Content template.HTML
	Data    interface{}
}

type StoryURLPath string

type Metadata struct {
	Stories []StoryMetadata `json:"stories"`
}

type StoryMetadata struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Date  string `json:"date"`
}

type Story struct {
	Title        string
	CreationTime time.Time
	LastEditTime time.Time
	Content      template.HTML
}

func main() {
	log.Println("Rendering templates...")
	templates["content"] = template.Must(template.ParseFiles(
		templFileLoc+"content.html",
		templFileLoc+"base.html",
	))
	templates["list"] = template.Must(template.ParseFiles(
		templFileLoc+"list.html",
		templFileLoc+"base.html",
	))

	log.Println("Parsing metadata...")
	metadataJSON, err := ioutil.ReadFile(metadatalFile)
	check(err)
	var metadata Metadata
	err = json.Unmarshal(metadataJSON, &metadata)
	check(err)

	for _, story := range metadata.Stories {
		storyPath := storiesLoc + story.Name + "." + markdownFileFormat
		if _, err := os.Stat(storyPath); os.IsNotExist(err) {
			log.Fatalf("Can't find story: %s", story.Name)
		}
		stories[StoryURLPath(story.Name)] = Story{
			Title:   story.Title,
			Content: readStory(storyPath),
		}
	}

	fmt.Println("Starting server on localhost:8080...")
	err = http.ListenAndServe(":8080", makeRouter())
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// readStory parses a story in markdown format and converts it to HTML.
func readStory(filePath string) template.HTML {
	data, err := ioutil.ReadFile(filePath)
	check(err)

	renderer := blackfriday.HtmlRenderer(commonHtmlFlags, "", "")
	unsafe := blackfriday.Markdown(data, renderer, markdownExtensions)
	policy := bluemonday.UGCPolicy()
	policy.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
	return template.HTML(policy.SanitizeBytes(unsafe))
}

func makeRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/{name}", storyHandler)
	const staticPathPrefix = "/static"
	r.PathPrefix(staticPathPrefix).Handler(http.StripPrefix(staticPathPrefix, http.FileServer(http.Dir(staticLoc))))
	return r
}

func renderTemplate(name string, wr io.Writer, data interface{}) error {
	return templates[name].ExecuteTemplate(wr, "base", data)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	type ListItem struct {
		Path  StoryURLPath
		Story Story
	}
	var items []ListItem
	for path, story := range stories {
		items = append(items, ListItem{Path: path, Story: story})
	}
	err := renderTemplate("list", w, Page{
		Title: "stdout",
		Data:  items,
	})
	check(err)
}

func storyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if story, ok := stories[StoryURLPath(vars["name"])]; ok {
		err := renderTemplate("content", w, Page{
			Title:   story.Title,
			Content: story.Content,
		})
		check(err)
	} else {
		http.NotFound(w, r)
	}
}
