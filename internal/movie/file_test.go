package movie

import (
	"slices"
	"testing"
	"time"

	"github.com/gabe565/ascii-movie/movies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFile(t *testing.T) {
	const TestFile = "sw1.txt"

	f, err := movies.Movies.Open(TestFile)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = f.Close()
	})

	movie := NewMovie()
	require.NoError(t, movie.LoadFile(TestFile, f, 1))

	assert.Equal(t, TestFile, movie.Filename)
	assert.EqualValues(t, 17*time.Minute+44*time.Second, movie.Duration().Truncate(time.Second))
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
