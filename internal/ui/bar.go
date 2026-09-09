package ui

import (
	"math"
	"strings"
)

// BarStyle represents character sets for progress bars
type BarStyle struct {
	LeftBracket  string
	RightBracket string
	FillChar     string
	EmptyChar    string
}

var styles = map[string]BarStyle{
	"unicode": {
		LeftBracket:  "[",
		RightBracket: "]",
		FillChar:     "█",
		EmptyChar:    "░",
	},
	"blocks": {
		LeftBracket:  "[",
		RightBracket: "]",
		FillChar:     "■",
		EmptyChar:    "□",
	},
	"ascii": {
		LeftBracket:  "[",
		RightBracket: "]",
		FillChar:     "#",
		EmptyChar:    "-",
	},
	"smooth": {
		LeftBracket:  "[",
		RightBracket: "]",
		FillChar:     "▰",
		EmptyChar:    "▱",
	},
}

// RenderBar generates a styled, colored progress bar for a given percentage (0.0 to 100.0)
func RenderBar(perc float64, width int, styleName string, useColors bool) string {
	if width < 4 {
		width = 10
	}

	st, ok := styles[strings.ToLower(styleName)]
	if !ok {
		st = styles["smooth"]
	}

	// Bar content width excluding brackets
	barContentWidth := width - 2
	if barContentWidth < 1 {
		barContentWidth = 1
	}

	// Clamp percentage
	if perc < 0 {
		perc = 0
	}
	if perc > 100 {
		perc = 100
	}

	fillCount := int(math.Round((perc / 100.0) * float64(barContentWidth)))
	if fillCount > barContentWidth {
		fillCount = barContentWidth
	}
	emptyCount := barContentWidth - fillCount

	fillStr := strings.Repeat(st.FillChar, fillCount)
	emptyStr := strings.Repeat(st.EmptyChar, emptyCount)

	if !useColors {
		return st.LeftBracket + fillStr + emptyStr + st.RightBracket
	}

	barColor := UsageColor(perc)
	bracketColor := FgBrightBlack
	emptyColor := FgBrightBlack

	return bracketColor + st.LeftBracket + Reset +
		barColor + fillStr + Reset +
		emptyColor + emptyStr + Reset +
		bracketColor + st.RightBracket + Reset
}
