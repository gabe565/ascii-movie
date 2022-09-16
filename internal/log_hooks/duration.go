package log_hooks

import (
	"time"
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

func (d Duration) String() string {
	return time.Since(d.start).Truncate(d.trunc).String()
}
