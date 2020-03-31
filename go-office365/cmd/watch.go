package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/devodev/go-office365/office365"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newCommandWatch())
}

func newCommandWatch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch [config]",
		Short: "Fetch audit events at regular intervals.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			confArg := args[0]

			watchConfig, err := loadConfig(confArg)
			if err != nil {
				logger.Printf("error occured loading config file: %s\n", err)
				return
			}

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)

			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				sigChan := getSigChan()
				for {
					select {
					case <-sigChan:
						cancel()
						return
					default:
					}
				}
			}()

			watcherConf := office365.SubscriptionWatcherConfig{
				FetcherCount:           watchConfig.Global.FetcherCount,
				LookBehindMinutes:      watchConfig.Global.FetcherLookBehindMinutes,
				FetcherIntervalSeconds: watchConfig.Global.FetcherIntervalSeconds,
				TickerIntervalSeconds:  watchConfig.Global.TickerIntervalSeconds,
			}

			resultChan, err := client.Watch.Watch(ctx, watcherConf)
			if err != nil {
				logger.Printf("error occured calling watch: %s\n", err)
				return
			}

			printer := office365.NewPrinter(defaultOutput)
			printer.Handle(resultChan)
		},
	}
	return cmd
}

// WatchConfig .
type WatchConfig struct {
	Global struct {
		TickerIntervalSeconds    int
		FetcherCount             int
		FetcherIntervalSeconds   int
		FetcherLookBehindMinutes int
		PubIdentifier            string
	}
}

func loadConfig(confPath string) (*WatchConfig, error) {
	vip := viper.New()
	vip.SetConfigFile(confPath)

	if err := vip.ReadInConfig(); err != nil {
		return nil, err
	}

	var config WatchConfig
	if err := vip.UnmarshalExact(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func getSigChan() chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	return sigChan
}
