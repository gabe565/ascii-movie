package cmd

import (
	"github.com/gabe565/ascii-movie/cmd/play"
	"github.com/gabe565/ascii-movie/cmd/serve"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ascii-movie",
		Short: "Command line ASCII movie player.",

		DisableAutoGenTag: true,
	}
	cmd.AddCommand(
		play.NewCommand(),
		serve.NewCommand(),
	)
	return cmd
}
