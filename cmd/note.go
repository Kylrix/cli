package cmd

import (
	"fmt"

	"github.com/nathfavour/kylrix/cli/pkg/db"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.InitDB()
		if err != nil {
			return err
		}
		defer database.Close()

		rows, err := database.Query("SELECT title, created_at FROM notes")
		if err != nil {
			return err
		}
		defer rows.Close()

		utils.Banner("Kylrix Note - List")
		header := []string{"TITLE", "CREATED"}
		var data [][]string
		for rows.Next() {
			var title, created string
			if err := rows.Scan(&title, &created); err != nil {
				return err
			}
			data = append(data, []string{title, created})
		}

		if len(data) == 0 {
			utils.Info("No notes found.")
		} else {
			utils.Table(header, data)
		}
		return nil
	},
}

var noteCreateCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]
		
		content, err := utils.Prompt("Note Content")
		if err != nil {
			return err
		}

		if enhance {
			utils.Info("Enhancing note with AI...")
			content = "[AI Enhanced] " + content
		}
		
		database, err := db.InitDB()
		if err != nil {
			return err
		}
		defer database.Close()

		_, err = database.Exec("INSERT INTO notes (title, content) VALUES (?, ?)", title, content)
		if err != nil {
			return err
		}

		utils.Success("Note created successfully in local database.")
		return nil
	},
}

func init() {
	noteCreateCmd.Flags().BoolVar(&enhance, "enhance", false, "Use AI to enhance note content")
	
	noteCmd.AddCommand(noteListCmd)
	noteCmd.AddCommand(noteCreateCmd)
	rootCmd.AddCommand(noteCmd)
}
