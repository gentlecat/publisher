package details

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"path"

	"go.roman.zone/publisher/reader"
	"go.roman.zone/publisher/writer"
)

// GenerateDetailsPages generates a page for a specific story.
func GenerateDetailsPages(stories *[]reader.Story, tpl *template.Template, outputDir string) {
	fmt.Println("> Generating details pages...")
	defer fmt.Println("  Generated all details pages!")

	for _, s := range *stories {
		generateStoryFile(s, tpl, outputDir)
	}
}

func generateStoryFile(s reader.Story, tpl *template.Template, outputDir string) {
	fmt.Printf("  - %s\n", s.Name)

	var templateOutput bytes.Buffer

	if err := tpl.ExecuteTemplate(&templateOutput, "base", s); err != nil {
		log.Fatalf("Failed to render details page for %s: %v", s.Name, err)
	}

	writer.WriteFile(path.Join(outputDir, fmt.Sprintf("%s.html", s.Name)), templateOutput.Bytes())
}
