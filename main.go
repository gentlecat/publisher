package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/otiai10/copy"
	"go.roman.zone/publisher/generator/details"
	"go.roman.zone/publisher/generator/index"
	"go.roman.zone/publisher/generator/robots"
	"go.roman.zone/publisher/generator/rss"
	"go.roman.zone/publisher/story"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

var (
	prodEnv    = flag.Bool("prod", false, "Whether the generator is running in production environment (if it is, draft stories won't be included)")
	outputDir  = flag.String("out", "./out", "Output directory for the content")
	contentLoc = flag.String("content", "./content", "Path to the content directory")

	storiesLoc  string
	templateLoc string
	staticLoc   string
	configLoc   string

	config Configuration

	// TODO: See if template stuff and serving needs to be in separate modules
	templates      map[string]*template.Template
	templatesMutex sync.Mutex
)

type Configuration struct {
	Feed rss.FeedConfiguration
}

func main() {
	flag.Parse()

	if _, err := os.Stat(*contentLoc); os.IsNotExist(err) {
		log.Fatalf("Content directory (%s) is missing.", *contentLoc)
	}

	storiesLoc = filepath.Join(*contentLoc, "stories")
	templateLoc = filepath.Join(*contentLoc, "templates")
	staticLoc = filepath.Join(*contentLoc, "static")
	configLoc = filepath.Join(*contentLoc, "config.json")

	// Measuring the time it takes to execute
	start := time.Now()
	defer fmt.Println(time.Since(start))

	log.Println("Initializing...")
	var err error
	config, err = readConfiguration(configLoc)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %s", err)
	}

	renderTemplates(templateLoc)
	processStories(storiesLoc, *prodEnv) // drafts are ignored in production
}

func readConfiguration(location string) (config Configuration, err error) {
	file, err := os.Open(location)
	if err != nil {
		return
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return
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

func processStories(dir string, ignoreDrafts bool) {
	log.Println("Processing stories...")
	defer log.Println("Done!")

	stories, err := story.ReadAll(dir, ignoreDrafts)
	check(err)

	// Creating the output directory before writing anything there
	check(os.MkdirAll(*outputDir, os.ModePerm))

	index.GenerateIndexPage(stories, templates["list"], *outputDir)
	details.GenerateDetailsPages(stories, templates["content"], *outputDir)
	rss.GenerateRSS(config.Feed, stories, *outputDir)
	robots.GenerateRobotsTxtFile(*outputDir)
	check(copy.Copy(staticLoc, path.Join(*outputDir, "static")))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
