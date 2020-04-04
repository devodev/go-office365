package cmd

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/devodev/go-office365/v0/office365"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(newCommandListSub())
	RootCmd.AddCommand(newCommandStartSub())
	RootCmd.AddCommand(newCommandStopSub())
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
				logger.Errorf("getting subscriptions: %s", err)
				return
			}
			for _, u := range subscriptions {
				payload, err := json.Marshal(u)
				if err != nil {
					logger.Errorf("marshalling subscriptions: %s", err)
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
				logger.Error("ContentType invalid")
				return
			}
			ct, err := office365.GetContentType(ctArg)
			if err != nil {
				logger.Error(err)
				return
			}

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier, logger)
			subscription, err := client.Subscription.Start(context.Background(), ct, nil)
			if err != nil {
				logger.Errorf("error getting subscriptions: %s", err)
				return
			}
			payload, err := json.Marshal(subscription)
			if err != nil {
				logger.Errorf("error marshalling subscription: %s", err)
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
				logger.Error("ContentType invalid")
				return
			}
			ct, err := office365.GetContentType(ctArg)
			if err != nil {
				logger.Error(err)
				return
			}

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier, logger)
			if err := client.Subscription.Stop(context.Background(), ct); err != nil {
				logger.Errorf("getting subscriptions: %s\n", err)
				return
			}
			WriteOut("subscription successfully stopped")
		},
	}
	return cmd
}
