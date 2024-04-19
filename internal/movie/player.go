package movie

import (
	"context"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabe565/ascii-movie/internal/loghooks"
	log "github.com/sirupsen/logrus"
)

func NewPlayer(m *Movie, logger *log.Entry, renderer *lipgloss.Renderer) Player {
	if renderer == nil {
		renderer = lipgloss.DefaultRenderer()
	}

	player := Player{
		movie:           m,
		renderer:        renderer,
		speed:           1,
		selectedOption:  3,
		activeOption:    4,
		styles:          NewStyles(m, renderer),
		optionViewStale: true,
	}
	player.play()
	if logger != nil {
		player.durationHook = loghooks.NewDuration()
		player.log = logger.WithField("duration", player.durationHook)
	}

	player.keymap = newKeymap()
	helpModel := help.New()
	helpModel.Styles.ShortSeparator = helpModel.Styles.ShortSeparator.Renderer(renderer)
	helpModel.Styles.ShortDesc = helpModel.Styles.ShortDesc.Renderer(renderer)
	helpModel.Styles.ShortKey = helpModel.Styles.ShortKey.Renderer(renderer)
	player.helpViewCache = helpModel.ShortHelpView([]key.Binding{
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
	durationHook loghooks.Duration
	renderer     *lipgloss.Renderer

	speed      float64
	playCtx    context.Context
	playCancel context.CancelFunc

	selectedOption  int
	activeOption    int
	styles          Styles
	optionViewCache string
	optionViewStale bool

	keymap        keymap
	helpViewCache string
}

func (p Player) Init() tea.Cmd {
	return tick(p.playCtx, p.movie.Frames[p.frame].Duration, frameTickMsg{})
}

//nolint:gocyclo
func (p Player) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case frameTickMsg:
		var frameDiff int
		switch {
		case p.speed >= 0:
			frameDiff = 1
			if p.frame+frameDiff >= len(p.movie.Frames) {
				if p.log != nil {
					p.log.Info("Finished movie")
				}
				return p, tea.Quit
			}
		case p.frame <= 0:
			p.speed = 1
			p.activeOption = 4
			p.pause()
			return p, nil
		default:
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
				p.pause()
				return p, nil
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
				p.selectedOption--
			}
		case key.Matches(msg, p.keymap.right):
			if p.selectedOption < len(playerOptions)-1 {
				p.selectedOption++
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
		p.clearTimeouts()
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
				p.pause()
				return p, nil
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
	case tea.WindowSizeMsg:
		p.styles.MarginX, p.styles.MarginY = "", ""
		if width := msg.Width/2 - p.movie.Width/2 - 1; width > 0 {
			p.styles.MarginX = strings.Repeat(" ", width)
		}
		if height := msg.Height/2 - lipgloss.Height(p.View())/2; height > 0 {
			p.styles.MarginY = strings.Repeat("\n", height)
		}
	}
	return p, nil
}

func (p Player) View() string {
	if p.optionViewStale {
		p.optionViewCache = p.OptionsView()
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		p.styles.MarginX,
		lipgloss.JoinVertical(
			lipgloss.Center,
			p.styles.MarginY,
			p.styles.Screen.Render(p.movie.Frames[p.frame].Data),
			p.styles.Progress.Render(p.movie.Frames[p.frame].Progress),
			p.optionViewCache,
			p.helpViewCache,
		),
	)

	return content
}

func (p *Player) OptionsView() string {
	p.optionViewStale = false

	options := make([]string, 0, len(playerOptions))
	for i, option := range playerOptions {
		if option == OptionPause && !p.isPlaying() {
			option = OptionPlay
		}
		var rendered string
		switch i {
		case p.selectedOption:
			rendered = p.styles.Selected.Render(string(option))
		case p.activeOption:
			rendered = p.styles.Active.Render(string(option))
		default:
			rendered = p.styles.Options.Render(string(option))
		}
		options = append(options, rendered)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, options...)
}

func (p *Player) pause() {
	p.clearTimeouts()
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
}
