package progressbar

import (
	"github.com/fatih/color"
	"strings"
	"time"
)

type ProgressBar struct {
	Formatter func(...any) string
	Phases    []string
}

func New() ProgressBar {
	return ProgressBar{
		Formatter: color.New(38).SprintFunc(),
		Phases: []string{
			" ",
			"▏",
			"▎",
			"▍",
			"▌",
			"▋",
			"▊",
			"▉",
			"█",
		},
	}
}

func (p *ProgressBar) Generate(n, total time.Duration, width int) string {
	width -= 2
	percent := float64(n) / float64(total)
	filledLen := percent * float64(width)
	filledNum := int(filledLen)
	phaseIdx := int((filledLen - float64(filledNum)) * float64(len(p.Phases)))
	emptyNum := width - filledNum

	result := "["
	result += strings.Repeat(p.Phases[len(p.Phases)-1], filledNum)
	if phaseIdx > 0 {
		result += p.Phases[phaseIdx]
		emptyNum -= 1
	}
	result += strings.Repeat(p.Phases[0], emptyNum)
	result += "]"
	return p.Formatter(result)
}
