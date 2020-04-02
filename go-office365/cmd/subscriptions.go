package cmd

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/devodev/go-office365/office365"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(newCommandListSub())
	rootCmd.AddCommand(newCommandStartSub())
	rootCmd.AddCommand(newCommandStopSub())
}

func newCommandListSub() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "List current subscriptions.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier, logger)
			subscriptions, err := client.Subscription.List(context.Background())
			if err != nil {
				logger.Printf("error getting subscriptions: %s\n", err)
				return
			}
			for _, u := range subscriptions {
				payload, err := json.Marshal(u)
				if err != nil {
					logger.Printf("error marshalling subscriptions: %s\n", err)
					continue
				}
				var out bytes.Buffer
				json.Indent(&out, payload, "", "\t")
				WriteOut(out.String())
			}
		},
	}
	return cmd
}

func newCommandStartSub() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start-sub [content-type]",
		Short: "Start a subscription for the provided Content Type.",
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

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier, logger)
			subscription, err := client.Subscription.Start(context.Background(), ct, nil)
			if err != nil {
				logger.Printf("error getting subscriptions: %s\n", err)
				return
			}
			payload, err := json.Marshal(subscription)
			if err != nil {
				logger.Printf("error marshalling subscription: %s\n", err)
				return
			}
			var out bytes.Buffer
			json.Indent(&out, payload, "", "\t")
			WriteOut(out.String())

			WriteOut("subscription successfully started")
		},
	}
	return cmd
}
func newCommandStopSub() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop-sub [content-type]",
		Short: "Stop a subscription for the provided Content Type.",
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

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier, logger)
			if err := client.Subscription.Stop(context.Background(), ct); err != nil {
				logger.Printf("error getting subscriptions: %s\n", err)
				return
			}

			WriteOut("subscription successfully stopped")
		},
	}
	return cmd
}
