package player

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

func newHelp(renderer *lipgloss.Renderer) help.Model {
	keyStyle := renderer.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "246",
		Dark:  "242",
	})

	descStyle := renderer.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "249",
		Dark:  "239",
	})

	sepStyle := renderer.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: "253",
		Dark:  "237",
	})

	return help.Model{
		ShortSeparator: " • ",
		FullSeparator:  "    ",
		Ellipsis:       "…",
		Styles: help.Styles{
			ShortKey:       keyStyle,
			ShortDesc:      descStyle,
			ShortSeparator: sepStyle,
			Ellipsis:       sepStyle,
			FullKey:        keyStyle,
			FullDesc:       descStyle,
			FullSeparator:  sepStyle,
		},
	}
}
