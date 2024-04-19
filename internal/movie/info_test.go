package movie

import (
	"testing"

	"github.com/gabe565/ascii-movie/movies"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetInfo(t *testing.T) {
	got, err := GetInfo(movies.Movies, "sw1.txt")
	require.NoError(t, err)
	assert.Equal(t, "sw1.txt", got.Path)
	assert.Equal(t, "sw1", got.Name)
	assert.True(t, got.Default)
	assert.NotZero(t, got.Size)
	assert.NotZero(t, got.Duration)
	assert.NotZero(t, got.NumFrames)
}

func TestListEmbedded(t *testing.T) {
	got, err := ListEmbedded()
	require.NoError(t, err)
	assert.NotEmpty(t, got)
}
