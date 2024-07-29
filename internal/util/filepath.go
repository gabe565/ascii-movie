package util

import (
	"strings"
)

func HasMovieExt(path string) bool {
	return strings.HasSuffix(path, ".txt") || strings.HasSuffix(path, ".txt.gz")
}
