package index

import (
	"bytes"
	"fmt"
	"go.roman.zone/publisher/generator"
	"go.roman.zone/publisher/story"
	"html/template"
	"log"
	"path"
)

func GenerateIndexPage(stories *[]story.Story, tpl *template.Template, outputDir string) {
	log.Println("Generating index page...")
	defer log.Println("Finished generating the index page!")

	var templateOutput bytes.Buffer

	type ListItem struct {
		Path  string
		Story *story.Story
	}
	var items []ListItem

	for i, s := range *stories {
		items = append(items, ListItem{Path: s.Name, Story: &(*stories)[i]})
	}
	type PageData struct {
		Stories []ListItem
	}

	if err := tpl.ExecuteTemplate(&templateOutput, "base", generator.PageContext{
		Data: PageData{
			Stories: items,
		},
	}); err != nil {
		log.Fatalf("Failed to render index page: %v", err)
	}

	generator.CheckedFileWriter(path.Join(outputDir, fmt.Sprintf("index.html")), templateOutput.Bytes())
}
