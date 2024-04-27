package movie

import (
	"context"
	"time"

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
