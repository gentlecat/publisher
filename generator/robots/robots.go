package robots

import (
	"fmt"
	"os"
	"path"

	"go.roman.zone/publisher/writer"
)

func GenerateRobotsTxtFile(robotsTxtPath, outputDir string) {
	fmt.Println("> Generating the robots.txt file...")

	existingRobotsTxt, err := os.ReadFile(robotsTxtPath)
	if err != nil {
		content := `User-agent: *
		Disallow: /static/`
		writer.WriteFile(path.Join(outputDir, "robots.txt"), []byte(content))
		fmt.Println("  Generated default robots.txt file!")
	} else {
		writer.WriteFile(path.Join(outputDir, "robots.txt"), existingRobotsTxt)
		fmt.Println("  Copied robots.txt file!")
	}
}
