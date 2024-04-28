package player

import (
	"context"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gabe565/ascii-movie/internal/loghooks"
	"github.com/gabe565/ascii-movie/internal/movie"
	zone "github.com/lrstanley/bubblezone"
	"github.com/rs/zerolog"
)

func NewPlayer(m *movie.Movie, logger zerolog.Logger, renderer *lipgloss.Renderer) *Player {
	if renderer == nil {
		renderer = lipgloss.DefaultRenderer()
	}

	playCtx, playCancel := context.WithCancel(context.Background())
	player := &Player{
		movie:           m,
		log:             logger.Hook(loghooks.NewDuration()),
		renderer:        renderer,
		zone:            zone.New(),
		speed:           1,
		selectedOption:  3,
		activeOption:    4,
		styles:          NewStyles(m, renderer),
		optionViewStale: true,
		playCtx:         playCtx,
		playCancel:      playCancel,
		keymap:          newKeymap(),
		help:            newHelp(renderer),
		helpViewStale:   true,
	}

	return player
}

type Player struct {
	movie    *movie.Movie
	frame    int
	log      zerolog.Logger
	renderer *lipgloss.Renderer
	zone     *zone.Manager

	speed      float64
	playCtx    context.Context
	playCancel context.CancelFunc

	selectedOption  int
	activeOption    int
	styles          Styles
	optionViewCache string
	optionViewStale bool

	keymap        keymap
	help          help.Model
	helpViewStale bool
	helpViewCache string
}

func (p *Player) Init() tea.Cmd {
	return tick(p.playCtx, p.movie.Frames[p.frame].Duration, frameTickMsg{})
}

//nolint:gocyclo
func (p *Player) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case frameTickMsg:
		var frameDiff int
		switch {
		case p.speed >= 0:
			frameDiff = 1
			if p.frame+frameDiff >= len(p.movie.Frames) {
				p.log.Info().Msg("Finished movie")
				return p, Quit
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
				p.log.Info().Msg("Finished movie")
				return p, Quit
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
		case key.Matches(msg, p.keymap.help):
			p.help.ShowAll = !p.help.ShowAll
			p.helpViewStale = true
		case key.Matches(msg, p.keymap.jumpPrev):
			var d time.Duration
			for d < 5*time.Second && p.frame > 0 {
				p.frame--
				d += p.movie.Frames[p.frame].Duration
			}
			if p.isPlaying() {
				return p, p.play()
			}
		case key.Matches(msg, p.keymap.jumpNext):
			var d time.Duration
			for d < 5*time.Second && p.frame < len(p.movie.Frames)-1 {
				p.frame++
				d += p.movie.Frames[p.frame].Duration
			}
			if p.isPlaying() {
				return p, p.play()
			}
		case key.Matches(msg, p.keymap.jumps...):
			for i, binding := range p.keymap.jumps {
				if key.Matches(msg, binding) {
					p.frame = p.movie.Sections[i*len(p.movie.Sections)/10]
					if p.isPlaying() {
						return p, p.play()
					}
				}
			}
		}
	case quitMsg:
		p.log.Info().Msg("Disconnected early")
		p.clearTimeouts()
		p.zone.Close()
		return p, tea.Quit //nolint:forbidigo
	case Option:
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
		return p, nil
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionRelease || msg.Button != tea.MouseButtonLeft {
			return p, nil
		}

		if p.zone.Get("progress").InBounds(msg) {
			x, _ := p.zone.Get("progress").Pos(msg)
			x--
			if x < 0 {
				x = 0
			}
			p.frame = p.movie.Sections[x]
			if p.isPlaying() {
				return p, p.play()
			}
		}

		if msg.Action != tea.MouseActionPress {
			return p, nil
		}

		switch {
		case p.zone.Get(string(Option3xRewind)).InBounds(msg):
			p.selectedOption = 0
			return p, chooseOption(Option3xRewind)
		case p.zone.Get(string(Option2xRewind)).InBounds(msg):
			p.selectedOption = 1
			return p, chooseOption(Option2xRewind)
		case p.zone.Get(string(Option1xRewind)).InBounds(msg):
			p.selectedOption = 2
			return p, chooseOption(Option1xRewind)
		case p.zone.Get(string(OptionPause)).InBounds(msg):
			p.selectedOption = 3
			return p, chooseOption(OptionPause)
		case p.zone.Get(string(OptionPlay)).InBounds(msg):
			p.selectedOption = 3
			return p, chooseOption(OptionPlay)
		case p.zone.Get(string(Option1xForward)).InBounds(msg):
			p.selectedOption = 4
			return p, chooseOption(Option1xForward)
		case p.zone.Get(string(Option2xForward)).InBounds(msg):
			p.selectedOption = 5
			return p, chooseOption(Option2xForward)
		case p.zone.Get(string(Option3xForward)).InBounds(msg):
			p.selectedOption = 6
			return p, chooseOption(Option3xForward)
		case p.zone.Get("help").InBounds(msg):
			p.help.ShowAll = !p.help.ShowAll
			p.helpViewStale = true
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

func (p *Player) View() string {
	if p.optionViewStale {
		p.optionViewCache = p.OptionsView()
	}
	if p.helpViewStale {
		p.helpViewCache = p.HelpView()
	}

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		p.styles.MarginX,
		lipgloss.JoinVertical(
			lipgloss.Center,
			p.styles.MarginY,
			p.styles.Screen.Render(p.movie.Frames[p.frame].Data),
			p.zone.Mark("progress", p.styles.Progress.Render(p.movie.Frames[p.frame].Progress)),
			p.optionViewCache,
			p.helpViewCache,
		),
	)

	return p.zone.Scan(content)
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
		options = append(options, p.zone.Mark(string(option), rendered))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, options...)
}

func (p *Player) HelpView() string {
	p.helpViewStale = false
	v := p.help.View(p.keymap)
	if p.help.ShowAll {
		sep := p.help.Styles.FullSeparator.Render(p.help.FullSeparator)
		sepSpaces := strings.Repeat(" ", lipgloss.Width(sep))
		// Remove first line separator
		v = strings.Replace(v, sep+"\n", "\n", 1)
		// Remove separator spaces form other lines
		v = strings.ReplaceAll(v, sepSpaces+"\n", "\n")
		// Remove separator spaces from final line
		v = strings.TrimSuffix(v, sepSpaces)
	}
	return p.zone.Mark("help", v)
}

func (p *Player) pause() {
	p.optionViewStale = true
	p.clearTimeouts()
}

func (p *Player) play() tea.Cmd {
	p.optionViewStale = true
	p.clearTimeouts()
	p.playCtx, p.playCancel = context.WithCancel(context.Background())
	return func() tea.Msg {
		return frameTickMsg{}
	}
}

func (p *Player) isPlaying() bool {
	return p.playCtx.Err() == nil
}

func (p *Player) clearTimeouts() {
	if p.playCancel != nil {
		p.playCancel()
	}
}
