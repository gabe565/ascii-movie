package movie

import (
	"github.com/charmbracelet/lipgloss"
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

func NewStyles(m *Movie, renderer *lipgloss.Renderer) Styles {
	borderColor := lipgloss.AdaptiveColor{Light: "7", Dark: "8"}
	activeColor := lipgloss.AdaptiveColor{Light: "8", Dark: "12"}
	optionsColor := lipgloss.AdaptiveColor{Light: "7", Dark: "8"}
	selectedColor := lipgloss.AdaptiveColor{Light: "12", Dark: "4"}

	s := Styles{
		Screen: lipgloss.NewStyle().
			Renderer(renderer).
			Width(m.Width).
			Height(m.Height).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"}),

		Progress: lipgloss.NewStyle().
			Renderer(renderer).
			Margin(1, 0).
			Foreground(borderColor).
			Border(lipgloss.InnerHalfBlockBorder(), false, true).
			BorderForeground(borderColor),

		Options: lipgloss.NewStyle().
			Renderer(renderer).
			Padding(0, 2).
			MarginBottom(1).
			Border(lipgloss.InnerHalfBlockBorder()).
			BorderForeground(optionsColor).
			Background(optionsColor),
	}

	s.Active = s.Options.Copy().
		Background(activeColor).
		BorderForeground(activeColor).
		Foreground(lipgloss.AdaptiveColor{Light: "15"}).
		Bold(true)

	s.Selected = s.Options.Copy().
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
