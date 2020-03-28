package cmd

import (
	"context"
	"encoding/json"

	"github.com/devodev/go-office365/office365"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCommandAudit())
}

func newCommandAudit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "audit [audit-id]",
		Short: "Retrieve events and/or actions for the provided audit-id.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			// command line args
			idArg := args[0]

			// validate args
			if idArg == "" {
				logger.Println("audit-id is empty")
				return
			}

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)
			audits, err := client.Subscriptions.Audit(context.Background(), idArg)
			if err != nil {
				logger.Printf("error getting audits: %s\n", err)
				return
			}
			for _, u := range audits {
				userData, err := json.Marshal(u)
				if err != nil {
					logger.Printf("error marshalling audits: %s\n", err)
					continue
				}
				WriteOut(string(userData))
			}
		},
	}
	return cmd
}
