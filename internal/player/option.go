package player

import (
	tea "github.com/charmbracelet/bubbletea"
)

//go:generate enumer -type Option -linecomment -output option_string.go

type Option uint8

const (
	Option3xRewind  Option = iota // <<<
	Option2xRewind                // <<
	Option1xRewind                // <
	OptionPlayPause               // |>||
	Option1xForward               // >
	Option2xForward               // >>
	Option3xForward               // >>>
)

func chooseOption(o Option) tea.Cmd {
	return func() tea.Msg {
		return o
	}
}

func (o Option) DynamicString(isPlaying bool) string {
	switch o {
	case OptionPlayPause:
		if isPlaying {
			return o.String()[2:]
		} else {
			return o.String()[:2]
		}
	default:
		return o.String()
	}
}

func (o Option) Speed() float64 {
	switch o {
	case Option3xRewind:
		return -15
	case Option2xRewind:
		return -3
	case Option1xRewind:
		return -1
	case Option1xForward:
		return 1
	case Option2xForward:
		return 3
	case Option3xForward:
		return 15
	default:
		return 0
	}
}
