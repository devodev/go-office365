package main

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/devodev/go-office365/v0/pkg/office365"
	"github.com/spf13/cobra"
)

func newCommandListSub() *cobra.Command {
	var (
		cfgFile string
	)

	cmd := &cobra.Command{
		Use:   "subscriptions",
		Short: "List current subscriptions.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := initConfig(cfgFile)
			if err != nil {
				return err
			}

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)
			_, subscriptions, err := client.Subscription.List(context.Background())
			if err != nil {
				return err
			}
			for _, u := range subscriptions {
				payload, err := json.Marshal(u)
				if err != nil {
					return err
				}
				var out bytes.Buffer
				json.Indent(&out, payload, "", "\t")
				writeOut(out.String())
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&cfgFile, "config", "", "Set configfile alternate location. Defaults are [$HOME/.go-office365.yaml, $CWD/.go-office365.yaml].")
	cmd.Flags().SortFlags = false
	return cmd
}
