package player

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"gabe565.com/ascii-movie/internal/movie"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

func NewPlayer(m *movie.Movie, logger *slog.Logger, renderer *lipgloss.Renderer) *Player {
	if renderer == nil {
		renderer = lipgloss.DefaultRenderer()
	}

	playCtx, playCancel := context.WithCancel(context.Background())
	player := &Player{
		movie:          m,
		log:            logger,
		start:          time.Now(),
		zone:           zone.New(),
		speed:          1,
		selectedOption: OptionPlayPause,
		activeOption:   Option1xForward,
		styles:         NewStyles(m, renderer),
		playCtx:        playCtx,
		playCancel:     playCancel,
		keymap:         newKeymap(),
		help:           newHelp(renderer),
	}
	player.optionsCache = NewCache(player.OptionsView)
	player.helpCache = NewCache(player.HelpView)

	return player
}

type Player struct {
	movie *movie.Movie
	frame int
	log   *slog.Logger
	start time.Time
	zone  *zone.Manager

	speed      float64
	playCtx    context.Context
	playCancel context.CancelFunc

	selectedOption Option
	activeOption   Option
	styles         Styles
	optionsCache   *ViewCache

	keymap    keymap
	help      help.Model
	helpCache *ViewCache
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
		p.optionsCache.Invalidate()
		switch {
		case key.Matches(msg, p.keymap.quit):
			return p, tea.Quit
		case key.Matches(msg, p.keymap.left):
			if (p.selectedOption - 1).IsAOption() {
				p.selectedOption--
			}
		case key.Matches(msg, p.keymap.right):
			if (p.selectedOption + 1).IsAOption() {
				p.selectedOption++
			}
		case key.Matches(msg, p.keymap.choose):
			return p, chooseOption(p.selectedOption)
		case key.Matches(msg, p.keymap.home):
			p.selectedOption = 0
		case key.Matches(msg, p.keymap.end):
			opts := OptionValues()
			p.selectedOption = Option(len(opts) - 1) //nolint:gosec
		case key.Matches(msg, p.keymap.help):
			p.help.ShowAll = !p.help.ShowAll
			p.helpCache.Invalidate()
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
		case key.Matches(msg, p.keymap.stepPrev):
			if !p.isPlaying() && p.frame > 0 {
				p.frame--
				return p, nil
			}
		case key.Matches(msg, p.keymap.stepNext):
			if !p.isPlaying() && p.frame < len(p.movie.Frames)-1 {
				p.frame++
				return p, nil
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
	case Option:
		p.optionsCache.Invalidate()
		switch msg {
		case OptionPlayPause:
			if p.isPlaying() {
				p.pause()
				return p, nil
			}
			return p, p.play()
		default:
			p.activeOption = msg
			p.speed = msg.Speed()
		}
		p.pause()
		return p, p.play()
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

		for _, opt := range OptionValues() {
			if p.zone.Get(opt.String()).InBounds(msg) {
				p.selectedOption = opt
				return p, chooseOption(p.selectedOption)
			}
		}

		if p.zone.Get("help").InBounds(msg) {
			p.help.ShowAll = !p.help.ShowAll
			p.helpCache.Invalidate()
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
	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		p.styles.MarginX,
		lipgloss.JoinVertical(
			lipgloss.Center,
			p.styles.MarginY,
			p.styles.Screen.Render(p.movie.Frames[p.frame].Data),
			p.zone.Mark("progress", p.styles.Progress.Render(p.movie.Frames[p.frame].Progress)),
			p.optionsCache.String(),
			p.helpCache.String(),
		),
	)

	return p.zone.Scan(content)
}

func (p *Player) OptionsView() string {
	opts := OptionValues()
	options := make([]string, 0, len(opts))
	isPlaying := p.isPlaying()
	for _, option := range opts {
		var rendered string
		switch option {
		case p.selectedOption:
			rendered = p.styles.Selected.Render(option.DynamicString(isPlaying))
		case p.activeOption:
			rendered = p.styles.Active.Render(option.DynamicString(isPlaying))
		default:
			rendered = p.styles.Options.Render(option.DynamicString(isPlaying))
		}
		options = append(options, p.zone.Mark(option.String(), rendered))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, options...)
}

func (p *Player) HelpView() string {
	return p.zone.Mark("help", p.help.View(p.keymap))
}

func (p *Player) pause() {
	p.optionsCache.Invalidate()
	p.clearTimeouts()
}

func (p *Player) play() tea.Cmd {
	p.optionsCache.Invalidate()
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

func (p *Player) Close() {
	p.log = p.log.With("duration", time.Since(p.start).Truncate(100*time.Millisecond))
	if p.frame >= len(p.movie.Frames)-1 {
		p.log.Info("Finished movie")
	} else {
		p.log.Info("Disconnected early")
	}
	p.clearTimeouts()
	p.zone.Close()
}
