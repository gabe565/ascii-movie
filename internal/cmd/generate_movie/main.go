package main

import (
	"bufio"
	_ "embed"
	"github.com/gabe565/ascii-movie/config"
	"github.com/gabe565/ascii-movie/internal/movie"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//go:embed movie.go.tmpl
var movieTmpl string

func main() {
	srcPath := filepath.Join(config.MovieDir, config.MovieFile)
	src, err := os.Open(srcPath)
	if err != nil {
		log.Panic(err)
	}
	defer func(src *os.File) {
		_ = src.Close()
	}(src)

	// Remove existing frames
	if err := filepath.Walk(config.OutputDir, func(path string, info fs.FileInfo, err error) error {
		if filepath.Ext(path) == ".go" && filepath.Base(path) != "stub.go" {
			return os.Remove(path)
		}
		return nil
	}); err != nil {
		log.Panic(err)
	}

	var frameCap int
	m := &movie.Movie{Filename: config.MovieFile}
	var f movie.Frame
	scanner := bufio.NewScanner(src)

	// Build part of every frame, excluding progress bar and bottom padding
	for lineNum := 0; scanner.Scan(); lineNum += 1 {
		frameLineNum := lineNum % config.FrameHeight
		if frameLineNum == 0 {
			f = movie.Frame{
				Num:  lineNum / config.FrameHeight,
				Data: strings.Repeat("\n", config.PadTop-1),
			}

			v, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Panic(err)
			}

			f.Duration = time.Duration(float64(v)*(1000.0/15.0)) * time.Millisecond
		} else {
			f.Data += "\n" + strings.Repeat(" ", config.PadLeft) + scanner.Text()
		}

		if frameLineNum == config.FrameHeight-1 {
			m.Frames = append(m.Frames, f)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Panic(err)
	}

	// Compute the total duration
	totalDuration := m.Duration(1)

	// Build the rest of every frame and write to disk
	var currentPosition time.Duration
	for i, f := range m.Frames {
		f.Data += strings.Repeat("\n", config.PadBottom)
		f.Data += strings.Repeat(" ", config.PadLeft-1)
		f.Data += progressBar(currentPosition+f.Duration/2, totalDuration, config.Width)
		f.Data += strings.Repeat(" ", config.PadLeft-1)
		f.Data += strings.Repeat("\n", config.PadBottom)
		f.Height = strings.Count(f.Data, "\n")
		m.Frames[i] = f
		if frameCap < len(f.Data) {
			frameCap = len(f.Data)
		}
		currentPosition += f.Duration
	}

	m.Cap = frameCap

	// Write frame list
	if err := writeMovie(m); err != nil {
		log.Panic(err)
	}
}
