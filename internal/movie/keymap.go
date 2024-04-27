package movie

import (
	"strconv"

	"github.com/charmbracelet/bubbles/key"
)

func newKeymap() keymap {
	jumps := make([]key.Binding, 0, 10)
	for i := range 10 {
		jumps = append(jumps, key.NewBinding(
			key.WithKeys(strconv.Itoa(i)),
		))
	}

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
		navigate: key.NewBinding(
			key.WithKeys("left", "h", "right", "l"),
			key.WithHelp("←→", "navigate"),
		),
		home: key.NewBinding(
			key.WithKeys("home"),
			key.WithHelp("home", "go to start"),
		),
		end: key.NewBinding(
			key.WithKeys("end"),
			key.WithHelp("end", "go to end"),
		),
		choose: key.NewBinding(
			key.WithKeys(" ", "enter"),
			key.WithHelp("enter", "select"),
		),
		chooseFull: key.NewBinding(
			key.WithKeys(" ", "enter"),
			key.WithHelp("enter/space", "select"),
		),
		jumps: jumps,
		jump: key.NewBinding(
			key.WithKeys("0", "1", "2", "3", "4", "5", "6", "7", "8", "9"),
			key.WithHelp("0-9", "jump to position"),
		),
		help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
	}
}

type keymap struct {
	quit       key.Binding
	left       key.Binding
	right      key.Binding
	navigate   key.Binding
	home       key.Binding
	end        key.Binding
	choose     key.Binding
	chooseFull key.Binding
	jumps      []key.Binding
	jump       key.Binding
	help       key.Binding
}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.navigate,
		k.choose,
		k.quit,
		k.help,
	}
}

func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{
			k.left,
			k.right,
			k.chooseFull,
			k.home,
			k.end,
		},
		{
			k.jump,
			k.help,
			k.quit,
		},
	}
}
