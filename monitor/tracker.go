package monitor

import (
	"fmt"
	"github.com/mmcdole/gofeed"
	"time"
)

type Tracker struct {
}

type Sighting struct {
	Map       string
	Selling   string
	Timestamp time.Time
}

func (t *Tracker) Find() (sightings []Sighting, err error) {
	var items []*gofeed.Item
	items, err = getFeed()
	if err != nil {
		return
	}

	for _, item := range items {
		var timestamp time.Time
		if item.PublishedParsed != nil {
			timestamp = *item.PublishedParsed
		}
		sightings = append(sightings, Sighting{
			Map:       item.Title,
			Selling:   item.Description,
			Timestamp: timestamp,
		})
	}

	return
}

func getFeed() (entries []*gofeed.Item, err error) {
	fp := gofeed.NewParser()
	var feed *gofeed.Feed
	feed, err = fp.ParseURL("https://whereismaurice.com/rss")
	if err != nil {
		return nil, fmt.Errorf("bad feed: %w", err)
	}

	return feed.Items, nil
}
