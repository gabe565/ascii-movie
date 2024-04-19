package movie

import (
	"testing"
	"time"

	"github.com/gabe565/ascii-movie/movies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFile(t *testing.T) {
	const TestFile = "short_intro.txt"

	f, err := movies.Movies.Open(TestFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = f.Close()
	})

	movie := NewMovie()
	require.NoError(t, movie.LoadFile(TestFile, f, 1))

	assert.Equal(t, TestFile, movie.Filename)
	assert.EqualValues(t, 3*time.Second, movie.Duration().Truncate(time.Second))
	assert.Len(t, movie.Frames, 45)
}
