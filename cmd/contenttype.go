package cmd

import (
	"github.com/devodev/go-office365/v0/office365"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(newCommandContentType())
}

func newCommandContentType() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "content-types",
		Short: "List content types accepted by the Microsoft API.",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			contentTypes := office365.GetContentTypes()
			for _, v := range contentTypes {
				WriteOut(v.String())
			}
		},
	}

	return cmd
}
