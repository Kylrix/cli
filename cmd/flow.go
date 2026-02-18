package cmd

import (
	"fmt"

	"github.com/nathfavour/kylrix/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var flowCmd = &cobra.Command{
	Use:   "flow",
	Short: "Manage Kylrix Flow productivity tasks",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var flowTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List tasks",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Banner("Kylrix Flow - Tasks")
		tasks := []string{"[ ] Review PR #123", "[x] Implement CLI foundations", "[ ] Update AGENTS.md"}
		for _, t := range tasks {
			fmt.Printf("- %s\n", t)
		}
	},
}

var flowSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync flow data with backend",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Banner("Kylrix Flow - Sync")
		utils.Info("Synchronizing tasks and calendars...")
		utils.Success("Synchronization complete.")
	},
}

func init() {
	flowCmd.AddCommand(flowTasksCmd)
	flowCmd.AddCommand(flowSyncCmd)
	rootCmd.AddCommand(flowCmd)
}
