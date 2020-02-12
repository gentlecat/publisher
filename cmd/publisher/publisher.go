// Publisher is a lightweight (by my standards) static website generator. Most
// of it is made for building my own websites, however all of it is meant to be
// reusable for different kinds of websites.
//
// Main purpose is to support blog-like websites, but you can obviously rewrite
// parts which don't work for your particular project.
//
// This binary is a default wrapper around the website generator. Implementation
// is simple enough that you should be able to write a custom wrapper, if you
// need to.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go.roman.zone/publisher/generator"
	"go.roman.zone/publisher/generator/rss"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	prodEnv    = flag.Bool("prod", false, "Whether the generator is running in production environment (if it is, draft stories won't be included)")
	outputDir  = flag.String("out", "./out", "Output directory for the content")
	contentDir = flag.String("content", "./content", "Path to the content directory")
)

// Configuration struct represents the expected structure of configuration file.
//
// Configuration file itself is expected to be in the content directory and be
// named `config.json`.
type Configuration struct {
	Feed rss.FeedConfiguration
}

func main() {
	flag.Parse()

	start := time.Now()
	fmt.Println("Generating the website...")
	defer fmt.Printf("Done in %v!", time.Since(start))

	config, err := readConfiguration(filepath.Join(*contentDir, "config.json"))
	if err != nil {
		log.Fatalf("Failed to read configuration file: %s", err)
	}

	templateDir := filepath.Join(*contentDir, "templates")

	generatorConfig := generator.WebsiteGeneratorConfig{
		IsProd:    *prodEnv,
		OutputDir: *outputDir,

		StoriesDir: filepath.Join(*contentDir, "stories"),
		StaticDir:  filepath.Join(*contentDir, "static"),

		IndexTemplate:        prepareTemplate(templateDir, "index.html"),
		DetailsTemplate:      prepareTemplate(templateDir, "details.html"),
		RSSFeedConfiguration: config.Feed,
	}
	generatorConfig.GenerateWebsite()
}

func readConfiguration(location string) (config Configuration, err error) {
	fmt.Println("> Reading configuration...")

	file, err := os.Open(location)
	if err != nil {
		return
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return
}

func prepareTemplate(templateDir string, fileName string) *template.Template {
	return template.Must(template.ParseFiles(
		filepath.Join(templateDir, "base.html"),
		filepath.Join(templateDir, fileName),
	))
}
