package get

import (
	"context"
	"net/url"

	"github.com/gabe565/ascii-movie/cmd/get/stream"
	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/gabe565/ascii-movie/internal/server"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Fetches data from a running server.",

		PersistentPreRunE: preRun,
	}

	cmd.PersistentFlags().String(server.ApiFlagPrefix+server.AddressFlag, "http://127.0.0.1:1977", "API address")

	cmd.AddCommand(
		stream.NewCommand(),
	)
	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {
	apiAddr, err := cmd.Flags().GetString(server.ApiFlagPrefix + server.AddressFlag)
	if err != nil {
		panic(err)
	}

	u, err := url.Parse(apiAddr)
	if err != nil {
		return err
	}

	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}
	ctx = context.WithValue(ctx, config.UrlContextKey, u)
	cmd.SetContext(ctx)

	return nil
}
