package player

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/muesli/termenv"
)

type Styles struct {
	Screen   lipgloss.Style
	Progress lipgloss.Style
	Options  lipgloss.Style
	Active   lipgloss.Style
	Selected lipgloss.Style

	MarginX, MarginY string
}

func NewStyles(m *movie.Movie, renderer *lipgloss.Renderer) Styles {
	borderColor := lipgloss.AdaptiveColor{Light: "7", Dark: "8"}
	activeColor := lipgloss.AdaptiveColor{Light: "8", Dark: "12"}
	optionsColor := lipgloss.AdaptiveColor{Light: "7", Dark: "8"}
	selectedColor := lipgloss.AdaptiveColor{Light: "12", Dark: "4"}

	isTTY := renderer.Output().TTY() != nil
	screenStyle := renderer.NewStyle().
		Width(m.Width).
		Height(m.Height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)
	if isTTY {
		screenStyle = screenStyle.Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"})
	}

	s := Styles{
		Screen: screenStyle,

		Progress: renderer.NewStyle().
			Margin(1, 0).
			Foreground(borderColor).
			Border(lipgloss.InnerHalfBlockBorder(), false, true).
			BorderForeground(borderColor),

		Options: renderer.NewStyle().
			Padding(0, 2).
			MarginBottom(1).
			Border(lipgloss.InnerHalfBlockBorder()).
			BorderForeground(optionsColor).
			Background(optionsColor).
			Foreground(lipgloss.AdaptiveColor{Light: "15", Dark: "7"}),
	}

	s.Active = s.Options.
		Background(activeColor).
		BorderForeground(activeColor).
		Bold(true)

	s.Selected = s.Options.
		Background(selectedColor).
		BorderForeground(selectedColor).
		Foreground(lipgloss.Color("15")).
		Bold(true)

	if renderer.ColorProfile() == termenv.Ascii {
		s.Options = s.Options.
			Padding(0, 2).
			Margin(1).
			Border(lipgloss.InnerHalfBlockBorder(), false)

		s.Active = s.Active.BorderStyle(lipgloss.DoubleBorder())
	}

	return s
}
