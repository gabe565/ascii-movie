package movie

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

func NewMovie() Movie {
	return Movie{}
}

type Movie struct {
	Filename string
	Cap      int
	Frames   []Frame
	Width    int

	screenStyle lipgloss.Style
}

func (m Movie) Duration() time.Duration {
	var totalDuration time.Duration
	for _, f := range m.Frames {
		totalDuration += f.Duration
	}
	return totalDuration
}
