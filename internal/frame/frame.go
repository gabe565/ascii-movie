package frame

import "time"

type Frame struct {
	Num    int
	Sleep  time.Duration
	Height int
	Data   string
}

func (f *Frame) ComputeSleep(multiplier float64) time.Duration {
	return time.Duration(float64(f.Sleep) / multiplier)
}
