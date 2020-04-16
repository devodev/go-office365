package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/devodev/go-office365/v0/pkg/office365"
	"github.com/devodev/go-office365/v0/pkg/office365/schema"
	"github.com/spf13/cobra"
)

func newCommandFetch() *cobra.Command {
	var (
		cfgFile   string
		startTime string
		endTime   string
	)

	cmd := &cobra.Command{
		Use:   "fetch [content-type]",
		Short: "Query audit records for the provided content-type.",
		Long:  fmt.Sprintf("Query audit records for the provided content-type.\n%s\n", timeArgsDescription),
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// command line args
			ctArg := args[0]

			// validate args
			if !schema.ContentTypeValid(ctArg) {
				return fmt.Errorf("ContentType invalid")
			}
			ct, err := schema.GetContentType(ctArg)
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

			// Create client
			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)

			// retrieve content
			_, content, err := client.Content.List(context.Background(), ct, startTime, endTime)
			if err != nil {
				return err
			}

			// retrieve audits
			var auditList []interface{}
			for _, c := range content {
				_, audits, err := client.Audit.List(context.Background(), c.ContentID, true)
				if err != nil {
					return err
				}
				auditList = append(auditList, audits...)
			}

			// output
			for _, a := range auditList {
				auditStr, err := json.Marshal(a)
				if err != nil {
					return err
				}
				writeOut(string(auditStr))
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
