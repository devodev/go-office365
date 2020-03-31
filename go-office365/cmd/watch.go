package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/devodev/go-office365/office365"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCommandWatch())
}

func newCommandWatch() *cobra.Command {
	var (
		tickerIntervalSeconds    int
		fetcherCount             int
		fetcherIntervalSeconds   int
		fetcherLookBehindMinutes int
	)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Fetch audit events at regular intervals.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
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
				FetcherCount:           fetcherCount,
				LookBehindMinutes:      fetcherLookBehindMinutes,
				FetcherIntervalSeconds: fetcherIntervalSeconds,
				TickerIntervalSeconds:  tickerIntervalSeconds,
			}

			resultChan, err := client.Subscription.Watch(ctx, watcherConf)
			if err != nil {
				logger.Printf("error occured calling watch: %s\n", err)
				return
			}

			printer := office365.NewPrinter(defaultOutput)
			printer.Handle(resultChan)
		},
	}
	cmd.Flags().IntVar(&tickerIntervalSeconds, "ticker-interval", 5, "TickerIntervalSeconds")
	cmd.Flags().IntVar(&fetcherCount, "fetcher-count", 10, "FetcherCount")
	cmd.Flags().IntVar(&fetcherIntervalSeconds, "fetcher-interval", 60, "FetcherIntervalSeconds")
	cmd.Flags().IntVar(&fetcherLookBehindMinutes, "fetcher-lookbehind", 1, "FetcherLookBehindMinutes")

	return cmd
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
