package writer

import (
	"log"
	"os"
	"path/filepath"
)

// WriteFile creates a specified path and writes provided bytes into it,
// creating any parent directories that don't exist yet.
func WriteFile(path string, content []byte) {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		check(err, path)
	}

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
