package movie

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle().
			Margin(2, 4)

	borderColor = lipgloss.AdaptiveColor{Light: "7", Dark: "8"}

	screenStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Foreground(lipgloss.AdaptiveColor{Light: "0", Dark: "15"})

	progressStyle = lipgloss.NewStyle().
			Margin(1, 0).
			Foreground(borderColor).
			Border(lipgloss.NormalBorder(), false, true).
			BorderForeground(borderColor)

	optionsColor = lipgloss.AdaptiveColor{Light: "7", Dark: "8"}
	optionsStyle = lipgloss.NewStyle().
			Padding(0, 2).
			MarginBottom(1).
			Border(lipgloss.InnerHalfBlockBorder()).
			BorderForeground(optionsColor).
			Background(optionsColor)

	activeColor = lipgloss.AdaptiveColor{Light: "8", Dark: "12"}
	activeStyle = optionsStyle.Copy().
			Background(activeColor).
			BorderForeground(activeColor).
			Foreground(lipgloss.AdaptiveColor{Light: "15"}).
			Bold(true)

	selectedColor = lipgloss.AdaptiveColor{Light: "12", Dark: "4"}
	selectedStyle = optionsStyle.Copy().
			Background(selectedColor).
			BorderForeground(selectedColor).
			Foreground(lipgloss.Color("15")).
			Bold(true)
)
