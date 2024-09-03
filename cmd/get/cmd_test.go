package get

import (
	"net/url"
	"testing"

	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/gabe565/ascii-movie/internal/server"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_preRun(t *testing.T) {
	type args struct {
		cmd *cobra.Command
		_   []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr require.ErrorAssertionFunc
	}{
		{"simple", args{cmd: NewCommand()}, "http://127.0.0.1", require.NoError},
		{"invalid", args{cmd: NewCommand()}, "\x00", require.Error},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.cmd.PersistentFlags().Set(server.APIFlagPrefix+server.AddressFlag, tt.want)
			require.NoError(t, err)

			tt.args.cmd.RunE = func(_ *cobra.Command, _ []string) error { return nil }
			err = tt.args.cmd.Execute()
			tt.wantErr(t, err)
			if err != nil {
				return
			}

			u, ok := tt.args.cmd.Context().Value(config.URLContextKey).(*url.URL)
			require.True(t, ok)
			assert.Equal(t, tt.want, u.String())
		})
	}
}
