package movie

import (
	"testing"

	"github.com/gabe565/ascii-movie/movies"
	"github.com/stretchr/testify/assert"
)

func TestGetInfo(t *testing.T) {
	got, err := GetInfo(movies.Movies, "sw1.txt")
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, "sw1.txt", got.Path)
	assert.Equal(t, "sw1", got.Name)
	assert.Equal(t, true, got.Default)
	assert.NotZero(t, got.Size)
	assert.NotZero(t, got.Duration)
	assert.NotZero(t, got.NumFrames)
}

func TestListEmbedded(t *testing.T) {
	got, err := ListEmbedded()
	if !assert.NoError(t, err) {
		return
	}

	assert.NotEmpty(t, got)
}
