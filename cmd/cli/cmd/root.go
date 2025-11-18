/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

// Request - struct representing the structure of the HTTP request
type Request struct {
	ID      string            `json:"id"`
	Path    string            `json:"path"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// structure of the response being sent
type Response struct {
	ID      string            `json:"id"`
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

var PORT int
var SERVER_URL = "ws://localhost:8080/connect"

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

		conn, _, err := websocket.DefaultDialer.Dial(SERVER_URL, nil)
		if err != nil {
			log.Fatal("Failed to connect to server", err)
		}
		defer conn.Close()

		// recieve the tunnel ID
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatal("Failed to read from the tunnel", err)
		}
		tunnelID := string(message[len("Tunnel established: "):])
		publicURL := fmt.Sprintf("http://localhost:8080/%s", tunnelID)
		fmt.Printf("Tunnel Established! Public URL: %s \n", publicURL)

		//  We are listening for requests
		for {
			_, message, err := conn.ReadMessage()

			if err != nil {
				log.Fatal("Websocket Read", err)
				break
			}

			// var Request - we are parsing it to match the request
			var request Request
			if err := json.Unmarshal(message, &request); err != nil {
				log.Println("Invalid Request", err)
				continue
			}

			// Forwarding our request to the localhost
			localHostURL := fmt.Sprintf("http://localhost:%d%s",)
		}

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
