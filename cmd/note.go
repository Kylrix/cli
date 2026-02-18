package cmd

import (
	"fmt"

	"github.com/nathfavour/kylrix/cli/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	enhance bool
)

var noteCmd = &cobra.Command{
	Use:   "note",
	Short: "Manage Kylrix Notes",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var noteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all notes",
	Run: func(cmd *cobra.Command, args []string) {
		utils.Banner("Kylrix Note - List")
		// Placeholder for actual API call
		notes := []string{"Personal Goals", "Work Project Alpha", "Quick Thoughts"}
		for i, note := range notes {
			fmt.Printf("%d. %s\n", i+1, note)
		}
	},
}

var noteCreateCmd = &cobra.Command{
	Use:   "create [title] [content]",
	Short: "Create a new note",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		title := args[0]
		content := ""
		if len(args) > 1 {
			content = args[1]
		}

		utils.Banner("Kylrix Note - Create")
		if enhance {
			utils.Info("Enhancing note with AI...")
			content = "[AI Enhanced] " + content
		}
		
		utils.Info(fmt.Sprintf("Creating note: %s", title))
		utils.Success("Note created successfully.")
		fmt.Printf("Content: %s\n", content)
	},
}

func init() {
	noteCreateCmd.Flags().BoolVar(&enhance, "enhance", false, "Use AI to enhance note content")
	
	noteCmd.AddCommand(noteListCmd)
	noteCmd.AddCommand(noteCreateCmd)
	rootCmd.AddCommand(noteCmd)
}
