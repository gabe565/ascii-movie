package main

import (
	"bytes"
	"github.com/gabe565/ascii-movie/config"
	"github.com/gabe565/ascii-movie/internal/movie"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func writeMovie(m *movie.Movie) error {
	filename := strings.TrimSuffix(m.Filename, filepath.Ext(m.Filename))
	filename = filepath.Join(config.OutputDir, filename+".go")

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		_ = out.Close()
	}(out)

	tmpl, err := template.New("").Parse(movieTmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer

	err = tmpl.Execute(&buf, map[string]any{
		"Package": filepath.Base(config.OutputDir),
		"Movie":   m,
	})
	if err != nil {
		return err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}

	if _, err := out.Write(formatted); err != nil {
		return err
	}

	if err := out.Close(); err != nil {
		return err
	}

	return nil
}
