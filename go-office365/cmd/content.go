package cmd

import (
	"context"
	"encoding/json"
	"time"

	"github.com/devodev/go-office365/v0/office365"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(newCommandContent())
}

func newCommandContent() *cobra.Command {
	var (
		startTime string
		endTime   string
	)

	cmd := &cobra.Command{
		Use:   "content [content-type]",
		Short: "List content that is available to be fetched for the provided content-type.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// command line args
			ctArg := args[0]

			// validate args
			if !office365.ContentTypeValid(ctArg) {
				logger.Error("ContentType invalid")
				return
			}
			ct, err := office365.GetContentType(ctArg)
			if err != nil {
				logger.Error(err)
				return
			}

			// parse optional args
			startTime := parseDate(startTime)
			endTime := parseDate(endTime)

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier, logger)
			content, err := client.Content.List(context.Background(), ct, startTime, endTime)
			if err != nil {
				logger.Errorf("getting content: %s", err)
				return
			}
			for _, u := range content {
				userData, err := json.Marshal(u)
				if err != nil {
					logger.Errorf("marshalling content: %s", err)
					continue
				}
				WriteOut(string(userData))
			}
		},
	}
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
