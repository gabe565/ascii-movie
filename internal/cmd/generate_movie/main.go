package main

import (
	_ "embed"
	"github.com/gabe565/ascii-movie/config"
	"github.com/gabe565/ascii-movie/internal/movie"
	"log"
	"path/filepath"
)

//go:embed movie.go.tmpl
var movieTmpl string

func main() {
	srcPath := filepath.Join(config.MovieDir, config.MovieFile)
	m, err := movie.NewFromFile(srcPath)
	if err != nil {
		log.Panic(err)
	}

	// Write frame list
	if err := writeMovie(m); err != nil {
		log.Panic(err)
	}
}
