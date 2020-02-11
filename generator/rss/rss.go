package rss

import (
	"fmt"
	"github.com/gorilla/feeds"
	"go.roman.zone/publisher/reader"
	"go.roman.zone/publisher/writer"
	"log"
	"path"
	"time"
)

type FeedConfiguration struct {
	Title       string
	Description string
	Author      feeds.Author
	Backlink    feeds.Link
}

// GenerateRSS generates an RSS feed file.
func GenerateRSS(config FeedConfiguration, stories *[]reader.Story, outputDir string) {
	fmt.Println("> Generating the RSS feed...")
	defer fmt.Println("  Finished generating the RSS feed!")

	feed := &feeds.Feed{
		Title:       config.Title,
		Link:        &config.Backlink,
		Description: config.Description,
		Author:      &config.Author,
		Updated:     time.Now(),
	}

	for _, s := range *stories {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       s.Title,
			Link:        &feeds.Link{Href: config.Backlink.Href + "/" + s.Name},
			Description: string(s.Content),
			Created:     s.PublicationDate,
		})
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.Fatalf("Failed to generate the RSS feed: %v", err)
	}

	writer.WriteFile(path.Join(outputDir, fmt.Sprintf("rss.xml")), []byte(rss))
}
