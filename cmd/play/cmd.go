package play

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/gabe565/ascii-movie/internal/movie"
	log "github.com/sirupsen/logrus"
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

	program := tea.NewProgram(movie.NewPlayer(&m, nil, nil))
	if _, err := program.Run(); err != nil {
		return err
	}
	return nil
}
