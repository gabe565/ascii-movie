package movie

import (
	"io/fs"
	"testing"
	"time"

	"github.com/gabe565/ascii-movie/movies"
	"github.com/stretchr/testify/assert"
)

func TestLoadFile(t *testing.T) {
	const TestFile = "short_intro.txt"

	f, err := movies.Movies.Open(TestFile)
	if !assert.NoError(t, err) {
		return
	}
	defer func(f fs.File) {
		_ = f.Close()
	}(f)

	movie := NewMovie()
	if err := movie.LoadFile(TestFile, f, 1); !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, TestFile, movie.Filename)
	assert.EqualValues(t, 3*time.Second, movie.Duration().Truncate(time.Second))
	assert.Len(t, movie.Frames, 45)
}
