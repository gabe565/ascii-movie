package player

import (
	tea "github.com/charmbracelet/bubbletea"
)

//go:generate go tool enumer -type Button -linecomment -output button_string.go

type Button uint8

const (
	Button3xRewind  Button = iota // <<<
	Button2xRewind                // <<
	Button1xRewind                // <
	ButtonPlayPause               // |>||
	Button1xForward               // >
	Button2xForward               // >>
	Button3xForward               // >>>
)

func chooseButton(o Button) tea.Cmd {
	return func() tea.Msg {
		return o
	}
}

func (o Button) DynamicString(isPlaying bool) string {
	switch o {
	case ButtonPlayPause:
		if isPlaying {
			return o.String()[2:]
		}
		return o.String()[:2]
	default:
		return o.String()
	}
}

func (o Button) Speed() float64 {
	switch o {
	case Button3xRewind:
		return -15
	case Button2xRewind:
		return -3
	case Button1xRewind:
		return -1
	case Button1xForward:
		return 1
	case Button2xForward:
		return 3
	case Button3xForward:
		return 15
	default:
		return 0
	}
}
