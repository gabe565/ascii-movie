package cmd

import (
	"github.com/gabe565/ascii-telnet-go/cmd/play"
	"github.com/gabe565/ascii-telnet-go/cmd/serve"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ascii-telnet",
	}
	cmd.AddCommand(
		play.NewCommand(),
		serve.NewCommand(),
	)
	return cmd
}
