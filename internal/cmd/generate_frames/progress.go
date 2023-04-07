package main

import (
	"github.com/fatih/color"
	_ "github.com/fatih/color"
	"math"
	"strings"
	"time"
)

var colorLightBlack = color.New(38).SprintFunc()

var parts = []rune{
	' ',
	'▏',
	'▎',
	'▍',
	'▌',
	'▋',
	'▊',
	'▉',
	'█',
}

func progressBar(n, total time.Duration, width int) string {
	percent := float64(n) / float64(total)
	fullWidth := percent * float64(width)
	part := parts[int(math.Round(math.Mod(fullWidth, 1.0)*8))]
	return colorLightBlack(
		"[" +
			strings.Repeat("█", int(fullWidth)) +
			string(part) +
			strings.Repeat(" ", width-int(fullWidth)) +
			"]",
	)
}
