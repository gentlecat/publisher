package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gorilla/mux"
	"go.roman.zone/publisher/story"
)

var (
	listenHost = flag.String("host", "127.0.0.1", "Host to listen on")
	listenPort = flag.Int("port", 8080, "Port to listen on")

	storiesMutex sync.Mutex
	stories      []story.Story
	storiesIndex map[string]*story.Story

	templatesMutex sync.Mutex
	templates      map[string]*template.Template

	contentLoc    = "content"
	metadatalFile = filepath.Join(contentLoc, "metadata.json")
	storiesLoc    = filepath.Join(contentLoc, "stories")
	templLoc      = filepath.Join(contentLoc, "templates")
	staticLoc     = filepath.Join(contentLoc, "static")
)

type Page struct {
	Name    string
	Title   string
	Date    time.Time
	Content template.HTML
	Tags    []string
	Data    interface{}
}

func main() {
	flag.Parse()
	if customContentLoc := os.Getenv("CONTENT_LOCATION"); customContentLoc != "" {
		updateContentLoc(customContentLoc)
	}

	log.Println("Initializing...")
	renderTemplates(templLoc)
	readStories(metadatalFile, storiesLoc)

	watcher, err := fsnotify.NewWatcher()
	check(err)
	defer watcher.Close()
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				handleEvent(event)
			case err := <-watcher.Errors:
				log.Println("Watcher error:", err)
			}
		}
	}()
	err = watcher.Add(metadatalFile)
	check(err)
	err = watcher.Add(storiesLoc)
	check(err)
	err = watcher.Add(templLoc)
	check(err)

	listenAddr := fmt.Sprintf("%s:%d", *listenHost, *listenPort)
	log.Printf("Starting server on http://%s...\n", listenAddr)
	err = http.ListenAndServe(listenAddr, makeRouter())
	check(err)
}

func renderTemplates(location string) {
	log.Println("Rendering templates...")
	defer log.Println("Done!")
	templatesMutex.Lock()
	defer templatesMutex.Unlock()
	templates = make(map[string]*template.Template)
	templates["content"] = template.Must(template.ParseFiles(
		filepath.Join(location, "content.html"),
		filepath.Join(location, "base.html"),
	))
	templates["list"] = template.Must(template.ParseFiles(
		filepath.Join(location, "list.html"),
		filepath.Join(location, "base.html"),
	))
}

func readStories(metadatalFile, storiesLoc string) {
	log.Println("Parsing stories...")
	defer log.Println("Done!")
	storiesMutex.Lock()
	defer storiesMutex.Unlock()
	stories = story.ReadStories(metadatalFile, storiesLoc)
	storiesIndex = make(map[string]*story.Story)
	for i, s := range stories {
		storiesIndex[s.Name] = &stories[i]
	}
}

func makeRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/robots.txt", robotsHandler)
	r.HandleFunc("/{name}", storyHandler)
	const staticPathPrefix = "/static"
	r.PathPrefix(staticPathPrefix).Handler(http.StripPrefix(staticPathPrefix, http.FileServer(http.Dir(staticLoc))))
	return r
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	type ListItem struct {
		Path  string
		Story *story.Story
	}
	var items []ListItem
	storiesMutex.Lock()
	for i, s := range stories {
		items = append(items, ListItem{Path: s.Name, Story: &stories[i]})
	}
	storiesMutex.Unlock()
	err := renderTemplate("list", w, Page{
		Data: items,
	})
	check(err)
}

func storyHandler(w http.ResponseWriter, r *http.Request) {
	storiesMutex.Lock()
	defer storiesMutex.Unlock()
	if s, ok := storiesIndex[mux.Vars(r)["name"]]; ok {
		err := renderTemplate("content", w, Page{
			Name:    s.Name,
			Title:   s.Title,
			Date:    s.Date,
			Content: s.Content,
			Tags:    s.Tags,
		})
		check(err)
	} else {
		http.NotFound(w, r)
	}
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `User-agent: *
Disallow: /static/`)
}

func renderTemplate(name string, wr io.Writer, data interface{}) error {
	templatesMutex.Lock()
	defer templatesMutex.Unlock()
	return templates[name].ExecuteTemplate(wr, "base", data)
}

func handleEvent(event fsnotify.Event) {
	if (event.Op&fsnotify.Create == fsnotify.Create) ||
		(event.Op&fsnotify.Remove == fsnotify.Remove) ||
		(event.Op&fsnotify.Write == fsnotify.Write) {
		log.Printf("Detected change in file: %s", event.Name)
		if strings.HasPrefix(event.Name, templLoc) {
			renderTemplates(templLoc)
		}
		if strings.HasPrefix(event.Name, metadatalFile) || strings.HasPrefix(event.Name, storiesLoc) {
			readStories(metadatalFile, storiesLoc)
		}
	}
}

func updateContentLoc(directoryPath string) {
	contentLoc = filepath.Clean(directoryPath)
	metadatalFile = filepath.Join(contentLoc, "metadata.json")
	storiesLoc = filepath.Join(contentLoc, "stories")
	templLoc = filepath.Join(contentLoc, "templates")
	staticLoc = filepath.Join(contentLoc, "static")
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
