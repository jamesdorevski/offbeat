package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "offbeat",
	Short: "Manage your Tempo worklogs",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	viper.SetConfigName("offbeat")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/offbeat")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Missing config file.")
		panic(err)
	}
	fmt.Println("Using config file: ", viper.ConfigFileUsed())
}