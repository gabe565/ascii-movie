package cmd

import (
	"context"

	"gabe565.com/ascii-movie/cmd/get"
	"gabe565.com/ascii-movie/cmd/ls"
	"gabe565.com/ascii-movie/cmd/play"
	"gabe565.com/ascii-movie/cmd/serve"
	"gabe565.com/ascii-movie/internal/config"
	"gabe565.com/utils/cobrax"
	"github.com/spf13/cobra"
)

func NewCommand(opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ascii-movie",
		Short: "Command line ASCII movie player.",

		DisableAutoGenTag: true,
		SilenceUsage:      true,
		SilenceErrors:     true,
	}

	conf := config.New()
	conf.RegisterFlags(cmd)
	if cmd.Context() == nil {
		cmd.SetContext(context.Background())
	}
	cmd.SetContext(config.NewContext(cmd.Context(), conf))

	cmd.AddCommand(
		play.NewCommand(conf),
		serve.NewCommand(conf),
		ls.NewCommand(),
		get.NewCommand(),
	)

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}
