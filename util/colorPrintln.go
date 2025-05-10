package util

import "github.com/fatih/color"

var (
	SuccessColor = color.New(color.FgGreen, color.Bold)
	ErrorColor   = color.New(color.FgRed, color.Bold)
	WarningColor = color.New(color.FgYellow, color.Bold)
	InfoColor    = color.New(color.FgCyan, color.Bold)
)
