package sitemap

import (
	"fmt"
	"go.roman.zone/publisher/generator/rss"
	"go.roman.zone/publisher/reader"
	"go.roman.zone/publisher/writer"
	"path"
	"strings"
)

func GenerateSitemap(config rss.FeedConfiguration, stories *[]reader.Story, outputDir string) {
	fmt.Println("> Generating the sitemap...")
	defer fmt.Println("  Finished generating the sitemap!")

	var links []string

	links = append(links, config.Backlink.Href)

	for _, s := range *stories {
		links = append(links, config.Backlink.Href+"/"+s.Name)
	}

	writer.WriteFile(path.Join(outputDir, fmt.Sprintf("sitemap.txt")), []byte(strings.Join(links, "\n")+"\n"))
}
