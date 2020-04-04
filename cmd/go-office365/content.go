package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/devodev/go-office365/v0/pkg/office365"
	"github.com/spf13/cobra"
)

func newCommandContent() *cobra.Command {
	var (
		cfgFile   string
		startTime string
		endTime   string
	)

	cmd := &cobra.Command{
		Use:   "content [content-type]",
		Short: "List content that is available to be fetched for the provided content-type.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// command line args
			ctArg := args[0]

			// validate args
			if !office365.ContentTypeValid(ctArg) {
				return fmt.Errorf("ContentType invalid")
			}
			ct, err := office365.GetContentType(ctArg)
			if err != nil {
				return err
			}

			config, err := initConfig(cfgFile)
			if err != nil {
				return err
			}

			// parse optional args
			startTime := parseDate(startTime)
			endTime := parseDate(endTime)

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)
			content, err := client.Content.List(context.Background(), ct, startTime, endTime)
			if err != nil {
				return err
			}
			for _, u := range content {
				userData, err := json.Marshal(u)
				if err != nil {
					return err
				}
				writeOut(string(userData))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&cfgFile, "config", "", "config file")
	cmd.Flags().StringVar(&startTime, "start", "", "Start time")
	cmd.Flags().StringVar(&endTime, "end", "", "End time")

	return cmd
}

func parseDate(param string) time.Time {
	formats := []string{
		office365.RequestDateFormat,
		office365.RequestDatetimeFormat,
		office365.RequestDatetimeLargeFormat,
	}
	for _, format := range formats {
		parsed, err := time.Parse(format, param)
		if err == nil {
			return parsed
		}
	}
	return time.Time{}
}
