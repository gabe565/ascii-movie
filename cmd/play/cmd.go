package play

import (
	"log/slog"

	"gabe565.com/ascii-movie/cmd/util"
	"gabe565.com/ascii-movie/internal/config"
	"gabe565.com/ascii-movie/internal/movie"
	"gabe565.com/ascii-movie/internal/player"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewCommand(opts ...util.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "play [movie]",
		Short: "Play an ASCII movie locally.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  run,

		ValidArgsFunction: movie.CompleteMovieName,
	}

	movie.Flags(cmd.Flags())

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if !cmd.Flags().Changed(config.LogLevelFlag) {
		if err := cmd.Flags().Set(config.LogLevelFlag, slog.LevelWarn.String()); err != nil {
			slog.Warn("Failed to decrease log level", "error", err)
		}
		config.InitLogCmd(cmd)
	}

	var path string
	if len(args) > 0 {
		path = args[0]
	}

	m, err := movie.FromFlags(cmd.Flags(), path)
	if err != nil {
		return err
	}

	p := player.NewPlayer(&m, slog.Default(), nil)
	defer p.Close()

	program := tea.NewProgram(p,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := program.Run(); err != nil {
		return err
	}
	return nil
}
