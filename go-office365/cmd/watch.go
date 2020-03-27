package cmd

import (
	"context"
	"encoding/json"
	"fmt"
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
			// command line args
			confArg := args[0]

			// load config file
			watchConfig, err := loadConfig(confArg)
			if err != nil {
				fmt.Printf("error occured loading config file: %s\n", err)
				return
			}

			// create office365 client
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

			resultChan := client.Subscriptions.Watch(ctx, watchConfig.Global.FetcherCount, watchConfig.Global.TickerIntervalMinutes)
			printer(resultChan)
		},
	}
	return cmd
}

func printer(in <-chan office365.Resource) {
	for r := range in {
		for _, a := range r.Records {
			auditStr, err := json.Marshal(a)
			if err != nil {
				fmt.Printf("error marshalling audit: %s\n", err)
				continue
			}
			fmt.Println(string(auditStr))
		}
	}
}

// WatchConfig .
type WatchConfig struct {
	Global struct {
		TickerIntervalMinutes int
		FetcherCount          int
		PubIdentifier         string
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
