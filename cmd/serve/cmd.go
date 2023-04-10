package serve

import (
	"context"
	"fmt"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/gabe565/ascii-movie/internal/server"
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

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)

	var serving bool

	if ssh := server.NewSSH(cmd.Flags()); ssh.Enabled {
		serving = true
		group.Go(func() error {
			return ssh.Listen(ctx, m)
		})
	}

	if telnet := server.NewTelnet(cmd.Flags()); telnet.Enabled {
		serving = true
		group.Go(func() error {
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
