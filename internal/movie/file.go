package movie

import (
	"bufio"
	"github.com/gabe565/ascii-movie/internal/progressbar"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func NewFromFile(path string, src io.Reader, pad Padding, progressPad Padding) (Movie, error) {
	m := Movie{
		Filename: filepath.Base(path),
		Speed:    1,
	}
	var f Frame
	var maxWidth int
	scanner := bufio.NewScanner(src)

	// Build part of every frame, excluding progress bar and bottom padding
	frameNum := -1
	frameHeadRe := regexp.MustCompile(`^\d+$`)
	for scanner.Scan() {
		if frameHeadRe.Match(scanner.Bytes()) {
			frameNum += 1
			if frameNum != 0 {
				m.Frames = append(m.Frames, f)
			}

			f = Frame{
				Data: strings.Repeat("\n", pad.Top),
			}

			v, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return m, err
			}

			f.Duration = time.Duration(float64(v)*(1000.0/15.0)) * time.Millisecond
		} else {
			if len(scanner.Bytes()) > maxWidth {
				maxWidth = len(scanner.Bytes())
			}
			f.Data += strings.Repeat(" ", pad.Left) + scanner.Text() + "\n"
		}
	}
	m.Frames = append(m.Frames, f)
	if err := scanner.Err(); err != nil {
		return m, err
	}

	// Compute the total duration
	var frameCap int
	bar := progressbar.New()
	totalDuration := m.Duration()

	// Build the rest of every frame and write to disk
	var currentPosition time.Duration
	for i, f := range m.Frames {
		f.Data += strings.Repeat("\n", pad.Bottom)
		if pad.Left != 0 {
			f.Data += strings.Repeat(" ", pad.Left)
		}
		f.Data += bar.Generate(currentPosition+f.Duration/2, totalDuration, maxWidth) + "\n"
		f.Data += strings.Repeat("\n", progressPad.Bottom)
		f.Height = strings.Count(f.Data, "\n")
		m.Frames[i] = f
		if frameCap < len(f.Data) {
			frameCap = len(f.Data)
		}
		currentPosition += f.Duration
	}

	m.Cap = frameCap

	return m, nil
}
