package player

import (
	"log/slog"

	"github.com/charmbracelet/lipgloss"
)

type Option func(*Player)

func WithLogger(l *slog.Logger) Option {
	return func(player *Player) {
		player.log = l
	}
}

func WithRenderer(r *lipgloss.Renderer) Option {
	return func(player *Player) {
		player.renderer = r
	}
}
