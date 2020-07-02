package robots

import (
	"fmt"
	"path"

	"go.roman.zone/publisher/writer"
)

func GenerateRobotsTxtFile(outputDir string) {
	fmt.Println("> Generating the robots.txt file...")
	defer fmt.Println("  Finished generating the robots.txt file!")

	content := `User-agent: *
Disallow: /static/`

	writer.WriteFile(path.Join(outputDir, "robots.txt"), []byte(content))
}
