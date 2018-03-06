package feeds

import (
	"github.com/gorilla/feeds"
	"go.roman.zone/publisher/story"
	"time"
)

type FeedConfiguration struct {
	Title       string
	Description string
	Author      feeds.Author
	Backlink    feeds.Link
}

func GenerateRSS(stories []story.Story, config FeedConfiguration) (rss string, err error) {
	feed := &feeds.Feed{
		Title:       config.Title,
		Link:        &config.Backlink,
		Description: config.Description,
		Author:      &config.Author,
		Updated:     time.Now(),
	}

	for _, s := range stories {
		feed.Items = append(feed.Items, &feeds.Item{
			Title:       s.Title,
			Link:        &feeds.Link{Href: config.Backlink.Href + "/" + s.Name},
			Description: string(s.Content),
			Created:     s.Date,
		})
	}

	rss, err = feed.ToRss()
	return
}
