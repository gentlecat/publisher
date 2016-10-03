package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	markdownFileFormat = "md"
	contentLoc         = "content/"
	templFileLoc       = contentLoc + "templates/"
	staticLoc          = contentLoc + "static/"
	storiesLoc         = contentLoc + "stories/"
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
	unsafe := blackfriday.MarkdownCommon(data)
	return template.HTML(bluemonday.UGCPolicy().SanitizeBytes(unsafe))
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
