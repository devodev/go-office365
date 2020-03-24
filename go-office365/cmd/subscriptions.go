package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/devodev/go-graph/office365"
	"github.com/spf13/cobra"
)

var (
	subscriptionsCmd = &cobra.Command{
		Use:   "subscriptions",
		Short: "List current subscriptions.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {

			// Micorosoft Office365 Management Activity Api Client
			client := office365.NewClientAuthenticated(&config.Credentials)

			pubIdentifier := ""

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
)

func init() {
	rootCmd.AddCommand(subscriptionsCmd)
}
