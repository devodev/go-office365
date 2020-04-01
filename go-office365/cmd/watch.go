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
		output            string
		humanReadable     bool
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
			writer, close, err := setupOutput(ctx, output)
			if err != nil {
				logger.Println(err)
				// ? TODO: Nested exit path
				return
			}
			defer close()
			if output != "" {
				logger.Printf("using output: %q\n", output)
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
	cmd.Flags().StringVar(&output, "output", "", "Target where to send audit records. Available schemes: file://path/to/file, udp://1.2.3.4:1234, tcp://1.2.3.4:1234")
	cmd.Flags().BoolVar(&humanReadable, "human-readable", false, "Human readable output format.")

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

func setupOutput(ctx context.Context, selection string) (io.Writer, func() error, error) {
	var writer io.Writer
	var deferred func() error

	filePrefix := "file://"
	udpPrefix := "udp://"
	tcpPrefix := "tcp://"

	switch {
	default:
		return nil, nil, fmt.Errorf("output invalid")
	case selection == "":
		writer = defaultOutput
		deferred = func() error { return nil }
	case strings.HasPrefix(selection, filePrefix):
		path := selection[len(filePrefix):len(selection)]
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, nil, fmt.Errorf("could not get absolute filepath for provided statefile: %s", err)
		}
		f, close, err := openOutputfile(path)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create/open file: %s", err)
		}
		writer = f
		deferred = close
	case strings.HasPrefix(selection, udpPrefix):
		path := selection[len(udpPrefix):len(selection)]
		var d net.Dialer
		conn, err := d.DialContext(ctx, "udp", path)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to dial udp: %s", err)
		}
		writer = conn
		deferred = conn.Close
	case strings.HasPrefix(selection, tcpPrefix):
		path := selection[len(tcpPrefix):len(selection)]
		var d net.Dialer
		conn, err := d.DialContext(ctx, "tcp", path)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to dial tcp: %s", err)
		}
		writer = conn
		deferred = conn.Close
	}
	return writer, deferred, nil
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
