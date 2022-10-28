package play

import (
	"github.com/gabe565/ascii-telnet-go/internal/server"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "play",
		Short: "Play the movie locally",
		RunE:  run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	return server.ServeAscii(cmd.OutOrStdout())
}
