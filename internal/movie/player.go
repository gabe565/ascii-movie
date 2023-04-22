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

func NewPlayer(m *Movie, logger *log.Entry) Player {
	player := Player{
		movie:           m,
		speed:           1,
		selectedOption:  3,
		activeOption:    4,
		optionViewStale: true,
	}
	player.playCtx, player.pause = context.WithCancel(context.Background())
	if logger != nil {
		player.durationHook = log_hooks.NewDuration()
		player.log = logger.WithField("duration", player.durationHook)
	}

	player.keymap = newKeymap()
	helpModel := help.New()
	player.helpViewCache = helpModel.ShortHelpView([]key.Binding{
		player.keymap.quit,
		player.keymap.left,
		player.keymap.right,
		player.keymap.choose,
	})

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

	selectedOption  int
	activeOption    int
	optionViewCache string
	optionViewStale bool

	keymap        keymap
	helpViewCache string
}

func (p Player) Init() tea.Cmd {
	return tick(p.playCtx, p.movie.Frames[p.frame].Duration)
}

func (p Player) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tickMsg:
		var frameDiff int
		if p.speed >= 0 {
			frameDiff = 1
			if p.frame+frameDiff >= len(p.movie.Frames) {
				if p.log != nil {
					p.log.Info("Finished movie")
				}
				return p, tea.Quit
			}
		} else if p.frame <= 0 {
			p.pause()
			p.speed = 1
			p.activeOption = 4
			return p, nil
		} else {
			frameDiff = -1
		}
		p.frame += frameDiff
		speed := p.speed
		if speed < 0 {
			speed *= -1
		}
		duration := p.movie.Frames[p.frame].CalcDuration(speed)
		for duration < time.Second/15 {
			if p.frame+frameDiff >= len(p.movie.Frames) {
				if p.log != nil {
					p.log.Info("Finished movie")
				}
				return p, tea.Quit
			} else if p.frame+frameDiff <= 0 {
				p.pause()
				p.speed = 1
				p.activeOption = 4
				return p, nil
			}
			p.frame += frameDiff
			duration += p.movie.Frames[p.frame].CalcDuration(speed)
		}
		return p, tick(p.playCtx, duration)
	case tea.KeyMsg:
		p.optionViewStale = true
		switch {
		case key.Matches(msg, p.keymap.quit):
			return p, Quit
		case key.Matches(msg, p.keymap.left):
			if p.selectedOption > 0 {
				p.selectedOption -= 1
			}
		case key.Matches(msg, p.keymap.right):
			if p.selectedOption < len(playerOptions)-1 {
				p.selectedOption += 1
			}
		case key.Matches(msg, p.keymap.choose):
			return p, chooseOption(playerOptions[p.selectedOption])
		}
	case quitMsg:
		p.pause()
		if p.log != nil {
			if time.Since(p.durationHook.GetStart()) < p.LogExcludeFaster {
				p.log.Trace("Disconnected early")
			} else {
				p.log.Info("Disconnected early")
			}
		}
		return p, tea.Quit
	case PlayerOption:
		p.optionViewStale = true
		switch msg {
		case Option3xRewind:
			p.activeOption = 0
			p.speed = -15
		case Option2xRewind:
			p.activeOption = 1
			p.speed = -3
		case Option1xRewind:
			p.activeOption = 2
			p.speed = -1
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
		case Option2xForward:
			p.activeOption = 5
			p.speed = 3
		case Option3xForward:
			p.activeOption = 6
			p.speed = 15
		}
		if p.playCtx.Err() != nil {
			p.playCtx, p.pause = context.WithCancel(context.Background())
			return p, tick(p.playCtx, 0)
		}
	}
	return p, nil
}

func (p Player) View() string {
	if p.optionViewStale {
		p.optionViewCache = p.OptionsView()
	}

	return appStyle.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		p.movie.screenStyle.Render(p.movie.Frames[p.frame].Data),
		progressStyle.Render(p.movie.Frames[p.frame].Progress),
		p.optionViewCache,
		p.helpViewCache,
	))
}

func (p Player) OptionsView() string {
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
	return lipgloss.JoinHorizontal(lipgloss.Top, options...)
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
