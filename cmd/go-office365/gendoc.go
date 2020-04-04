package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func newCommandGenDoc() *cobra.Command {
	var (
		dir string
	)

	cmd := &cobra.Command{
		Use:   "gendoc",
		Short: "Generate markdown documentation for the go-office365 CLI.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fpath, err := filepath.Abs(dir)
			if err != nil {
				return err
			}
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}

			if err := doc.GenMarkdownTree(newCommandRoot(), fpath); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&dir, "dir", "./docs", "directory where to write the doc.")
	return cmd
}
