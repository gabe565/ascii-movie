package stream

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"gabe565.com/ascii-movie/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_preRun(t *testing.T) {
	t.Run("without count", func(t *testing.T) {
		cmd := NewCommand()
		require.NoError(t, preRun(cmd, []string{}))
		got, err := cmd.Flags().GetString("count")
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("with active count", func(t *testing.T) {
		cmd := NewCommand()
		require.NoError(t, preRun(cmd, []string{"count"}))
		got, err := cmd.Flags().GetString("count")
		require.NoError(t, err)
		assert.Equal(t, "active", got)
	})
}

func Test_run(t *testing.T) {
	countCmd := func() *cobra.Command {
		cmd := NewCommand()
		require.NoError(t, cmd.Flags().Set("count", "active"))
		return cmd
	}

	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want500 bool
		wantErr require.ErrorAssertionFunc
	}{
		{"0", args{cmd: countCmd()}, "0\n", false, require.NoError},
		{"1", args{cmd: countCmd()}, "1\n", false, require.NoError},
		{"http error", args{cmd: countCmd()}, "", true, require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				if tt.want500 {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				_, _ = w.Write([]byte(`{"active":` + strings.TrimSuffix(tt.want, "\n") + `}`))
			}))
			t.Cleanup(svr.Close)

			u, err := url.Parse(svr.URL)
			require.NoError(t, err)

			ctx := context.WithValue(t.Context(), config.URLContextKey, u)
			tt.args.cmd.SetContext(ctx)

			var buf strings.Builder
			tt.args.cmd.SetOut(&buf)

			tt.wantErr(t, run(tt.args.cmd, tt.args.args))

			assert.Equal(t, tt.want, buf.String())
		})
	}
}
