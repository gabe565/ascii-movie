package movie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromFlags(t *testing.T) {
	t.Run("default embedded", func(t *testing.T) {
		t.Parallel()
		_, err := Load("", 1)
		require.NoError(t, err)
	})

	t.Run("short_intro embedded", func(t *testing.T) {
		t.Parallel()
		_, err := Load("short_intro", 1)
		require.NoError(t, err)
	})

	t.Run("short_intro file", func(t *testing.T) {
		t.Parallel()
		_, err := Load("../../movies/short_intro.txt", 1)
		require.NoError(t, err)
	})

	t.Run("invalid speed", func(t *testing.T) {
		t.Parallel()
		_, err := Load("", -1)
		require.Error(t, err)
	})
}
