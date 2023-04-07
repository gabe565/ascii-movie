package serve

import (
	"github.com/gabe565/ascii-telnet-go/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server", "listen"},
		Short:   "Serve movie to telnet clients",
		RunE:    run,
	}

	cmd.Flags().StringP("address", "a", ":23", "Listen address")
	server.Flags(cmd.Flags())

	return cmd
}

func run(cmd *cobra.Command, args []string) (err error) {
	handler, err := server.New(cmd.Flags())
	if err != nil {
		return err
	}

	log.WithField("duration", handler.MovieDuration()).Info("total movie duration")

	addr, err := cmd.Flags().GetString("address")
	if err != nil {
		return err
	}

	if err := handler.Listen(cmd.Context(), addr); err != nil {
		return err
	}

	return nil
}
