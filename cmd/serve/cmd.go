package serve

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gabe565.com/ascii-movie/cmd/util"
	"gabe565.com/ascii-movie/internal/movie"
	"gabe565.com/ascii-movie/internal/server"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewCommand(opts ...util.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve [movie]",
		Aliases: []string{"server", "listen"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Serve an ASCII movie over Telnet and SSH.",
		RunE:    run,

		ValidArgsFunction: movie.CompleteMovieName,
	}

	movie.Flags(cmd.Flags())
	server.Flags(cmd.Flags())

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

var ErrAllDisabled = errors.New("all server types are disabled")

func run(cmd *cobra.Command, args []string) error {
	if parent := cmd.Parent(); parent != nil {
		slog.Info("ASCII Movie",
			"version", parent.Annotations["version"],
			"commit", parent.Annotations["commit"],
		)
	}

	var path string
	if len(args) > 0 {
		path = args[0]
	}

	lipgloss.SetColorProfile(termenv.ANSI256)

	m, err := movie.FromFlags(cmd.Flags(), path)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)

	api := server.NewAPI(cmd.Flags())

	if ssh := server.NewSSH(cmd.Flags()); ssh.Enabled {
		api.SSHEnabled = true
		group.Go(func() error {
			return ssh.Listen(ctx, &m)
		})
	}

	if telnet := server.NewTelnet(cmd.Flags()); telnet.Enabled {
		api.TelnetEnabled = true
		server.LoadDeprecated(cmd.Flags())
		group.Go(func() error {
			return telnet.Listen(ctx, &m)
		})
	}

	if !api.SSHEnabled && !api.TelnetEnabled {
		return ErrAllDisabled
	}

	if api.Enabled {
		group.Go(func() error {
			return api.Listen(ctx)
		})
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
