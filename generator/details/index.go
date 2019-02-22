package details

import (
	"bytes"
	"fmt"
	"go.roman.zone/publisher/generator"
	"go.roman.zone/publisher/story"
	"html/template"
	"log"
	"path"
)

func GenerateDetailsPages(stories *[]story.Story, tpl *template.Template, outputDir string) {
	log.Println("Generating details pages...")
	defer log.Println("Finished generating all details pages!")

	for _, s := range *stories {
		generateStoryFile(s, tpl, outputDir)
	}
}

func generateStoryFile(s story.Story, tpl *template.Template, outputDir string) {
	log.Printf("Generating details page for %s", s.Name)

	var templateOutput bytes.Buffer

	if err := tpl.ExecuteTemplate(&templateOutput, "base", generator.PageContext{
		Name:    s.Name,
		Title:   s.Title,
		Date:    s.PublicationDate,
		Content: s.Content,
		Tags:    s.Tags,
	}); err != nil {
		log.Fatalf("Failed to render details page for %s: %v", s.Name, err)
	}

	generator.CheckedFileWriter(path.Join(outputDir, fmt.Sprintf("%s.html", s.Name)), templateOutput.Bytes())
}
