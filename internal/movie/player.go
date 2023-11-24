package movie

import (
	"context"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabe565/ascii-movie/internal/log_hooks"
	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
)

func NewPlayer(m *Movie, logger *log.Entry, profile termenv.Profile) Player {
	player := Player{
		movie:           m,
		profile:         profile,
		speed:           1,
		selectedOption:  3,
		activeOption:    4,
		optionsStyle:    &optionsStyle,
		activeStyle:     &activeStyle,
		optionViewStale: true,
	}
	if profile == termenv.Ascii {
		player.optionsStyle = &optionsStyleAnsi
		player.activeStyle = &activeStyleAnsi
	}
	player.play()
	if logger != nil {
		player.durationHook = log_hooks.NewDuration()
		player.log = logger.WithField("duration", player.durationHook)
	}

	player.keymap = newKeymap()
	helpModel := help.New()
	player.helpViewCache = RenderHelpWithProfile(profile, helpModel, []key.Binding{
		player.keymap.quit,
		player.keymap.left,
		player.keymap.right,
		player.keymap.choose,
	})

	return player
}

type Player struct {
	movie        *Movie
	frame        int
	log          *log.Entry
	durationHook log_hooks.Duration
	profile      termenv.Profile

	speed      float64
	playCtx    context.Context
	playCancel context.CancelFunc

	selectedOption  int
	activeOption    int
	optionsStyle    *lipgloss.Style
	activeStyle     *lipgloss.Style
	optionViewCache string
	optionViewStale bool

	keymap        keymap
	helpViewCache string

	timeoutCtx    context.Context
	timeoutCancel context.CancelFunc
}

func (p Player) Init() tea.Cmd {
	return tick(p.playCtx, p.movie.Frames[p.frame].Duration, frameTickMsg{})
}

func (p Player) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case frameTickMsg:
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
			p.speed = 1
			p.activeOption = 4
			return p, p.pause()
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
				p.speed = 1
				p.activeOption = 4
				return p, p.pause()
			}
			p.frame += frameDiff
			duration += p.movie.Frames[p.frame].CalcDuration(speed)
		}
		return p, tick(p.playCtx, duration, frameTickMsg{})
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
		case key.Matches(msg, p.keymap.home):
			p.selectedOption = 0
		case key.Matches(msg, p.keymap.end):
			p.selectedOption = len(playerOptions) - 1
		case key.Matches(msg, p.keymap.jumps...):
			for i, binding := range p.keymap.jumps {
				if key.Matches(msg, binding) {
					p.frame = p.movie.Sections[i]
					if p.isPlaying() {
						return p, p.play()
					}
				}
			}
		}
	case quitMsg:
		if p.log != nil {
			p.log.Info("Disconnected early")
		}
		if p.playCancel != nil {
			p.playCancel()
		}
		if p.timeoutCancel != nil {
			p.timeoutCancel()
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
			if p.isPlaying() {
				return p, p.pause()
			} else {
				return p, p.play()
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
		if !p.isPlaying() {
			return p, p.play()
		}
	}
	return p, nil
}

func (p Player) View() string {
	if p.optionViewStale {
		p.optionViewCache = p.OptionsView()
	}

	return appStyle.RenderWithProfile(p.profile, lipgloss.JoinVertical(
		lipgloss.Center,
		p.movie.screenStyle.RenderWithProfile(p.profile, p.movie.Frames[p.frame].Data),
		progressStyle.RenderWithProfile(p.profile, p.movie.Frames[p.frame].Progress),
		p.optionViewCache,
		p.helpViewCache,
	))
}

func (p *Player) OptionsView() string {
	p.optionViewStale = false

	options := make([]string, 0, len(playerOptions))
	for i, option := range playerOptions {
		if option == OptionPause && !p.isPlaying() {
			option = OptionPlay
		}
		var rendered string
		if i == p.selectedOption {
			rendered = selectedStyle.RenderWithProfile(p.profile, string(option))
		} else if i == p.activeOption {
			rendered = p.activeStyle.RenderWithProfile(p.profile, string(option))
		} else {
			rendered = p.optionsStyle.RenderWithProfile(p.profile, string(option))
		}
		options = append(options, rendered)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, options...)
}

func (p *Player) pause() tea.Cmd {
	p.clearTimeouts()
	p.timeoutCtx, p.timeoutCancel = context.WithCancel(context.Background())
	return tick(p.timeoutCtx, 15*time.Minute, Quit())
}

func (p *Player) play() tea.Cmd {
	p.clearTimeouts()
	p.playCtx, p.playCancel = context.WithCancel(context.Background())
	return func() tea.Msg {
		return frameTickMsg{}
	}
}

func (p Player) isPlaying() bool {
	return p.playCtx.Err() == nil
}

func (p *Player) clearTimeouts() {
	if p.playCancel != nil {
		p.playCancel()
	}
	if p.timeoutCancel != nil {
		p.timeoutCancel()
	}
}
