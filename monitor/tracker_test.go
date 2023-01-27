package monitor_test

import (
	"context"
	"github.com/clambin/maurice/monitor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTracker_Find(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(testContent))
	}))
	defer s.Close()
	tracker := monitor.Tracker{URL: s.URL}

	items, err := tracker.Find(context.Background())
	require.NoError(t, err)
	assert.Equal(t, []monitor.Sighting{
		{
			Map:       "Jakobs Estate, Eden-6",
			Selling:   "[Ruby's Wrath](https://www.lootlemon.com/weapon/rubys-wrath-bl3), [The Horizon](https://www.lootlemon.com/weapon/the-horizon-bl3) and [The Tidal Wave](https://www.lootlemon.com/weapon/the-tidal-wave-bl3)",
			Timestamp: time.Date(2023, time.January, 19, 17, 0, 0, 0, time.UTC),
		},
	}, items)
}

const testContent = `
<?xml version="1.0" encoding="UTF-8"?>
<rss xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
  <channel>
    <title><![CDATA[Where is Maurice?]]></title>
    <description><![CDATA[Where in the heck is Maurice's Black Market vending machine?]]></description>
    <link>https://whereismaurice.com</link>
    <generator>GatsbyJS</generator>
    <lastBuildDate>Thu, 19 Jan 2023 22:43:22 GMT</lastBuildDate>
    <item>
      <title><![CDATA[Jakobs Estate, Eden-6]]></title>
      <description><![CDATA[[Ruby's Wrath](https://www.lootlemon.com/weapon/rubys-wrath-bl3), [The Horizon](https://www.lootlemon.com/weapon/the-horizon-bl3) and [The Tidal Wave](https://www.lootlemon.com/weapon/the-tidal-wave-bl3)]]></description>
      <link>https://whereismaurice.com/#map</link>
      <guid isPermaLink="false">WIM-GUID-2023-01-19T12:00-05:00</guid>
      <pubDate>Thu, 19 Jan 2023 17:00:00 GMT</pubDate>
    </item>
  </channel>
</rss>%`
