package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/devodev/go-office365/office365"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	config  Config
	logger  *log.Logger

	loggerOutput  = os.Stderr
	defaultOutput = os.Stdout

	rootCmd = &cobra.Command{
		Use:     "go-office365",
		Short:   "Query the Microsoft Office365 Management Activity API.",
		Long:    "Query the Microsoft Office365 Management Activity API.",
		Version: "0.1",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// WriteOut .
func WriteOut(line string) {
	fmt.Fprintln(defaultOutput, line)
}

func init() {
	cobra.OnInitialize(initLogging, initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			logger.Fatalln(err)
		}

		viper.AddConfigPath(wd)
		viper.SetConfigName(".go-office365")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Println("Using config file:", viper.ConfigFileUsed())

	if err := viper.UnmarshalExact(&config); err != nil {
		logger.Fatalln(err)
	}
}

func initLogging() {
	logger = log.New(loggerOutput, "[go-office365] ", log.Flags())
}

// Config stores credentials and application
// specific attributes.
type Config struct {
	Global struct {
		Identifier string
	}
	Credentials office365.Credentials
}
