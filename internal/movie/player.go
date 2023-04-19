package movie

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
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

	player.keymap = newKeymap()
	player.help = help.New()

	return player
}

type Player struct {
	movie            *Movie
	frame            int
	log              *log.Entry
	durationHook     log_hooks.Duration
	LogExcludeFaster time.Duration

	keymap keymap
	help   help.Model
}

func (p Player) Init() tea.Cmd {
	return tick(p.movie.Frames[p.frame].Duration)
}

func (p Player) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, p.keymap.quit):
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
	shortHelp := p.help.ShortHelpView([]key.Binding{
		p.keymap.quit,
	})

	return p.movie.BodyStyle.Render(lipgloss.JoinVertical(
		lipgloss.Left,
		p.movie.Frames[p.frame].Data,
		p.movie.ProgressStyle.Render(p.movie.Frames[p.frame].Progress),
		shortHelp,
	))
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

type keymap struct {
	quit key.Binding
}

func newKeymap() keymap {
	return keymap{
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c", "ctrl+d"),
			key.WithHelp("q", "quit"),
		),
	}
}
