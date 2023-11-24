package movie

import (
	"testing"

	"github.com/muesli/termenv"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewPlayer(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		movie := NewMovie()
		player := NewPlayer(&movie, log.WithField("test", t.Name()), termenv.ColorProfile())
		assert.Equal(t, &movie, player.movie)
		assert.NotNil(t, player.log)
		assert.NotEmpty(t, player.durationHook)
	})

	t.Run("no logger", func(t *testing.T) {
		movie := NewMovie()
		player := NewPlayer(&movie, nil, termenv.ColorProfile())
		assert.Equal(t, &movie, player.movie)
		assert.Nil(t, player.log)
		assert.Empty(t, player.durationHook)
	})
}
