package movie

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabe565/ascii-movie/internal/log_hooks"
	log "github.com/sirupsen/logrus"
)

var (
	appStyle = lipgloss.NewStyle().
			Margin(2, 4)

	screenStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#312A24"))

	progressStyle = lipgloss.NewStyle().
			Margin(1, 0).
			Background(lipgloss.Color("#111")).
			Foreground(lipgloss.Color("#626262")).
			Border(lipgloss.NormalBorder(), false, true).
			BorderForeground(lipgloss.Color("#594F46"))

	optionsStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Margin(0, 1).
			Background(lipgloss.Color("#111"))

	activeStyle = optionsStyle.Copy().
			Background(lipgloss.Color("#222"))

	selectedStyle = optionsStyle.Copy().
			Background(lipgloss.Color("#2C3C55"))
)

func NewPlayer(m *Movie, logger *log.Entry) Player {
	player := Player{
		movie:          m,
		speed:          1,
		selectedOption: 3,
		activeOption:   4,
		screenStyle:    screenStyle.Copy().Width(m.Width),
	}
	player.playCtx, player.pause = context.WithCancel(context.Background())
	if logger != nil {
		player.durationHook = log_hooks.NewDuration()
		player.log = logger.WithField("duration", player.durationHook)
	}

	player.keymap = newKeymap()
	player.help = help.New()
	player.helpKeys = []key.Binding{
		player.keymap.quit,
		player.keymap.left,
		player.keymap.right,
		player.keymap.choose,
	}

	return player
}

type Player struct {
	movie            *Movie
	frame            int
	log              *log.Entry
	durationHook     log_hooks.Duration
	LogExcludeFaster time.Duration

	speed   float64
	playCtx context.Context
	pause   context.CancelFunc

	selectedOption int
	activeOption   int

	screenStyle lipgloss.Style

	keymap   keymap
	help     help.Model
	helpKeys []key.Binding
}

func (p Player) Init() tea.Cmd {
	return tick(p.playCtx, p.movie.Frames[p.frame].Duration)
}

func (p Player) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, p.keymap.quit):
			return p, Quit
		case key.Matches(msg, p.keymap.left):
			if p.selectedOption > 0 {
				p.selectedOption -= 1
			}
		case key.Matches(msg, p.keymap.right):
			if p.selectedOption < 6 {
				p.selectedOption += 1
			}
		case key.Matches(msg, p.keymap.choose):
			return p, chooseOption(playerOptions[p.selectedOption])
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
		if p.speed >= 0 {
			p.frame += 1
		} else if p.frame == 0 {
			p.pause()
			p.speed = 1
			p.activeOption = 4
			return p, nil
		} else {
			p.frame -= 1
		}
		duration := p.movie.Frames[p.frame].Duration
		speed := p.speed
		if speed < 0 {
			speed *= -1
		}
		duration = time.Duration(float64(duration) / speed)
		return p, tick(p.playCtx, duration)
	case PlayerOption:
		switch msg {
		case Option3xRewind:
			p.activeOption = 0
			p.speed = -15
			return p, nil
		case Option2xRewind:
			p.activeOption = 1
			p.speed = -3
			return p, nil
		case Option1xRewind:
			p.activeOption = 2
			p.speed = -1
			return p, nil
		case OptionPause, OptionPlay:
			if p.playCtx.Err() == nil {
				p.pause()
				return p, nil
			} else {
				p.playCtx, p.pause = context.WithCancel(context.Background())
				return p, tick(p.playCtx, 0)
			}
		case Option1xForward:
			p.activeOption = 4
			p.speed = 1
			return p, nil
		case Option2xForward:
			p.activeOption = 5
			p.speed = 3
			return p, nil
		case Option3xForward:
			p.activeOption = 6
			p.speed = 15
			return p, nil
		}
	}
	return p, nil
}

func (p Player) View() string {
	options := make([]string, 0, len(playerOptions))
	for i, option := range playerOptions {
		if option == OptionPause && p.playCtx.Err() != nil {
			option = OptionPlay
		}
		var rendered string
		if i == p.selectedOption {
			rendered = selectedStyle.Render(string(option))
		} else if i == p.activeOption {
			rendered = activeStyle.Render(string(option))
		} else {
			rendered = optionsStyle.Render(string(option))
		}
		options = append(options, rendered)
	}
	optionsView := lipgloss.JoinHorizontal(lipgloss.Top, options...)

	shortHelp := p.help.ShortHelpView(p.helpKeys)

	return appStyle.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		p.screenStyle.Render(p.movie.Frames[p.frame].Data),
		progressStyle.Render(p.movie.Frames[p.frame].Progress),
		optionsView,
		shortHelp,
	))
}

func newKeymap() keymap {
	return keymap{
		quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c", "ctrl+d"),
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
