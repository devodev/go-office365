package main

import (
	"context"
	"fmt"

	"github.com/devodev/go-office365/v0/pkg/office365"
	"github.com/spf13/cobra"
)

func newCommandStopSub() *cobra.Command {
	var (
		cfgFile string
	)

	cmd := &cobra.Command{
		Use:   "stop-sub [content-type]",
		Short: "Stop a subscription for the provided Content Type.",
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

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)
			if err := client.Subscription.Stop(context.Background(), ct); err != nil {
				return err
			}
			writeOut("subscription successfully stopped")

			return nil
		},
	}
	cmd.Flags().StringVar(&cfgFile, "config", "", "config file")
	return cmd
}
