package writer

import (
	"log"
	"os"
)

// WriteFile creates a specified path and writes provided bytes into it.
func WriteFile(path string, content []byte) {
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
