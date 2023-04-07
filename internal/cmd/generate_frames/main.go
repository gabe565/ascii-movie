package main

import (
	"bufio"
	_ "embed"
	"github.com/gabe565/ascii-telnet-go/config"
	"github.com/gabe565/ascii-telnet-go/internal/frame"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

//go:embed all.go.tmpl
var allTmpl string

//go:embed frame.go.tmpl
var frameTmpl string

func main() {
	filename := filepath.Join("movies", filepath.Base(config.MovieFile))

	in, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func(in *os.File) {
		_ = in.Close()
	}(in)

	if err := filepath.Walk(config.OutputDir, func(path string, info fs.FileInfo, err error) error {
		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "stub.go") {
			return os.Remove(path)
		}
		return nil
	}); err != nil {
		log.Fatal(err)
	}

	totalLines, err := countNewlines(in)
	if err != nil {
		log.Fatal(err)
	}

	var frameCap int
	var frames []frame.Frame
	var f frame.Frame
	var i int
	scan := bufio.NewScanner(in)
	for scan.Scan() {
		j := i % config.FrameHeight
		if j == 0 {
			f = frame.Frame{
				Num:  i / config.FrameHeight,
				Data: strings.Repeat("\n", config.PadTop-1),
			}

			v, err := strconv.ParseInt(scan.Text(), 0, 32)
			if err != nil {
				log.Fatal(err)
			}

			f.Sleep = time.Duration(float64(v)*(1000.0/15.0)) * time.Millisecond
		} else {
			f.Data += "\n" + strings.Repeat(" ", config.PadLeft) + scan.Text()
		}

		if j == config.FrameHeight-1 {
			frames = append(frames, f)
		}

		i += 1
	}

	var totalDuration time.Duration
	for _, f := range frames {
		totalDuration += f.Sleep
	}

	var currentPosition time.Duration
	for _, f := range frames {
		f.Data += strings.Repeat("\n", config.PadBottom)
		f.Data += strings.Repeat(" ", config.PadLeft-1)
		f.Data += progressBar(currentPosition, totalDuration, config.Width)
		f.Data += strings.Repeat(" ", config.PadLeft-1)
		f.Data += strings.Repeat("\n", config.PadBottom)
		f.Height = strings.Count(f.Data, "\n")
		if frameCap < len(f.Data) {
			frameCap = len(f.Data)
		}
		if err := writeFrame(f); err != nil {
			log.Fatal(err)
		}
		currentPosition += f.Sleep
	}

	if err := writeFrameList(totalLines/config.FrameHeight, frameCap); err != nil {
		log.Fatal(err)
	}
}
