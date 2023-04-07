package frame

import "time"

type Frame struct {
	Num      int
	Duration time.Duration
	Height   int
	Data     string
}

func (f *Frame) CalcDuration(multiplier float64) time.Duration {
	return time.Duration(float64(f.Duration) / multiplier)
}
