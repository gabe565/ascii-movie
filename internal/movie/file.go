package movie

import (
	"bufio"
	"bytes"
	"io"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gabe565/ascii-movie/internal/progressbar"
)

func (m *Movie) LoadFile(path string, src io.Reader, speed float64) error {
	m.Filename = filepath.Base(path)

	m.Frames = make([]Frame, 0, 2000)
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
				if frameWidth := lipgloss.Width(f.Data); frameWidth > m.Width {
					m.Width = frameWidth
				}
				if frameHeight := lipgloss.Height(f.Data); frameHeight > m.Height {
					m.Height = frameHeight
				}
				buf.Reset()
				m.Frames = append(m.Frames, f)
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
			buf.WriteString(scanner.Text() + "\n")
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	m.Frames = slices.Clip(m.Frames)

	// Compute the total duration
	bar := progressbar.New()
	totalDuration := m.Duration()

	// Build the rest of every frame and write to disk
	var currentPosition time.Duration
	m.Sections = make([]int, m.Width+1)
	for i, f := range m.Frames {
		f.Progress = bar.Generate(currentPosition+f.Duration/2, totalDuration, m.Width+2)
		m.Frames[i] = f
		percent := int(currentPosition * time.Duration(m.Width) / totalDuration)
		if percent < len(m.Sections)-1 {
			m.Sections[percent+1] = i
		}
		currentPosition += f.Duration
	}

	return nil
}
