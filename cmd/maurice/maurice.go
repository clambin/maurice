package main

import (
	"context"
	"github.com/clambin/maurice/monitor"
	"github.com/clambin/maurice/version"
	"github.com/xonvanetta/shutdown/pkg/shutdown"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
	"path/filepath"
)

func main() {
	var (
		debug       bool
		notifierURL string
	)

	a := kingpin.New(filepath.Base(os.Args[0]), "maurice")
	a.Version(version.BuildVersion)
	a.HelpFlag.Short('h')
	a.VersionFlag.Short('v')
	a.Flag("debug", "Log debug messages").BoolVar(&debug)
	a.Flag("notifier", "Notification URL (shoutrrr format)").StringVar(&notifierURL)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		panic(err)
	}

	m := monitor.New(notifierURL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go m.Run(ctx)

	<-shutdown.Chan()
}
