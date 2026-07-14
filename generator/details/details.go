package details

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path"

	"github.com/otiai10/copy"
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
	// The page is served at /<name> (from <name>.html), and its assets live in a
	// sibling /<name>/ directory, referenced from the content with absolute
	// paths such as /<name>/image.png.
	writer.WriteFile(path.Join(outputDir, s.Name+".html"), templateOutput.Bytes())

	copyAssets(s, path.Join(outputDir, s.Name))
}

// copyAssets copies the story's static assets (everything in its bundle other
// than the content and metadata files) into the story's asset directory. It
// does nothing if the story has no assets, so no empty directory is created.
func copyAssets(s reader.Story, assetDir string) {
	if s.SourceDir == "" || !hasAssets(s.SourceDir) {
		return
	}
	err := copy.Copy(s.SourceDir, assetDir, copy.Options{
		Skip: func(_ os.FileInfo, src, _ string) (bool, error) {
			return isBundleFile(path.Base(src)), nil
		},
	})
	if err != nil {
		log.Fatalf("Failed to copy assets for %s: %v", s.Name, err)
	}
}

// hasAssets reports whether the bundle directory contains any files other than
// the content and metadata files.
func hasAssets(sourceDir string) bool {
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		log.Fatalf("Failed to read bundle directory %s: %v", sourceDir, err)
	}
	for _, e := range entries {
		if !isBundleFile(e.Name()) {
			return true
		}
	}
	return false
}

func isBundleFile(name string) bool {
	return name == reader.ContentFileName || name == reader.MetadataFileName
}
