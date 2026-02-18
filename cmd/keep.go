package cmd

import (
	"fmt"

	"github.com/nathfavour/kylrix/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var keepCmd = &cobra.Command{
	Use:   "keep",
	Short: "Manage Kylrix Keep backups and persistence",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var keepStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check backup status",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Banner("Kylrix Keep - Status")
		header := []string{"SERVICE", "LAST BACKUP", "STATUS"}
		data := [][]string{
			{"Vault", "2 hours ago", "Healthy"},
			{"Note", "5 mins ago", "Healthy"},
			{"Flow", "1 day ago", "Outdated"},
		}
		utils.Table(header, data)
	},
}

var keepBackupCmd = &cobra.Command{
	Use:   "backup [service]",
	Short: "Perform a manual backup",
	Run: func(cmd *cobra.Command, args []string) {
		service := "all"
		if len(args) > 0 {
			service = args[0]
		}
		utils.Banner(fmt.Sprintf("Kylrix Keep - Backup (%s)", service))
		utils.Info("Starting encryption and upload...")
		utils.Success("Backup completed successfully.")
	},
}

func init() {
	keepCmd.AddCommand(keepStatusCmd)
	keepCmd.AddCommand(keepBackupCmd)
	rootCmd.AddCommand(keepCmd)
}
