package stream

import (
	"github.com/gabe565/ascii-movie/cmd/get/stream/count"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stream",
		Aliases: []string{"streams", "connection", "connections", "client", "clients"},
		Short:   "Fetches stream metrics from a running server.",
	}
	cmd.AddCommand(
		count.NewCommand(),
	)
	return cmd
}
