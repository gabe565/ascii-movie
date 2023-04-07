package movie

import (
	"time"
)

type Movie struct {
	Filename string
	Cap      int
	Frames   []Frame
}

func (m Movie) Duration(speed float64) time.Duration {
	var totalDuration time.Duration
	for _, f := range m.Frames {
		totalDuration += f.CalcDuration(speed)
	}
	return totalDuration
}
