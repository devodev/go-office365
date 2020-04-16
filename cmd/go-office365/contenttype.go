package main

import (
	"github.com/devodev/go-office365/v0/pkg/office365/schema"
	"github.com/spf13/cobra"
)

func newCommandContentType() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content-types",
		Short: "List content types accepted by the Microsoft API.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			contentTypes := schema.GetContentTypes()
			for _, v := range contentTypes {
				writeOut(v.String())
			}
			return nil
		},
	}
	return cmd
}
