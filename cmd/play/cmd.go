package play

import (
	"context"
	"errors"
	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/gabe565/ascii-movie/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "play [movie]",
		Short: "Play an ASCII movie locally.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  run,
	}

	movie.Flags(cmd.Flags())

	return cmd
}

func run(cmd *cobra.Command, args []string) (err error) {
	if !cmd.Flags().Changed(config.LogLevelFlag) {
		log.SetLevel(log.WarnLevel)
	}

	var path string
	if len(args) > 0 {
		path = args[0]
	}

	m, err := movie.FromFlags(cmd.Flags(), path)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(cmd.Context())

	go server.ListenForExit(ctx, cancel, cmd.InOrStdin())

	if err := m.Stream(ctx, cmd.OutOrStdout()); err != nil {
		if !errors.Is(err, context.Canceled) {
			return err
		}
	}

	return nil
}
