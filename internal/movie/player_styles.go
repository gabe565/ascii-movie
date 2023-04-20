package movie

import "github.com/charmbracelet/lipgloss"

var (
	appStyle = lipgloss.NewStyle().
			Margin(2, 4)

	screenStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.CompleteAdaptiveColor{
			Light: lipgloss.CompleteColor{TrueColor: "#CABDB3", ANSI256: "236", ANSI: "8"},
			Dark:  lipgloss.CompleteColor{TrueColor: "#312A24", ANSI256: "236", ANSI: "8"},
		})

	progressStyle = lipgloss.NewStyle().
			Margin(1, 0).
			Background(lipgloss.CompleteAdaptiveColor{
			Light: lipgloss.CompleteColor{TrueColor: "#eee", ANSI256: "255", ANSI: "15"},
			Dark:  lipgloss.CompleteColor{TrueColor: "#111", ANSI256: "233", ANSI: "0"},
		}).
		Foreground(lipgloss.Color("#626262")).
		Border(lipgloss.NormalBorder(), false, true).
		BorderForeground(lipgloss.Color("#594F46"))

	optionsStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(0, 1).
			Background(lipgloss.CompleteAdaptiveColor{
			Light: lipgloss.CompleteColor{TrueColor: "#eee", ANSI256: "254", ANSI: "15"},
			Dark:  lipgloss.CompleteColor{TrueColor: "#111", ANSI256: "235", ANSI: "0"},
		})

	activeStyle = optionsStyle.Copy().
			Background(lipgloss.CompleteAdaptiveColor{
			Light: lipgloss.CompleteColor{TrueColor: "#ccc", ANSI256: "246", ANSI: "7"},
			Dark:  lipgloss.CompleteColor{TrueColor: "#222", ANSI256: "237", ANSI: "8"},
		}).
		Foreground(lipgloss.CompleteAdaptiveColor{
			Light: lipgloss.CompleteColor{TrueColor: "0", ANSI256: "15", ANSI: "0"},
			Dark:  lipgloss.CompleteColor{TrueColor: "15", ANSI256: "15", ANSI: "0"},
		})

	selectedStyle = optionsStyle.Copy().
			Background(lipgloss.CompleteAdaptiveColor{
			Light: lipgloss.CompleteColor{TrueColor: "#55729D", ANSI256: "27", ANSI: "12"},
			Dark:  lipgloss.CompleteColor{TrueColor: "#2C3C55", ANSI256: "19", ANSI: "4"},
		}).
		Foreground(lipgloss.CompleteAdaptiveColor{
			Light: lipgloss.CompleteColor{TrueColor: "15", ANSI256: "15", ANSI: "15"},
			Dark:  lipgloss.CompleteColor{TrueColor: "15", ANSI256: "15", ANSI: "15"},
		})
)
