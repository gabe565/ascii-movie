package movie

import (
	"context"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type frameTickMsg time.Time

func tick(ctx context.Context, d time.Duration, msg tea.Msg) tea.Cmd {
	return func() tea.Msg {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(d):
			return msg
		}
	}
}

type quitMsg struct{}

func Quit() tea.Msg {
	return quitMsg{}
}

type keymap struct {
	quit     key.Binding
	left     key.Binding
	right    key.Binding
	navigate key.Binding
	home     key.Binding
	end      key.Binding
	choose   key.Binding
	jumps    []key.Binding
}

func newKeymap() keymap {
	jumps := make([]key.Binding, 0, 10)
	for i := range 10 {
		jumps = append(jumps, key.NewBinding(
			key.WithKeys(strconv.Itoa(i)),
		))
	}

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
		navigate: key.NewBinding(
			key.WithKeys("left", "h", "right", "l"),
			key.WithHelp("←→", "navigate"),
		),
		home: key.NewBinding(
			key.WithKeys("home"),
		),
		end: key.NewBinding(
			key.WithKeys("end"),
		),
		choose: key.NewBinding(
			key.WithKeys(" ", "enter"),
			key.WithHelp("enter", "select"),
		),
		jumps: jumps,
	}
}
