package serve

import (
	"context"
	"github.com/gabe565/ascii-movie/internal/generated_movie"
	"github.com/gabe565/ascii-movie/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve",
		Aliases: []string{"server", "listen"},
		Short:   "Serve an ASCII movie over Telnet and SSH.",
		RunE:    run,
	}

	server.PlayFlags(cmd.Flags())
	server.ServeFlags(cmd.Flags())

	return cmd
}

func run(cmd *cobra.Command, args []string) (err error) {
	handler, err := server.New(cmd.Flags(), true)
	if err != nil {
		return err
	}

	log.WithField("duration", generated_movie.Movie.Duration(handler.Speed)).
		Info("Movie info")

	ctx, cancel := context.WithCancel(cmd.Context())
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		select {
		case <-ctx.Done():
		case <-sig:
			// Trigger shutdown
			cancel()
		}
		return nil
	})

	group.Go(func() error {
		if handler.SSHConfig.Enabled {
			return handler.ListenSSH(ctx)
		}
		return nil
	})

	group.Go(func() error {
		if handler.TelnetConfig.Enabled {
			return handler.ListenTelnet(ctx)
		}
		return nil
	})

	return group.Wait()
}
