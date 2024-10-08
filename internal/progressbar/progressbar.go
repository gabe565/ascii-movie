package progressbar

import (
	"strings"
	"time"
)

type ProgressBar struct {
	Phases []string
}

func DefaultPhases() []string {
	return []string{
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
}

func New() ProgressBar {
	return ProgressBar{
		Phases: DefaultPhases(),
	}
}

func (p ProgressBar) Generate(n, total time.Duration, width int) string {
	width -= 2

	percent := float64(n) / float64(total)
	filledLen := percent * float64(width)
	filledNum := int(filledLen)
	phaseIdx := int((filledLen - float64(filledNum)) * float64(len(p.Phases)))
	emptyNum := width - filledNum

	var buf strings.Builder
	buf.Grow(width * 3)

	buf.WriteString(strings.Repeat(p.Phases[len(p.Phases)-1], filledNum))
	if phaseIdx > 0 {
		buf.WriteString(p.Phases[phaseIdx])
		emptyNum--
	}
	buf.WriteString(strings.Repeat(p.Phases[0], emptyNum))
	return buf.String()
}
