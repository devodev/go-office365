package main

import (
	"context"
	"encoding/json"
	"fmt"

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
		Short: "Query content for the provided content-type.",
		Long:  fmt.Sprintf("Query content for the provided content-type.\n%s\n", timeArgsDescription),
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
			_, content, err := client.Content.List(context.Background(), ct, startTime, endTime)
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
	cmd.Flags().StringVar(&cfgFile, "config", "", "Set configfile alternate location. Defaults are [$HOME/.go-office365.yaml, $CWD/.go-office365.yaml].")
	cmd.Flags().StringVar(&startTime, "start", "", "Start time.")
	cmd.Flags().StringVar(&endTime, "end", "", "End time.")
	cmd.Flags().SortFlags = false
	return cmd
}
