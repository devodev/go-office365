package cmd

import (
	"fmt"
	"os"

	"github.com/devodev/go-graph/office365"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	config  Config

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

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(wd)
		viper.SetConfigName(".go-office365")
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Using config file:", viper.ConfigFileUsed())

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Config stores credentials and application
// specific attributes.
type Config struct {
	Credentials office365.Credentials
}
