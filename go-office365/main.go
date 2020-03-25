package main

import (
	"os"

	"github.com/devodev/go-office365/go-office365/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
