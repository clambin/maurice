package monitor

import (
	"context"
	"github.com/containrrr/shoutrrr"
	"github.com/containrrr/shoutrrr/pkg/router"
	log "github.com/sirupsen/logrus"
	"time"
)

type Server struct {
	Tracker
	Notifier
	lastUpdate time.Time
}

func New(notifierURL string) *Server {
	var r *router.ServiceRouter
	if notifierURL != "" {
		var err error
		r, err = shoutrrr.CreateSender(notifierURL)
		if err != nil {
			log.WithError(err).Error("unable to create notification sender. Ignoring.")
		}
	}

	return &Server{
		Notifier: Notifier{router: r},
	}
}

func (s *Server) Run(ctx context.Context) {
	s.check()

	interval := time.NewTicker(5 * time.Minute)
	for running := true; running; {
		select {
		case <-ctx.Done():
			running = false
		case <-interval.C:
			s.check()
		}
	}
}

func (s *Server) check() {
	items, err := s.Tracker.Find()
	if err != nil {
		log.WithError(err).Error("failed to check feed")
		return
	}
	for _, item := range items {
		if item.Timestamp.After(s.lastUpdate) {
			_ = s.Send("Maurice is now at "+item.Map, "Selling: "+item.Selling+"\nhttps://whereismaurice.com")
			log.Infof("Maurice is now at %s, selling: %s", item.Map, item.Selling)
			s.lastUpdate = item.Timestamp
		}
	}
}
