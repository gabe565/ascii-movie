package frame

import "time"

type Frame struct {
	Num    int
	Sleep  time.Duration
	Height int
	Data   string
}
