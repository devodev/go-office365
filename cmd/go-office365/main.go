package main

import (
	"fmt"
	"os"

	"github.com/devodev/go-office365/v0/pkg/office365"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	loggerOutput  = os.Stderr
	defaultOutput = os.Stdout
)

// Execute executes the root command.
func Execute() error {
	rootCmd := newCommandRoot()
	return rootCmd.Execute()
}

func writeOut(line string) {
	fmt.Fprintln(defaultOutput, line)
}

func newCommandRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "go-office365",
		Short:   "Query the Microsoft Office365 Management Activity API.",
		Long:    "Query the Microsoft Office365 Management Activity API.",
		Version: "0.1.0-alpha.1",
	}
	cmd.AddCommand(
		newCommandAudit(),
		newCommandContent(),
		newCommandContentType(),
		newCommandFetch(),
		newCommandGenDoc(),
		newCommandListSub(),
		newCommandStartSub(),
		newCommandStopSub(),
		newCommandWatch(),
	)
	return cmd
}

func initConfig(cfgFile string) (*Config, error) {
	viperInstance := viper.New()
	if cfgFile != "" {
		viperInstance.SetConfigFile(cfgFile)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		hd, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}

		viperInstance.AddConfigPath(wd)
		viperInstance.AddConfigPath(hd)
		viperInstance.SetConfigName(".go-office365")
		viperInstance.SetConfigType("yaml")
	}

	viperInstance.AutomaticEnv()

	err := viperInstance.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var config Config
	if err := viperInstance.UnmarshalExact(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// Config stores credentials and application
// specific attributes.
type Config struct {
	Global struct {
		Identifier string
	}
	Credentials office365.Credentials
}

func main() {
	if err := Execute(); err != nil {
		writeOut(err.Error())
		os.Exit(1)
	}
}
