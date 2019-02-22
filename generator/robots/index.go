package robots

import (
	"go.roman.zone/publisher/generator"
	"log"
	"path"
)

func GenerateRobotsTxtFile(outputDir string) {
	log.Println("Generating robots.txt file...")
	defer log.Println("Finished generating robots.txt file!")

	content := `User-agent: *
Disallow: /static/`

	generator.CheckedFileWriter(path.Join(outputDir, "robots.txt"), []byte(content))

}
