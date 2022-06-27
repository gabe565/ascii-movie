package main

import (
	"bufio"
	_ "embed"
	"github.com/gabe565/ascii-telnet-go/config"
	"github.com/gabe565/ascii-telnet-go/internal/frame"
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

	if err := os.RemoveAll(config.OutputDir); err != nil {
		log.Fatal(err)
	}

	if err := os.Mkdir(config.OutputDir, 0777); err != nil {
		log.Fatal(err)
	}

	totalLines, err := countNewlines(in)
	if err != nil {
		log.Fatal(err)
	}

	var f *frame.Frame
	var i int
	scan := bufio.NewScanner(in)
	for scan.Scan() {
		j := i % config.FrameHeight
		if j == 0 {
			f = &frame.Frame{
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
			f.Data += strings.Repeat("\n", config.PadBottom)
			f.Data += strings.Repeat(" ", config.PadLeft-1)
			f.Data += progressBar(i, totalLines, config.Width)
			f.Data += strings.Repeat(" ", config.PadLeft-1)
			f.Data += strings.Repeat("\n", config.PadBottom)
			f.Height = strings.Count(f.Data, "\n")
			if err := writeFrame(*f); err != nil {
				log.Fatal(err)
			}
		}

		i += 1
	}

	if err := writeFrameList(totalLines / config.FrameHeight); err != nil {
		log.Fatal(err)
	}
}
