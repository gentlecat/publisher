package robots

import (
	"go.roman.zone/publisher/writer"
	"log"
	"path"
)

func GenerateRobotsTxtFile(outputDir string) {
	log.Println("Generating robots.txt file...")
	defer log.Println("Finished generating robots.txt file!")

	content := `User-agent: *
Disallow: /static/`

	writer.WriteFile(path.Join(outputDir, "robots.txt"), []byte(content))
}
