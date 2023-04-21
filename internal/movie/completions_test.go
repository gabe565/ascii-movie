package movie

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestCompleteMovieName(t *testing.T) {
	got, shellComp := CompleteMovieName(&cobra.Command{}, []string{}, "")
	assert.NotEmpty(t, got)
	assert.Equal(t, cobra.ShellCompDirectiveDefault, shellComp)

	got, shellComp = CompleteMovieName(&cobra.Command{}, []string{}, "movie-that-does-not-exist")
	assert.Equal(t, []string{"txt", "txt.gz"}, got)
	assert.Equal(t, cobra.ShellCompDirectiveFilterFileExt, shellComp)
}
