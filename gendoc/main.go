package main

import (
	"log"
	"path/filepath"

	"github.com/devodev/go-office365/v0/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	fpath, err := filepath.Abs("./docs")
	if err != nil {
		log.Fatal(err)
	}
	err = doc.GenMarkdownTree(cmd.RootCmd, fpath)
	if err != nil {
		log.Fatal(err)
	}
}
