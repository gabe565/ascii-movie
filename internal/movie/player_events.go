package movie

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

func tick(ctx context.Context, d time.Duration) tea.Cmd {
	return func() tea.Msg {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(d):
			return tickMsg{}
		}
	}
}

type quitMsg struct{}

func Quit() tea.Msg {
	return quitMsg{}
}

type keymap struct {
	quit   key.Binding
	left   key.Binding
	right  key.Binding
	choose key.Binding
}

func newKeymap() keymap {
	return keymap{
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c", "ctrl+d", "esc"),
			key.WithHelp("q", "quit"),
		),
		left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "left"),
		),
		right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "right"),
		),
		choose: key.NewBinding(
			key.WithKeys(" ", "enter"),
			key.WithHelp("enter", "choose"),
		),
	}
}
