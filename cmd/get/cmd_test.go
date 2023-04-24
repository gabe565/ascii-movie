package get

import (
	"net/url"
	"testing"

	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/gabe565/ascii-movie/internal/server"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func Test_preRun(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{"simple", args{cmd: NewCommand()}, "http://127.0.0.1", assert.NoError},
		{"invalid", args{cmd: NewCommand()}, "\x00", assert.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.cmd.PersistentFlags().Set(server.ApiFlagPrefix+server.AddressFlag, tt.want)
			if !assert.NoError(t, err) {
				return
			}

			if err := tt.args.cmd.ParseFlags(tt.args.args); !assert.NoError(t, err) {
				return
			}

			err = preRun(tt.args.cmd, tt.args.args)
			tt.wantErr(t, err)
			if err != nil {
				return
			}

			u, ok := tt.args.cmd.Context().Value(config.UrlContextKey).(*url.URL)
			if !assert.True(t, ok) {
				return
			}

			assert.Equal(t, tt.want, u.String())
		})
	}
}
