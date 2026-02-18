package cmd

import (
	"fmt"

	"github.com/nathfavour/kylrix/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	readOnly  bool
	readWrite bool
	targetTo  string
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Manage Kylrix Connect P2P communications",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share a session with another user",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Banner("Kylrix Connect - Share")
		mode := "read-only"
		if readWrite {
			mode = "read-write"
		}
		
		msg := fmt.Sprintf("Sharing session in %s mode", mode)
		if targetTo != "" {
			msg += fmt.Sprintf(" with user: %s", targetTo)
		}
		utils.Info(msg)
		utils.Success("Sharing session started at https://kylrix.com/shared/xyz-123")
	},
}

var joinCmd = &cobra.Command{
	Use:   "join [id]",
	Short: "Join a shared session",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sessionID := args[0]
		utils.Banner("Kylrix Connect - Join")
		utils.Info(fmt.Sprintf("Joining session %s...", sessionID))
		utils.Success("Connected to session.")
	},
}

func init() {
	shareCmd.Flags().BoolVar(&readOnly, "ro", true, "Share in read-only mode")
	shareCmd.Flags().BoolVar(&readWrite, "rw", false, "Share in read-write mode")
	shareCmd.Flags().StringVar(&targetTo, "to", "", "Target user ID for sharing")
	
	connectCmd.AddCommand(shareCmd)
	connectCmd.AddCommand(joinCmd)
	rootCmd.AddCommand(connectCmd)
}
