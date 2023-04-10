package movie

import (
	"bufio"
	"github.com/gabe565/ascii-movie/internal/progressbar"
	"io"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func NewFromFile(path string, src io.Reader, frameHeight int, pad Padding, progressPad Padding) (*Movie, error) {
	m := Movie{
		Filename: filepath.Base(path),
		Speed:    1,
	}
	var f Frame
	var maxWidth int
	scanner := bufio.NewScanner(src)

	// Build part of every frame, excluding progress bar and bottom padding
	for lineNum := 0; scanner.Scan(); lineNum += 1 {
		frameLineNum := lineNum % frameHeight
		if frameLineNum == 0 {
			f = Frame{
				Num:  lineNum / frameHeight,
				Data: strings.Repeat("\n", pad.Top-1),
			}

			v, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return nil, err
			}

			f.Duration = time.Duration(float64(v)*(1000.0/15.0)) * time.Millisecond
		} else {
			if len(scanner.Bytes()) > maxWidth {
				maxWidth = len(scanner.Bytes())
			}
			f.Data += "\n" + strings.Repeat(" ", pad.Left) + scanner.Text()
		}

		if frameLineNum == frameHeight-1 {
			m.Frames = append(m.Frames, f)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Compute the total duration
	var frameCap int
	bar := progressbar.New()
	totalDuration := m.Duration()

	// Build the rest of every frame and write to disk
	var currentPosition time.Duration
	for i, f := range m.Frames {
		f.Data += strings.Repeat("\n", pad.Bottom)
		f.Data += strings.Repeat(" ", pad.Left-1)
		f.Data += bar.Generate(currentPosition+f.Duration/2, totalDuration, maxWidth)
		f.Data += strings.Repeat("\n", progressPad.Bottom)
		f.Height = strings.Count(f.Data, "\n")
		m.Frames[i] = f
		if frameCap < len(f.Data) {
			frameCap = len(f.Data)
		}
		currentPosition += f.Duration
	}

	m.Cap = frameCap

	return &m, nil
}
