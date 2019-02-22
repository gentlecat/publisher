package generator

import (
	"html/template"
	"log"
	"os"
	"time"
)

type Output struct {
	Content  []byte
	FileType string
}

type ContentGenerator interface {
	generate() []Output
}

// PageContext contains actual content that gets sent to a template.
type PageContext struct {
	Title string
	Data  interface{} // Additional data that doesn't fit into any other field.
	// TODO: Perhaps look into moving fields defined below into `Data`:
	Name    string
	Content template.HTML
	Tags    []string
	Date    time.Time
}

func CheckedFileWriter(path string, content []byte) {
	f, err := os.Create(path)
	check(err, path)
	defer f.Close()

	_, err = f.Write(content)
	check(err, path)
}
func check(e error, filePath string) {
	if e != nil {
		log.Fatalf("Failed to write file %s", filePath)
	}
}
