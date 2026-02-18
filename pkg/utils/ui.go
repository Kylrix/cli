package utils

import (
	"os"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/olekukonko/tablewriter"
)

func Success(msg string) {
	color.Green("✓ %s", msg)
}

func Error(msg string) {
	color.Red("✗ %s", msg)
}

func Info(msg string) {
	color.Blue("ℹ %s", msg)
}

func Warning(msg string) {
	color.Yellow("⚠ %s", msg)
}

func Banner(msg string) {
	color.Cyan("=== %s ===", msg)
}

func Table(header []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	// Some versions use AppendBulk, some use AddRow
	// We'll try the most basic approach
	for _, row := range data {
		table.Append(row)
	}
	table.Render()
}

func Prompt(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}
	return prompt.Run()
}

func PasswordPrompt(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
		Mask:  '*',
	}
	return prompt.Run()
}

func Select(label string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	return prompt.Run()
}
