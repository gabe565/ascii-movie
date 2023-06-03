package stream

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_preRun(t *testing.T) {
	t.Run("without count", func(t *testing.T) {
		cmd := NewCommand()
		if err := preRun(cmd, []string{}); !assert.NoError(t, err) {
			return
		}
		got, err := cmd.Flags().GetString("count")
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "", got)
	})

	t.Run("with active count", func(t *testing.T) {
		cmd := NewCommand()
		if err := preRun(cmd, []string{"count"}); !assert.NoError(t, err) {
			return
		}
		got, err := cmd.Flags().GetString("count")
		if !assert.NoError(t, err) {
			return
		}
		assert.Equal(t, "active", got)
	})
}

func Test_run(t *testing.T) {
	countCmd := func() *cobra.Command {
		cmd := NewCommand()
		if err := cmd.Flags().Set("count", "active"); !assert.NoError(t, err) {
			return cmd
		}
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
		wantErr assert.ErrorAssertionFunc
	}{
		{"0", args{cmd: countCmd()}, "0\n", false, assert.NoError},
		{"1", args{cmd: countCmd()}, "1\n", false, assert.NoError},
		{"http error", args{cmd: countCmd()}, "", true, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.want500 {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				_, _ = w.Write([]byte(`{"active":` + strings.TrimSuffix(tt.want, "\n") + `}`))
			}))
			defer svr.Close()

			u, err := url.Parse(svr.URL)
			if !assert.NoError(t, err) {
				return
			}

			ctx := context.WithValue(context.Background(), config.UrlContextKey, u)
			tt.args.cmd.SetContext(ctx)

			var buf strings.Builder
			tt.args.cmd.SetOut(&buf)

			if err := run(tt.args.cmd, tt.args.args); !tt.wantErr(t, err) {
				return
			}

			assert.Equal(t, tt.want, buf.String())
		})
	}
}
