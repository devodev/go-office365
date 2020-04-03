package cmd

import (
	"fmt"
	"os"

	"github.com/devodev/go-office365/office365"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	debug       bool
	jsonLogging bool
)

var (
	config Config
	logger *logrus.Logger

	loggerOutput  = os.Stderr
	defaultOutput = os.Stdout

	// RootCmd made public so that gendoc can access it.
	RootCmd = &cobra.Command{
		Use:     "go-office365",
		Short:   "Query the Microsoft Office365 Management Activity API.",
		Long:    "Query the Microsoft Office365 Management Activity API.",
		Version: "0.1",
	}
)

// Execute executes the root command.
func Execute() error {
	return RootCmd.Execute()
}

// WriteOut .
func WriteOut(line string) {
	fmt.Fprintln(defaultOutput, line)
}

func init() {
	cobra.OnInitialize(initLogging, initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "set log level to DEBUG")
	RootCmd.PersistentFlags().BoolVar(&jsonLogging, "json", false, "set log formatter to JSON")
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
	logger.Infof("using config file: %s", viper.ConfigFileUsed())

	if err := viper.UnmarshalExact(&config); err != nil {
		logger.Fatalln(err)
	}
}

func initLogging() {
	logger = logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	if jsonLogging {
		logger.SetFormatter(&logrus.JSONFormatter{})
	}
	logger.SetLevel(logrus.InfoLevel)
	if debug {
		logger.SetLevel(logrus.DebugLevel)
	}
}

// Config stores credentials and application
// specific attributes.
type Config struct {
	Global struct {
		Identifier string
	}
	Credentials office365.Credentials
}
