package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"

	"github.com/otiai10/copy"
	"go.roman.zone/publisher/generator/details"
	"go.roman.zone/publisher/generator/index"
	"go.roman.zone/publisher/generator/robots"
	"go.roman.zone/publisher/generator/rss"
	"go.roman.zone/publisher/reader"
)

type WebsiteGeneratorConfig struct {
	// IsDraft value indicates whether the generator is running in the draft
	// environment. If it is, draft stories won't be included in the generated
	// website.
	IsDraft bool

	StoriesDir    string
	StaticDir     string
	RobotsTxtPath string

	// OutputDir indicates output directory for the website files. At the end of
	// the execution output will be set up in a way that can be used on pretty
	// much any static website hosting.
	OutputDir string

	IndexTemplate   *template.Template
	DetailsTemplate *template.Template

	RSSFeedConfiguration rss.FeedConfiguration
}

func (c *WebsiteGeneratorConfig) GenerateWebsite() {

	if _, err := os.Stat(c.StoriesDir); os.IsNotExist(err) {
		log.Fatalf("Stories directory (%s) is missing.", c.StoriesDir)
	}

	// Creating the output directory before writing anything there
	err := os.MkdirAll(c.OutputDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("> Processing stories...")
	if c.IsDraft {
		fmt.Println("> ⚠️ DRAFT MODE")
	}

	r := reader.NewReader()
	r.SkipDrafts = !c.IsDraft
	stories, err := r.ReadAll(c.StoriesDir)
	if err != nil {
		log.Fatal(err)
	}

	generateAll(stories, c)
	copyStaticFiles(c)
}

func generateAll(stories *[]reader.Story, c *WebsiteGeneratorConfig) {
	index.GenerateIndexPage(stories, c.IndexTemplate, c.OutputDir)
	details.GenerateDetailsPages(stories, c.DetailsTemplate, c.OutputDir)
	rss.GenerateRSS(c.RSSFeedConfiguration, stories, c.OutputDir)
	robots.GenerateRobotsTxtFile(c.RobotsTxtPath, c.OutputDir)
}

func copyStaticFiles(c *WebsiteGeneratorConfig) {
	err := copy.Copy(c.StaticDir, path.Join(c.OutputDir, "static"))
	if err != nil {
		log.Fatal(err)
	}
}
