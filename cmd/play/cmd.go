package play

import (
	"github.com/gabe565/ascii-movie/internal/movie"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "play",
		Short: "Play an ASCII movie locally.",
		RunE:  run,
	}

	movie.Flags(cmd.Flags())

	return cmd
}

func run(cmd *cobra.Command, args []string) (err error) {
	log.SetLevel(log.WarnLevel)

	m, err := movie.FromFlags(cmd.Flags())
	if err != nil {
		return err
	}

	return m.Stream(cmd.OutOrStdout())
}
