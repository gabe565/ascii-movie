package movie

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabe565/ascii-movie/internal/log_hooks"
	log "github.com/sirupsen/logrus"
)

func NewPlayer(m *Movie, logger *log.Entry) Player {
	player := Player{movie: m}
	if logger != nil {
		player.durationHook = log_hooks.NewDuration()
		player.log = logger.WithField("duration", player.durationHook)
	}
	return player
}

type Player struct {
	movie            *Movie
	frame            int
	log              *log.Entry
	durationHook     log_hooks.Duration
	LogExcludeFaster time.Duration
}

func (p Player) Init() tea.Cmd {
	return tick(p.movie.Frames[p.frame].Duration)
}

func (p Player) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "ctrl+d":
			return p, Quit
		}
	case quitMsg:
		if p.log != nil {
			if time.Since(p.durationHook.GetStart()) < p.LogExcludeFaster {
				p.log.Trace("Disconnected early")
			} else {
				p.log.Info("Disconnected early")
			}
		}
		return p, tea.Quit
	case tickMsg:
		if p.frame+1 >= len(p.movie.Frames) {
			if p.log != nil {
				p.log.Info("Finished movie")
			}
			return p, tea.Quit
		}
		p.frame += 1
		return p, tick(p.movie.Frames[p.frame].Duration)
	}
	return p, nil
}

func (p Player) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		p.movie.BodyStyle.Render(p.movie.Frames[p.frame].Data),
		p.movie.ProgressStyle.Render(p.movie.Frames[p.frame].Progress),
	)
}

type tickMsg time.Time

func tick(d time.Duration) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(d)
		return tickMsg{}
	}
}

type quitMsg struct{}

func Quit() tea.Msg {
	return quitMsg{}
}
