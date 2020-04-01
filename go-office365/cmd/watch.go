package cmd

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
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
		statefile         string
		humanReadable     bool
		output            string
	)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Fetch audit events at regular intervals.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// create cancelling context based on signals
			ctx, cancel := context.WithCancel(context.Background())
			go func() {
				sigChan := getSigChan()
				for {
					select {
					case <-sigChan:
						cancel()
						return
					}
				}
			}()

			// Create state instance
			state := office365.NewGOBState()
			if statefile != "" {
				statefileAbs, writeStateDefer, err := setupStatefile(state, statefile)
				if err != nil {
					logger.Println(err)
					// ? TODO: Nested exit path
					return
				}
				defer writeStateDefer()
				logger.Printf("using statefile: %q\n", statefileAbs)
			}

			// Select output target
			// TODO: move into setup function
			filePrefix := "file://"
			udpPrefix := "udp://"
			tcpPrefix := "tcp://"

			var writer io.Writer
			switch {
			default:
				logger.Println("output invalid")
				return
			case output == "":
				writer = defaultOutput
			case strings.HasPrefix(output, filePrefix):
				path := output[len(filePrefix):len(output)]
				f, close, err := openOutputfile(path)
				if err != nil {
					logger.Printf("failed to create/open file: %s\n", err)
					// ? TODO: Nested exit path
					return
				}
				defer close()
				logger.Printf("using file output: %q\n", path)
				writer = f
			case strings.HasPrefix(output, udpPrefix):
				path := output[len(udpPrefix):len(output)]
				var d net.Dialer
				conn, err := d.DialContext(ctx, "udp", path)
				if err != nil {
					logger.Printf("Failed to dial udp: %s\n", err)
					// ? TODO: Nested exit path
					return
				}
				defer conn.Close()
			case strings.HasPrefix(output, tcpPrefix):
				path := output[len(tcpPrefix):len(output)]
				var d net.Dialer
				conn, err := d.DialContext(ctx, "tcp", path)
				if err != nil {
					logger.Printf("Failed to dial tcp: %s\n", err)
					// ? TODO: Nested exit path
					return
				}
				defer conn.Close()
			}

			// Select resource handler
			var handler office365.ResourceHandler
			if humanReadable {
				handler = office365.NewHumanReadableHandler(writer)
			} else {
				handler = office365.NewJSONHandler(writer, logger)
			}

			// Create client and launch watcher
			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)

			watcherConf := office365.SubscriptionWatcherConfig{
				LookBehindMinutes:     lookBehindMinutes,
				TickerIntervalSeconds: intervalSeconds,
			}
			if err := client.Subscription.Watch(ctx, watcherConf, state, handler); err != nil {
				logger.Printf("error occured calling watch: %s\n", err)
			}
		},
	}
	cmd.Flags().IntVar(&intervalSeconds, "interval", 5, "TickerIntervalSeconds")
	cmd.Flags().IntVar(&lookBehindMinutes, "lookbehind", 1, "Number of minutes from request time used when fetching available content.")
	cmd.Flags().StringVar(&statefile, "statefile", "", "File used to read/save state on start/exit.")
	cmd.Flags().BoolVar(&humanReadable, "human-readable", false, "Human readable output format.")
	cmd.Flags().StringVar(&output, "output", "", "Target where to send audit records. Available scheme: file://path/to/file, udp://1.2.3.4:1234, tcp://1.2.3.4:1234")

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

func openOutputfile(fpath string) (*os.File, func() error, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return nil, nil, err
	}
	return f, f.Close, nil
}

func setupStatefile(state *office365.GOBState, fpath string) (string, func() error, error) {
	statefile, err := filepath.Abs(fpath)
	if err != nil {
		return "", nil, fmt.Errorf("could not get absolute filepath for provided statefile: %s", err)
	}

	if err := readState(state, statefile); err != nil {
		return "", nil, fmt.Errorf("error occured setuping statefile: %s", err)
	}

	deferred := func() error {
		return writeState(state, statefile)
	}
	return statefile, deferred, nil
}

func openStatefile(fpath string) (*os.File, func() error, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, nil, err
	}
	return f, f.Close, nil
}

func readState(state *office365.GOBState, fpath string) error {
	f, close, err := openStatefile(fpath)
	if err != nil {
		return err
	}
	defer close()

	err = state.Read(f)
	if err != nil {
		logger.Println("state empty or invalid. Start fresh!")
	}
	return nil
}

func writeState(state *office365.GOBState, fpath string) error {
	f, close, err := openStatefile(fpath)
	if err != nil {
		return err
	}
	defer close()

	err = state.Write(f)
	if err != nil {
		return err
	}
	return nil
}
