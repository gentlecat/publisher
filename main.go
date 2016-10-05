package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"regexp"
)

const (
	contentLoc   = "content/"
	templFileLoc = contentLoc + "templates/"
	staticLoc    = contentLoc + "static/"
	storiesLoc   = contentLoc + "stories/"

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
		blackfriday.HTML_TOC |
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

	log.Println("Reading stories...")
	files, err := ioutil.ReadDir(storiesLoc)
	check(err)

	format := strings.ToLower("." + markdownFileFormat)
	formatLen := len(format)
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), format) {
			if len(file.Name())-formatLen < 1 {
				log.Fatal("Strange file name: ", file.Name())
			}
			name := file.Name()[:len(file.Name())-formatLen]
			log.Printf(" - %s", name)
			// TODO: Figure out how to set up titles, timestamps.
			stories[StoryURLPath(name)] = Story{
				Content: readStory(file.Name()),
			}
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
func readStory(fileName string) template.HTML {
	data, err := ioutil.ReadFile(storiesLoc + fileName)
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
