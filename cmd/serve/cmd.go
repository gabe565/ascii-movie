package serve

import (
	"context"
	"fmt"
	"github.com/gabe565/ascii-movie/internal/movie"
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

	movie.Flags(cmd.Flags())
	server.Flags(cmd.Flags())

	return cmd
}

func run(cmd *cobra.Command, args []string) (err error) {
	m, err := movie.FromFlags(cmd.Flags())
	if err != nil {
		return err
	}

	log.WithField("duration", m.Duration()).
		Info("Movie info")

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)

	var serving bool

	if server.SSHEnabled(cmd.Flags()) {
		serving = true
		group.Go(func() error {
			ssh, err := server.NewSSH(cmd.Flags())
			if err != nil {
				return err
			}
			return ssh.Listen(ctx, m)
		})
	}

	if server.TelnetEnabled(cmd.Flags()) {
		serving = true
		group.Go(func() error {
			telnet, err := server.NewTelnet(cmd.Flags())
			if err != nil {
				return err
			}
			return telnet.Listen(ctx, m)
		})
	}

	if !serving {
		return fmt.Errorf("all server types were disabled")
	}

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

	return group.Wait()
}
