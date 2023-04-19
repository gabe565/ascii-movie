package movie

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

func NewMovie() Movie {
	return Movie{
		BodyStyle:     lipgloss.NewStyle(),
		ProgressStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")),
	}
}

type Movie struct {
	Filename string
	Cap      int
	Frames   []Frame

	BodyStyle     lipgloss.Style
	ProgressStyle lipgloss.Style
}

func (m Movie) Duration() time.Duration {
	var totalDuration time.Duration
	for _, f := range m.Frames {
		totalDuration += f.Duration
	}
	return totalDuration
}
