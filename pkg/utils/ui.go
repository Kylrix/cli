package utils

import (
	"github.com/fatih/color"
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
