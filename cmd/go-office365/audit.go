package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/devodev/go-office365/v0/pkg/office365"
	"github.com/spf13/cobra"
)

func newCommandAudit() *cobra.Command {
	var (
		cfgFile         string
		extendedSchemas bool
	)

	cmd := &cobra.Command{
		Use:   "audit [audit-id]",
		Short: "Query audit records for the provided audit-id.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// command line args
			idArg := args[0]

			// validate args
			if idArg == "" {
				return fmt.Errorf("audit-id is empty")
			}

			config, err := initConfig(cfgFile)
			if err != nil {
				return err
			}

			client := office365.NewClientAuthenticated(&config.Credentials, config.Global.Identifier)
			_, audits, err := client.Audit.List(context.Background(), idArg, extendedSchemas)
			if err != nil {
				return err
			}
			for _, u := range audits {
				userData, err := json.Marshal(u)
				if err != nil {
					return err
				}
				writeOut(string(userData))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&cfgFile, "config", "", "Set configfile alternate location. Defaults are [$HOME/.go-office365.yaml, $CWD/.go-office365.yaml].")
	cmd.Flags().BoolVar(&extendedSchemas, "extended-schemas", false, "Set whether to add extended schemas to the output of the record or not.")
	cmd.Flags().SortFlags = false
	return cmd
}
