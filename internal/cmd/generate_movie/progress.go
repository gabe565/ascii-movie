package main

import (
	"github.com/fatih/color"
	_ "github.com/fatih/color"
	"strings"
	"time"
)

var colorLightBlack = color.New(38).SprintFunc()

var phases = [...]string{
	" ",
	"▏",
	"▎",
	"▍",
	"▌",
	"▋",
	"▊",
	"▉",
	"█",
}

func progressBar(n, total time.Duration, width int) string {
	percent := float64(n) / float64(total)
	filledLen := percent * float64(width)
	filledNum := int(filledLen)
	phaseIdx := int((filledLen - float64(filledNum)) * float64(len(phases)))
	emptyNum := width - filledNum

	result := "["
	result += strings.Repeat(phases[len(phases)-1], filledNum)
	if phaseIdx > 0 {
		result += phases[phaseIdx]
		emptyNum -= 1
	}
	result += strings.Repeat(phases[0], emptyNum)
	result += "]"
	return colorLightBlack(result)
}
