package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/devodev/go-graph/office365"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(newCommandWatch())
}

func newCommandWatch() *cobra.Command {
	var (
		pubIdentifier string
	)

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

			// parse optional args
			if pubIdentifier == "" {
				pubIdentifier = config.Credentials.ClientID
			}

			// setup ticker using config interval
			// TODO: change time.Second into time.Minute. This is to ease testing.
			tickerDur := time.Duration(watchConfig.Global.TickerIntervalMinutes) * time.Second
			ticker := time.NewTicker(tickerDur)

			// setup signal handling
			sigChan := getSigChan()

			// create office365 client
			client := office365.NewClientAuthenticated(&config.Credentials)

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			doneChan := make(chan bool)

			// main loop
			go main(ctx, client, pubIdentifier, ticker, tickerDur, doneChan, sigChan)

			<-doneChan
		},
	}
	cmd.Flags().StringVar(&pubIdentifier, "identifier", "", "Publisher Identifier")

	return cmd
}

// WatchConfig .
type WatchConfig struct {
	Global struct {
		TickerIntervalMinutes int
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

func main(ctx context.Context, o365Client *office365.Client, pubIdentifier string, tc *time.Ticker, tDur time.Duration, done chan bool, exit <-chan os.Signal) {
	for {
		select {
		case <-exit:
			done <- true
			return
		case t := <-tc.C:

			subscriptions, err := o365Client.Subscriptions.List(ctx, pubIdentifier)
			if err != nil {
				fmt.Printf("error getting subscriptions: %s\n", err)
				break
			}

			// TODO: remove time.Minute
			startTime := t.Add(-(tDur + time.Minute))
			endTime := t

			for _, s := range subscriptions {

				ct, err := office365.GetContentType(s.ContentType)
				if err != nil {
					fmt.Println(err)
					continue
				}

				queue := make(chan []office365.AuditRecord)

				go fetcher(ctx, o365Client, pubIdentifier, ct, startTime, endTime, queue)
				go printer(queue)
			}
		}
	}
}

func fetcher(ctx context.Context, o365Client *office365.Client, pubIdentifier string, ct *office365.ContentType, start time.Time, end time.Time, queue chan []office365.AuditRecord) {
	content, err := o365Client.Subscriptions.Content(ctx, pubIdentifier, ct, start, end)
	if err != nil {
		fmt.Printf("error getting content: %s\n", err)
		return
	}

	var auditList []office365.AuditRecord
	for _, c := range content {
		audits, err := o365Client.Subscriptions.Audit(ctx, c.ContentID)
		if err != nil {
			fmt.Printf("error getting audits: %s\n", err)
			continue
		}
		auditList = append(auditList, audits...)
	}

	queue <- auditList
	close(queue)
}

func printer(queue <-chan []office365.AuditRecord) {
	for audits := range queue {
		for _, a := range audits {
			auditStr, err := json.Marshal(a)
			if err != nil {
				fmt.Printf("error marshalling audit: %s\n", err)
				continue
			}
			fmt.Println(string(auditStr))
		}
	}
}
