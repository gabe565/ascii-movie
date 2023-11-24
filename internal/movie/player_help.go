package movie

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func RenderHelpWithProfile(p termenv.Profile, m help.Model, bindings []key.Binding) string {
	if len(bindings) == 0 {
		return ""
	}

	var b strings.Builder
	var totalWidth int
	separator := m.Styles.ShortSeparator.Inline(true).RenderWithProfile(p, m.ShortSeparator)

	for i, kb := range bindings {
		if !kb.Enabled() {
			continue
		}

		var sep string
		if totalWidth > 0 && i < len(bindings) {
			sep = separator
		}

		str := sep +
			m.Styles.ShortKey.Inline(true).RenderWithProfile(p, kb.Help().Key) + " " +
			m.Styles.ShortDesc.Inline(true).RenderWithProfile(p, kb.Help().Desc)

		w := lipgloss.Width(str)

		// If adding this help item would go over the available width, stop
		// drawing.
		if m.Width > 0 && totalWidth+w > m.Width {
			// Although if there's room for an ellipsis, print that.
			tail := " " + m.Styles.Ellipsis.Inline(true).RenderWithProfile(p, m.Ellipsis)
			tailWidth := lipgloss.Width(tail)

			if totalWidth+tailWidth < m.Width {
				b.WriteString(tail)
			}

			break
		}

		totalWidth += w
		b.WriteString(str)
	}

	return b.String()
}
