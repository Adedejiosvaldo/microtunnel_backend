/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var PORT int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mt",
	Short: "A tunneling solution",
	Long: `
Microtunnel is a cli application that
gives developers a public webhook endpoint
that they can inspect, replay, and forward to
localhost — without needing ngrok`,

	Run: func(cmd *cobra.Command, args []string) {
		id := createUUID()
		fmt.Println(id)
		fmt.Printf("Hola %s", args[0])

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.micro_tunnel.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().IntVar(&PORT, "Port Number", 8000, "The port to tunnel to")
}
