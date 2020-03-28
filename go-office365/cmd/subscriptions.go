package cmd

import (
	"context"
	"encoding/json"

	"github.com/devodev/go-office365/office365"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCommandSubscriptions())
}

func newCommandSubscriptions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "List current subscriptions.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)
			subscriptions, err := client.Subscriptions.List(context.Background())
			if err != nil {
				logger.Printf("error getting subscriptions: %s\n", err)
				return
			}
			for _, u := range subscriptions {
				userData, err := json.Marshal(u)
				if err != nil {
					logger.Printf("error marshalling subscriptions: %s\n", err)
					continue
				}
				WriteOut(string(userData))
			}
		},
	}
	return cmd
}
