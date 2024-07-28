package play

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/gabe565/ascii-movie/internal/player"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "play [movie]",
		Short: "Play an ASCII movie locally.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  run,

		ValidArgsFunction: movie.CompleteMovieName,
	}

	movie.Flags(cmd.Flags())

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	if !cmd.Flags().Changed(config.LogLevelFlag) {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}

	var path string
	if len(args) > 0 {
		path = args[0]
	}

	m, err := movie.FromFlags(cmd.Flags(), path)
	if err != nil {
		return err
	}

	p := player.NewPlayer(&m, log.Level(zerolog.ErrorLevel), nil)
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
