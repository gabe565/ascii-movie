package util

import (
	"strings"

	"github.com/muesli/termenv"
)

func Profile(term string) termenv.Profile {
	term = strings.ToLower(term)
	switch {
	case strings.Contains(term, "256color"):
		return termenv.ANSI256
	case term == "",
		strings.Contains(term, "color"),
		strings.Contains(term, "xterm"),
		strings.Contains(term, "ansi"),
		strings.Contains(term, "tmux"),
		strings.Contains(term, "screen"),
		strings.Contains(term, "cygwin"),
		strings.Contains(term, "rxvt"):
		return termenv.ANSI
	default:
		return termenv.Ascii
	}
}
