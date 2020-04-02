package cmd

import (
	"context"
	"encoding/json"

	"github.com/devodev/go-office365/office365"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCommandFetch())
}

func newCommandFetch() *cobra.Command {
	var (
		startTime string
		endTime   string
	)

	cmd := &cobra.Command{
		Use:   "fetch [content-type]",
		Short: "Combination of content and audit commands.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// command line args
			ctArg := args[0]

			// validate args
			if !office365.ContentTypeValid(ctArg) {
				logger.Println("ContentType invalid")
				return
			}
			ct, err := office365.GetContentType(ctArg)
			if err != nil {
				logger.Println(err)
				return
			}

			// parse optional args
			startTime := parseDate(startTime)
			endTime := parseDate(endTime)

			// Create client
			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier, logger)

			// retrieve content
			content, err := client.Content.List(context.Background(), ct, startTime, endTime)
			if err != nil {
				logger.Printf("error getting content: %s\n", err)
				return
			}

			// retrieve audits
			var auditList []office365.AuditRecord
			for _, c := range content {
				audits, err := client.Audit.List(context.Background(), c.ContentID)
				if err != nil {
					logger.Printf("error getting audits: %s\n", err)
					continue
				}
				auditList = append(auditList, audits...)
			}

			// output
			for _, a := range auditList {
				auditStr, err := json.Marshal(a)
				if err != nil {
					logger.Printf("error marshalling audit: %s\n", err)
					continue
				}
				WriteOut(string(auditStr))
			}

		},
	}
	cmd.Flags().StringVar(&startTime, "start", "", "Start time")
	cmd.Flags().StringVar(&endTime, "end", "", "End time")

	return cmd
}
