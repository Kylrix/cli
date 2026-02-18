package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "kylrix",
	Short: "Kylrix Ecosystem CLI",
	Long:  `A robust CLI tool for managing the Kylrix ecosystem (Note, Vault, Connect, Flow).`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		checkAnyisland()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to the Kylrix Ecosystem CLI")
		fmt.Println("Use 'kylrix --help' for available commands.")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kylrix/config.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func checkAnyisland() {
	home, _ := os.UserHomeDir()
	socketPath := home + "/.anyisland/anyisland.sock"

	if _, err := os.Stat(socketPath); os.IsNotExist(err) {
		return
	}

	conn, err := net.DialTimeout("unix", socketPath, 100*time.Millisecond)
	if err != nil {
		return
	}
	defer conn.Close()

	// Pulse Handshake
	handshake := map[string]string{"op": "HANDSHAKE"}
	json.NewEncoder(conn).Encode(handshake)

	var response map[string]interface{}
	if err := json.NewDecoder(conn).Decode(&response); err == nil {
		if response["status"] == "MANAGED" && verbose {
			fmt.Printf("[Anyisland] Managed as %v (%v)\n", response["tool_id"], response["version"])
		}
	}
}
