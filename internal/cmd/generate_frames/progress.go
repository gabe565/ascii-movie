package main

import (
	"github.com/fatih/color"
	_ "github.com/fatih/color"
	"math"
	"strings"
)

var colorLightBlack = color.New(38).SprintFunc()

func progressBar(n, total, width int) string {
	percent := float64(n) / float64(total)
	fullWidth := percent * float64(width)
	var partChar rune
	switch int(math.Mod(fullWidth, 1.0) * 8) {
	case 0:
		partChar = ' '
	case 1:
		partChar = '▏'
	case 2:
		partChar = '▎'
	case 3:
		partChar = '▍'
	case 4:
		partChar = '▌'
	case 5:
		partChar = '▋'
	case 6:
		partChar = '▊'
	case 7:
		partChar = '▉'
	default:
		partChar = '█'
	}
	return colorLightBlack(
		"[" +
			strings.Repeat("█", int(fullWidth)) +
			string(partChar) +
			strings.Repeat(" ", width-int(fullWidth)) +
			"]",
	)
}
