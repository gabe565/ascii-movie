package telnet

import (
	"io"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func MakeRenderer(w io.Writer, profile termenv.Profile) *lipgloss.Renderer {
	r := lipgloss.NewRenderer(w, termenv.WithColorCache(true))
	r.SetColorProfile(profile)
	return r
}
