package serve

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"gabe565.com/ascii-movie/internal/config"
	"gabe565.com/ascii-movie/internal/movie"
	"gabe565.com/ascii-movie/internal/server"
	"gabe565.com/utils/cobrax"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewCommand(conf *config.Config, opts ...cobrax.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "serve [movie]",
		Aliases: []string{"server", "listen"},
		Args:    cobra.MaximumNArgs(1),
		Short:   "Serve an ASCII movie over Telnet and SSH.",
		RunE:    run,

		ValidArgsFunction: movie.CompleteMovieName,
	}

	conf.RegisterPlayFlags(cmd)
	conf.Server.RegisterFlags(cmd)

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

var ErrAllDisabled = errors.New("all server types are disabled")

func run(cmd *cobra.Command, args []string) error {
	conf, err := config.Load(cmd)
	if err != nil {
		return err
	}

	if parent := cmd.Parent(); parent != nil {
		slog.Info("ASCII Movie",
			"version", cobrax.GetVersion(cmd),
			"commit", cobrax.GetCommit(cmd),
		)
	}

	var path string
	if len(args) > 0 {
		path = args[0]
	}

	lipgloss.SetColorProfile(termenv.ANSI256)

	m, err := movie.Load(path, conf.Speed)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(cmd.Context())
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)

	api := server.NewAPI(conf)

	if conf.Server.SSH.Enabled {
		ssh := server.NewSSH(conf, api.Info)
		api.SSHEnabled = true
		group.Go(func() error {
			return ssh.Listen(ctx, &m)
		})
	}

	if conf.Server.Telnet.Enabled {
		telnet := server.NewTelnet(conf, api.Info)
		api.TelnetEnabled = true
		group.Go(func() error {
			return telnet.Listen(ctx, &m)
		})
	}

	if !api.SSHEnabled && !api.TelnetEnabled {
		return ErrAllDisabled
	}

	if conf.Server.API.Enabled {
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
