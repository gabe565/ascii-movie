package movie

import (
	tea "github.com/charmbracelet/bubbletea"
)

type PlayerOption string

const (
	Option3xRewind  PlayerOption = "<<<"
	Option2xRewind  PlayerOption = "<<"
	Option1xRewind  PlayerOption = "<"
	OptionPause     PlayerOption = "||"
	OptionPlay      PlayerOption = "|>"
	Option1xForward PlayerOption = ">"
	Option2xForward PlayerOption = ">>"
	Option3xForward PlayerOption = ">>>"
)

var playerOptions = [...]PlayerOption{
	Option3xRewind,
	Option2xRewind,
	Option1xRewind,
	OptionPause,
	Option1xForward,
	Option2xForward,
	Option3xForward,
}

func chooseOption(o PlayerOption) tea.Cmd {
	return func() tea.Msg {
		return o
	}
}
