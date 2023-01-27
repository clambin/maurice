package monitor_test

import (
	"context"
	"github.com/clambin/maurice/monitor"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestServer_Run(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(testContent))
	}))
	defer s.Close()

	m := monitor.New("")
	m.Tracker.URL = s.URL
	ch := make(chan notification)
	m.Notifier = &fakeNotifier{ch: ch}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() { defer wg.Done(); m.Run(ctx) }()

	n := <-ch
	assert.Equal(t, notification{
		title:   "Maurice is now at Jakobs Estate, Eden-6",
		message: "Selling: [Ruby's Wrath](https://www.lootlemon.com/weapon/rubys-wrath-bl3), [The Horizon](https://www.lootlemon.com/weapon/the-horizon-bl3) and [The Tidal Wave](https://www.lootlemon.com/weapon/the-tidal-wave-bl3)\nhttps://whereismaurice.com",
	}, n)

	cancel()
	wg.Wait()
}

var _ monitor.Notifier = &fakeNotifier{}

type fakeNotifier struct {
	ch chan notification
}

type notification struct {
	title   string
	message string
}

func (f fakeNotifier) Send(title, message string) []error {
	f.ch <- notification{title: title, message: message}
	return []error{}
}
