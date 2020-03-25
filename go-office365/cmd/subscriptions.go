package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/devodev/go-graph/office365"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCommandSubscriptions())
}

func newCommandSubscriptions() *cobra.Command {
	var pubIdentifier string

	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "List current subscriptions.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			// parse optional args
			if pubIdentifier == "" {
				pubIdentifier = config.Credentials.ClientID
			}

			client := office365.NewClientAuthenticated(&config.Credentials)
			subscriptions, err := client.Subscriptions.List(context.Background(), pubIdentifier)
			if err != nil {
				fmt.Printf("error getting subscriptions: %s\n", err)
				return
			}
			for _, u := range subscriptions {
				userData, err := json.Marshal(u)
				if err != nil {
					fmt.Printf("error marshalling subscriptions: %s\n", err)
					continue
				}
				fmt.Println(string(userData))
			}
		},
	}
	cmd.Flags().StringVar(&pubIdentifier, "identifier", "", "Publisher Identifier")

	return cmd
}
