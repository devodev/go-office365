package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// create office365 client
			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)

			resourceChan := make(chan *Resource)
			resultChan := make(chan *ResourceResponse)
			defer close(resultChan)

			for i := 0; i < 3; i++ {
				go fetcher(ctx, client, resourceChan, resultChan)
			}
			go printer(resultChan)

			var wg sync.WaitGroup
			wg.Add(1)

			go main(ctx, client, watchConfig.Global.TickerIntervalMinutes, resourceChan, &wg)

			wg.Wait()
		},
	}
	return cmd
}

func main(ctx context.Context, o365Client *office365.Client, intervalMinutes int, out chan *Resource, wg *sync.WaitGroup) {
	sigChan := getSigChan()

	// TODO: change time.Second into time.Minute. This is to ease testing.
	tickerDur := time.Duration(intervalMinutes) * time.Second
	ticker := time.NewTicker(tickerDur)

	for {
		select {
		case <-sigChan:
			wg.Done()
			return
		case t := <-ticker.C:
			subscriptions, err := o365Client.Subscriptions.List(ctx)
			if err != nil {
				fmt.Printf("error getting subscriptions: %s\n", err)
				break
			}

			// TODO: remove time.Minute
			startTime := t.Add(-(tickerDur + time.Minute))
			endTime := t

			for _, s := range subscriptions {
				ct, err := office365.GetContentType(s.ContentType)
				if err != nil {
					fmt.Println(err)
					continue
				}
				resource := &Resource{
					contentType: ct,
					startTime:   startTime,
					endTime:     endTime,
				}
				out <- resource
			}
		}
	}
}

func fetcher(ctx context.Context, client *office365.Client, in <-chan *Resource, out chan *ResourceResponse) {
	for r := range in {
		content, err := client.Subscriptions.Content(ctx, r.contentType, r.startTime, r.endTime)
		if err != nil {
			fmt.Printf("error getting content: %s\n", err)
			continue
		}

		var auditList []office365.AuditRecord
		for _, c := range content {
			audits, err := client.Subscriptions.Audit(ctx, c.ContentID)
			if err != nil {
				fmt.Printf("error getting audits: %s\n", err)
				continue
			}
			auditList = append(auditList, audits...)
		}
		out <- &ResourceResponse{Records: auditList}
	}

}

func printer(in <-chan *ResourceResponse) {
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
		PubIdentifier         string
	}
}

// Resource .
type Resource struct {
	contentType *office365.ContentType
	startTime   time.Time
	endTime     time.Time
}

// ResourceResponse .
type ResourceResponse struct {
	Records []office365.AuditRecord
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
