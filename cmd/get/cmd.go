package get

import (
	"context"
	"net/url"

	"github.com/gabe565/ascii-movie/cmd/get/stream"
	"github.com/gabe565/ascii-movie/cmd/util"
	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/gabe565/ascii-movie/internal/server"
	"github.com/spf13/cobra"
)

func NewCommand(opts ...util.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Fetches data from a running server.",

		PersistentPreRunE: preRun,
	}

	cmd.PersistentFlags().String(server.APIFlagPrefix+server.AddressFlag, "http://127.0.0.1:1977", "API address")

	cmd.AddCommand(
		stream.NewCommand(),
	)
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func preRun(cmd *cobra.Command, _ []string) error {
	apiAddr, err := cmd.Flags().GetString(server.APIFlagPrefix + server.AddressFlag)
	if err != nil {
		panic(err)
	}

	u, err := url.Parse(apiAddr)
	if err != nil {
		return err
	}

	cmd.SetContext(context.WithValue(cmd.Context(), config.URLContextKey, u))
	return nil
}
