package cmd

import (
	"bytes"
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
		intervalSeconds   int
		lookBehindMinutes int
	)

	// TODO: we could take a state file as param.
	// TODO: If param is passed, we should try to
	// TODO: load state upon starting and write state
	// TODO: upon closing.
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
				LookBehindMinutes:     lookBehindMinutes,
				TickerIntervalSeconds: intervalSeconds,
			}

			buf := bytes.NewBuffer(nil)
			state := office365.NewGOBState()

			defer func() {
				err := state.Write(buf)
				if err != nil {
					WriteOut("could not encode state to buffer")
				}
			}()

			err := state.Read(buf)
			if err != nil {
				WriteOut("could not decode state from buffer")
			}

			resultChan, err := client.Subscription.Watch(ctx, watcherConf, state)
			if err != nil {
				logger.Printf("error occured calling watch: %s\n", err)
				return
			}

			printer := office365.NewPrinter(defaultOutput)
			printer.Handle(resultChan)
		},
	}
	cmd.Flags().IntVar(&intervalSeconds, "interval", 5, "TickerIntervalSeconds")
	cmd.Flags().IntVar(&lookBehindMinutes, "lookbehind", 1, "Number of minutes from request time used when fetching available content.")

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
