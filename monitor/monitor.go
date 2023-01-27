package monitor

import (
	"context"
	"github.com/containrrr/shoutrrr"
	"golang.org/x/exp/slog"
	"time"
)

type Server struct {
	Tracker
	Notifier
	lastUpdate time.Time
}

type Notifier interface {
	Send(title, message string) []error
}

func New(notifierURL string) *Server {
	s := Server{Tracker: Tracker{URL: MauriceURL}}

	if notifierURL != "" {
		router, err := shoutrrr.CreateSender(notifierURL)
		if err != nil {
			slog.Error("unable to create notification sender. Ignoring.", err)
		}
		s.Notifier = &ShoutrrrNotifier{router: router}
	}

	return &s
}

func (s *Server) Run(ctx context.Context) {
	s.check(ctx)
	interval := time.NewTicker(5 * time.Minute)
	defer interval.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-interval.C:
			s.check(ctx)
		}
	}
}

func (s *Server) check(ctx context.Context) {
	items, err := s.Tracker.Find(ctx)
	if err != nil {
		slog.Error("failed to check feed", err)
		return
	}
	for _, item := range items {
		if item.Timestamp.After(s.lastUpdate) {
			_ = s.Send("Maurice is now at "+item.Map, "Selling: "+item.Selling+"\nhttps://whereismaurice.com")
			slog.Info("Maurice found at new location", "location", item.Map)
			s.lastUpdate = item.Timestamp
		}
	}
}
