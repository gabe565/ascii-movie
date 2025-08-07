package movie

import (
	"slices"
	"testing"
	"time"

	"gabe565.com/ascii-movie/movies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFile(t *testing.T) {
	const testFile = "sw1.txt"

	f, err := movies.Movies.Open(testFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = f.Close()
	})

	movie := NewMovie()
	require.NoError(t, movie.LoadFile(testFile, f, 1))

	assert.Equal(t, testFile, movie.Filename)
	assert.Equal(t, 17*time.Minute+44*time.Second, movie.Duration().Truncate(time.Second))
	assert.Len(t, movie.Frames, 3410)
	assert.Equal(t, 67, movie.Width)
	assert.Equal(t, 13, movie.Height)
	assert.Len(t, movie.Sections, 68)
	var current time.Duration
	totalDuration := movie.Duration()
	for i, frame := range movie.Frames {
		current += frame.Duration
		if sectionIdx := slices.Index(movie.Sections, i); sectionIdx != -1 {
			timeSection := int(current * time.Duration(movie.Width) / totalDuration)
			assert.Equal(t, sectionIdx, timeSection)
		}
	}
}
