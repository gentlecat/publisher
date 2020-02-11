package details

import (
	"bytes"
	"fmt"
	"go.roman.zone/publisher/reader"
	"go.roman.zone/publisher/writer"
	"html/template"
	"log"
	"path"
)

func GenerateDetailsPages(stories *[]reader.Story, tpl *template.Template, outputDir string) {
	log.Println("Generating details pages...")
	defer log.Println("Finished generating all details pages!")

	for _, s := range *stories {
		generateStoryFile(s, tpl, outputDir)
	}
}

func generateStoryFile(s reader.Story, tpl *template.Template, outputDir string) {
	log.Printf("Generating details page for %s", s.Name)

	var templateOutput bytes.Buffer

	if err := tpl.ExecuteTemplate(&templateOutput, "base", s); err != nil {
		log.Fatalf("Failed to render details page for %s: %v", s.Name, err)
	}

	writer.WriteFile(path.Join(outputDir, fmt.Sprintf("%s.html", s.Name)), templateOutput.Bytes())
}
