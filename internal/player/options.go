package player

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Option string

const (
	Option3xRewind  Option = "<<<"
	Option2xRewind  Option = "<<"
	Option1xRewind  Option = "<"
	OptionPause     Option = "||"
	OptionPlay      Option = "|>"
	Option1xForward Option = ">"
	Option2xForward Option = ">>"
	Option3xForward Option = ">>>"
)

var playerOptions = [...]Option{ //nolint:gochecknoglobals
	Option3xRewind,
	Option2xRewind,
	Option1xRewind,
	OptionPause,
	Option1xForward,
	Option2xForward,
	Option3xForward,
}

func chooseOption(o Option) tea.Cmd {
	return func() tea.Msg {
		return o
	}
}
