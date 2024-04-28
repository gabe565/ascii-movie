package loghooks

import (
	"time"

	"github.com/rs/zerolog"
)

func NewDuration() Duration {
	return Duration{
		start: time.Now(),
		trunc: time.Millisecond * 10,
	}
}

type Duration struct {
	start time.Time
	trunc time.Duration
}

func (d Duration) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Str("duration", time.Since(d.start).Truncate(d.trunc).String())
}
