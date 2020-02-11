package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go.roman.zone/publisher/generator"
	"go.roman.zone/publisher/generator/rss"
	"go.roman.zone/publisher/reader"
	"html/template"
	"log"
	"os"
	"path/filepath"
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

	fmt.Println("Generating the website...")
	defer fmt.Println("Done!")

	var err error
	config, err = readConfiguration(configLoc)
	if err != nil {
		log.Fatalf("Failed to read configuration file: %s", err)
	}

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

func processStories(dir string, skipDrafts bool) {
	fmt.Println("> Processing stories...")

	r := reader.NewReader()
	r.SkipDrafts = skipDrafts

	stories, err := r.ReadAll(dir)
	check(err)

	generatorConfig := generator.WebsiteGeneratorConfig{
		IndexTemplate:       prepareTemplate("index.html"),
		DetailsTemplate:     prepareTemplate("details.html"),
		StaticFilesLocation: staticLoc,
		RSSFeedConfig:       config.Feed,
	}
	generatorConfig.GenerateWebsite(stories, *outputDir)
}

func prepareTemplate(fileName string) *template.Template {
	return template.Must(template.ParseFiles(
		filepath.Join(templateLoc, "base.html"),
		filepath.Join(templateLoc, fileName),
	))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
