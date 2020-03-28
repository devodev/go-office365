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

			resultChan, err := client.Subscriptions.Watch(ctx,
				watchConfig.Global.FetcherCount,
				watchConfig.Global.FetcherLookBehindMinutes,
				watchConfig.Global.TickerIntervalSeconds)
			if err != nil {
				logger.Printf("error occured calling watch: %s\n", err)
				return
			}
			printer(resultChan)
		},
	}
	return cmd
}

func printer(in <-chan office365.Resource) {
	for r := range in {
		for idx, e := range r.Errors {
			WriteOut(fmt.Sprintf("[%s] Error%d: %s\n", r.Request.ContentType, idx, e.Error()))
		}
		for _, a := range r.Response.Records {
			auditStr, err := json.Marshal(a)
			if err != nil {
				logger.Printf("error marshalling audit: %s\n", err)
				continue
			}
			WriteOut(fmt.Sprintf("[%s] %s\n", r.Request.ContentType, string(auditStr)))
		}
	}
}

// WatchConfig .
type WatchConfig struct {
	Global struct {
		TickerIntervalSeconds    int
		FetcherCount             int
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
