package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/devodev/go-office365/v0/pkg/office365"
	"github.com/devodev/go-office365/v0/pkg/office365/schema"
	"github.com/spf13/cobra"
)

func newCommandStartSub() *cobra.Command {
	var (
		cfgFile string
	)

	cmd := &cobra.Command{
		Use:   "start-sub [content-type]",
		Short: "Start a subscription for the provided Content Type.",
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

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)
			_, subscription, err := client.Subscription.Start(context.Background(), ct, nil)
			if err != nil {
				return err
			}
			payload, err := json.Marshal(subscription)
			if err != nil {
				return err
			}
			var out bytes.Buffer
			json.Indent(&out, payload, "", "\t")
			writeOut(out.String())
			writeOut("subscription successfully started")

			return nil
		},
	}
	cmd.Flags().StringVar(&cfgFile, "config", "", "Set configfile alternate location. Defaults are [$HOME/.go-office365.yaml, $CWD/.go-office365.yaml].")
	cmd.Flags().SortFlags = false
	return cmd
}
