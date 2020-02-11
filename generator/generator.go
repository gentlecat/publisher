package generator

import (
	"github.com/otiai10/copy"
	"go.roman.zone/publisher/generator/details"
	"go.roman.zone/publisher/generator/index"
	"go.roman.zone/publisher/generator/robots"
	"go.roman.zone/publisher/generator/rss"
	"go.roman.zone/publisher/reader"
	"html/template"
	"log"
	"os"
	"path"
)

type WebsiteGeneratorConfig struct {
	IndexTemplate   *template.Template
	DetailsTemplate *template.Template

	StaticFilesLocation string

	RSSFeedConfig rss.FeedConfiguration
}

func (c *WebsiteGeneratorConfig) GenerateWebsite(stories *[]reader.Story, outputDir string) {
	// Creating the output directory before writing anything there
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	index.GenerateIndexPage(stories, c.IndexTemplate, outputDir)
	details.GenerateDetailsPages(stories, c.DetailsTemplate, outputDir)

	rss.GenerateRSS(c.RSSFeedConfig, stories, outputDir)

	robots.GenerateRobotsTxtFile(outputDir)

	err = copy.Copy(c.StaticFilesLocation, path.Join(outputDir, "static"))
	if err != nil {
		log.Fatal(err)
	}
}
