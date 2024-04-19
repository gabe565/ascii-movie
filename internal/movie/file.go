package movie

import (
	"bufio"
	"bytes"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gabe565/ascii-movie/internal/progressbar"
)

func (m *Movie) LoadFile(path string, src io.Reader, speed float64) error {
	m.Filename = filepath.Base(path)

	frames := make([]Frame, 0, 2000)
	var f Frame
	var buf bytes.Buffer
	scanner := bufio.NewScanner(src)

	// Build part of every frame, excluding progress bar and bottom padding
	frameNum := -1
	frameHeadRe := regexp.MustCompile(`^\d+$`)
	for {
		ok := scanner.Scan()

		if frameHeadRe.Match(scanner.Bytes()) || !ok {
			frameNum++
			if frameNum != 0 {
				f.Data = strings.TrimSuffix(buf.String(), "\n")
				if frameHeight := strings.Count(f.Data, "\n"); m.Height < frameHeight {
					m.Height = frameHeight
				}
				buf.Reset()
				frames = append(frames, f)
			}
			if !ok {
				break
			}

			f = Frame{}

			v, err := strconv.Atoi(scanner.Text())
			if err != nil {
				return err
			}

			f.Duration = time.Duration(v) * time.Second / 15
			f.Duration = time.Duration(float64(f.Duration) / speed)
		} else {
			if len(scanner.Bytes()) > m.Width {
				m.Width = len(scanner.Bytes())
			}
			buf.WriteString(scanner.Text() + "\n")
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	m.Frames = make([]Frame, len(frames))
	copy(m.Frames, frames)

	// Compute the total duration
	var frameCap int
	bar := progressbar.New()
	totalDuration := m.Duration()

	// Build the rest of every frame and write to disk
	var currentPosition time.Duration
	for i, f := range m.Frames {
		f.Progress = bar.Generate(currentPosition+f.Duration/2, totalDuration, m.Width+2)
		m.Frames[i] = f
		percent := int(currentPosition * 10 / totalDuration)
		if percent < len(m.Sections)-1 {
			m.Sections[percent+1] = i
		}
		if frameCap < len(f.Data) {
			frameCap = len(f.Data)
		}
		currentPosition += f.Duration
	}

	m.Cap = frameCap
	return nil
}
