package ui

import (
	"os"

	"golang.org/x/term"
)

// ANSI Color Codes
const (
	Reset       = "\033[0m"
	Bold        = "\033[1m"
	Dim         = "\033[2m"
	Italic      = "\033[3m"
	Underline   = "\033[4m"

	// Foreground colors
	FgBlack   = "\033[30m"
	FgRed     = "\033[31m"
	FgGreen   = "\033[32m"
	FgYellow  = "\033[33m"
	FgBlue    = "\033[34m"
	FgMagenta = "\033[35m"
	FgCyan    = "\033[36m"
	FgWhite   = "\033[37m"

	// Bright foreground colors
	FgBrightBlack   = "\033[90m" // Gray
	FgBrightRed     = "\033[91m"
	FgBrightGreen   = "\033[92m"
	FgBrightYellow  = "\033[93m"
	FgBrightBlue    = "\033[94m"
	FgBrightMagenta = "\033[95m"
	FgBrightCyan    = "\033[96m"
	FgBrightWhite   = "\033[97m"

	// Background colors
	BgRed   = "\033[41m"
	BgGreen = "\033[42m"
)

// ColorMode determines if colors should be applied
type ColorMode bool

// IsTerminal returns true if stdout is an interactive terminal
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// ShouldUseColor checks if colors should be enabled based on user config, env vars, and tty
func ShouldUseColor(userPref bool) bool {
	if !userPref {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("TERM") == "dumb" {
		return false
	}
	return IsTerminal()
}

// Colorize wraps text in ANSI color codes if enabled
func Colorize(text string, color string, enabled bool) string {
	if !enabled || color == "" {
		return text
	}
	return color + text + Reset
}

// UsageColor returns the appropriate color for a given usage percentage
func UsageColor(perc float64) string {
	if perc >= 90.0 {
		return FgBrightRed + Bold
	} else if perc >= 75.0 {
		return FgBrightYellow
	}
	return FgBrightGreen
}

// GetTerminalWidth returns the width of the terminal or a default of 80
func GetTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		return 80
	}
	return width
}
