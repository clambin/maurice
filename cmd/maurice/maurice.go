package main

import (
	"context"
	"github.com/clambin/maurice/monitor"
	"github.com/clambin/maurice/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
	"os"
	"os/signal"
	"syscall"
)

var (
	cmd = &cobra.Command{
		Use:     "maurice",
		Short:   "monitors whereismaurice.com and notifies changes via Slack",
		Version: version.BuildVersion,
		Run:     Main,
	}
)

func main() {
	if err := cmd.Execute(); err != nil {
		slog.Error("failed to start", err)
		os.Exit(1)
	}
}

func Main(_ *cobra.Command, _ []string) {
	m := monitor.New(viper.GetString("notifier"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go m.Run(ctx)

	ctx, done := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer done()
	<-ctx.Done()
}

func init() {
	cobra.OnInitialize(initConfig)
	cmd.Flags().StringP("notifier", "n", "", "Notifications URL (in ShoutRrr format)")
	_ = cmd.MarkFlagRequired("notifier")
	_ = viper.BindPFlag("notifier", cmd.Flags().Lookup("notifier"))
}

func initConfig() {
	viper.AddConfigPath("/etc/maurice/")
	viper.AddConfigPath("$HOME/.maurice")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	viper.SetEnvPrefix("MAURICE")
	viper.AutomaticEnv()

	_ = viper.ReadInConfig()
}
