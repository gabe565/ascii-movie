package count

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

func Test_run(t *testing.T) {
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
		{"0", args{cmd: NewCommand()}, "0\n", false, assert.NoError},
		{"1", args{cmd: NewCommand()}, "1\n", false, assert.NoError},
		{"http error", args{cmd: NewCommand()}, "", true, assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.want500 {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				_, _ = w.Write([]byte(`{"count":` + strings.TrimSuffix(tt.want, "\n") + `}`))
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
