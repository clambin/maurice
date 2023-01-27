package monitor

import (
	"context"
	"fmt"
	"github.com/mmcdole/gofeed"
	"time"
)

const MauriceURL = "https://whereismaurice.com/rss"

type Tracker struct {
	URL string
}

type Sighting struct {
	Map       string
	Selling   string
	Timestamp time.Time
}

func (t *Tracker) Find(ctx context.Context) ([]Sighting, error) {
	var sightings []Sighting
	items, err := getFeed(ctx, t.URL)
	if err == nil {
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
	}
	return sightings, err
}

func getFeed(ctx context.Context, feedURL string) ([]*gofeed.Item, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURLWithContext(feedURL, ctx)
	if err != nil {
		return nil, fmt.Errorf("bad feed: %w", err)
	}
	return feed.Items, nil
}
