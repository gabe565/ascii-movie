package movie

import "time"

type Frame struct {
	Duration time.Duration
	Data     string
	Progress string
}

func (f *Frame) CalcDuration(multiplier float64) time.Duration {
	return time.Duration(float64(f.Duration) / multiplier)
}
