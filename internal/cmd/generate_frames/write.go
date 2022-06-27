package main

import (
	"fmt"
	"github.com/gabe565/ascii-telnet-go/config"
	"github.com/gabe565/ascii-telnet-go/internal/frame"
	"os"
	"path/filepath"
	"text/template"
)

func writeFrame(f frame.Frame) error {
	filename := filepath.Join(config.OutputDir, fmt.Sprintf("frame%d.go", f.Num))

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	tmpl, err := template.New("").Parse(frameTmpl)
	if err != nil {
		return err
	}

	err = tmpl.Execute(out, map[string]any{
		"Package": config.OutputDir,
		"Frame":   f,
	})
	if err != nil {
		return err
	}

	if err := out.Close(); err != nil {
		return err
	}

	return nil
}

func writeFrameList(n int) error {
	filename := filepath.Join(config.OutputDir, "0_frame_list.go")

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	tmpl, err := template.New("").Parse(allTmpl)
	if err != nil {
		return err
	}

	frames := make([]struct{}, n)

	err = tmpl.Execute(out, map[string]any{
		"Package": config.OutputDir,
		"Frames":  frames,
	})
	if err != nil {
		return err
	}

	if err := out.Close(); err != nil {
		return err
	}

	return nil
}
