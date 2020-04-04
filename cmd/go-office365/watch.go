package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/devodev/go-office365/v0/pkg/office365"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	errInvalidStatefile = errors.New("statefile content empty or invalid, starting fresh")
)

func newCommandWatch() *cobra.Command {
	var (
		logFile   string
		cfgFile   string
		stateFile string

		intervalSeconds   int
		lookBehindMinutes int
		output            string
		indent            bool
		debug             bool
		jsonLogging       bool
	)

	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Query audit records at regular intervals.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// init logger and config
			logger, err := initLogger(cmd, logFile, debug, jsonLogging)
			if err != nil {
				return err
			}
			config, err := initConfig(cfgFile)
			if err != nil {
				return err
			}

			// create cancelling context using signals
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

			// create state instance
			state := office365.NewMemoryState()
			if stateFile != "" {
				statefileAbs, writeStateDefer, err := setupStatefile(state, stateFile)
				if err != nil {
					if err != errInvalidStatefile {
						return err
					}
					logger.Info(err)
				}
				defer writeStateDefer()
				logger.Infof("using statefile: %s", statefileAbs)
			}

			// setup output target
			writer, close, err := setupOutput(ctx, output)
			if err != nil {
				return err
			}
			defer close()
			if output != "" {
				logger.Infof("using output: %s", output)
			}

			// create watcher and start it
			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)
			handler := office365.NewJSONHandler(writer, logger, indent)

			watcherConf := office365.SubscriptionWatcherConfig{
				LookBehindMinutes:     lookBehindMinutes,
				TickerIntervalSeconds: intervalSeconds,
			}
			watcher, err := office365.NewSubscriptionWatcher(client, watcherConf, state, handler, logger)
			if err != nil {
				return err
			}
			return watcher.Run(ctx)
		},
	}
	cmd.Flags().StringVar(&cfgFile, "config", "", "Set configfile alternate location. Defaults are [$HOME/.go-office365.yaml, $CWD/.go-office365.yaml].")

	cmd.Flags().StringVar(&logFile, "log", "", "Set logging output to provided file. Default is stderr.")
	cmd.Flags().StringVar(&stateFile, "state", "", "Set state output to provided file. Default is to not persist state.")
	cmd.Flags().StringVar(&output, "output", "", "Set records output. Available schemes: file://path/to/file, udp://1.2.3.4:1234, tcp://1.2.3.4:1234")

	cmd.Flags().IntVar(&intervalSeconds, "interval", 5, "Ticker interval used to trigger fetch pipelines, in second(s).")
	cmd.Flags().IntVar(&lookBehindMinutes, "lookbehind", 1, "Minimum interval used by fetch actions, in minute(s).")
	cmd.Flags().BoolVar(&indent, "indent", false, "Set records output to be indented.")
	cmd.Flags().BoolVar(&debug, "debug", false, "Set log level to DEBUG.")
	cmd.Flags().BoolVar(&jsonLogging, "json", false, "Set log formatter to JSON.")
	cmd.Flags().SortFlags = false
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

func initLogger(cmd *cobra.Command, logFile string, setDebug, setJSON bool) (*logrus.Logger, error) {
	logger := logrus.New()

	logger.SetLevel(logrus.InfoLevel)
	if setDebug {
		logger.SetLevel(logrus.DebugLevel)
	}

	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		DisableSorting:         true,
	})
	if setJSON {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}

	logger.SetOutput(loggerOutput)
	if logFile != "" {
		logFile, err := filepath.Abs(logFile)
		if err != nil {
			return nil, fmt.Errorf("could not get absolute filepath for provided logfile: %s", err)
		}
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
		if err != nil {
			return nil, fmt.Errorf("could not use provided logfile: %s", err)
		}
		logger.SetOutput(f)
		cmd.PostRunE = func(cmd *cobra.Command, args []string) error {
			return f.Close()
		}
	}
	return logger, nil
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
		path := selection[len(filePrefix):]
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
		path := selection[len(udpPrefix):]
		var d net.Dialer
		conn, err := d.DialContext(ctx, "udp", path)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to dial udp: %s", err)
		}
		writer = conn
		deferred = conn.Close
	case strings.HasPrefix(selection, tcpPrefix):
		path := selection[len(tcpPrefix):]
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
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	if err != nil {
		return nil, nil, err
	}
	return f, f.Close, nil
}

func setupStatefile(state *office365.MemoryState, fpath string) (string, func() error, error) {
	statefile, err := filepath.Abs(fpath)
	if err != nil {
		return "", nil, fmt.Errorf("could not get absolute filepath for provided statefile: %s", err)
	}

	deferred := func() error {
		return writeState(state, statefile)
	}
	// readstate could return errInvalidStatefile
	// which we handle gracefully, therefore,
	// just send down the err and let the caller handle it
	err = readState(state, statefile)

	return statefile, deferred, err
}

func openStatefile(fpath string) (*os.File, func() error, error) {
	f, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return nil, nil, err
	}
	return f, f.Close, nil
}

func readState(state *office365.MemoryState, fpath string) error {
	f, close, err := openStatefile(fpath)
	if err != nil {
		return err
	}
	defer close()

	err = state.Read(f)
	if err != nil {
		return errInvalidStatefile
	}
	return nil
}

func writeState(state *office365.MemoryState, fpath string) error {
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
