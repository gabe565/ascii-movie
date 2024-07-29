package ls

import (
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls [PATH]...",
		Aliases: []string{"ls-embedded"},
		Short:   "Lists movie files and metadata.",
		Long:    "Lists movie files and metadata.\nIf no path is given, embedded movies are listed.",
		RunE:    run,

		ValidArgsFunction: func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
			return []string{".txt", ".txt.gz"}, cobra.ShellCompDirectiveFilterFileExt
		},
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	movieInfos := make([]movie.Info, 0, len(args))

	if len(args) > 0 {
		for _, arg := range args {
			movieInfo, err := movie.GetInfo(nil, arg)
			if err != nil {
				log.Warn().Err(err).Str("path", arg).Msg("Failed to get movie info")
				continue
			}
			movieInfos = append(movieInfos, movieInfo)
		}
	} else {
		var err error
		movieInfos, err = movie.ListEmbedded()
		if err != nil {
			return err
		}
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
	if _, err := fmt.Fprintln(w, "NAME\tSIZE\tDEFAULT\tDURATION\tFRAME COUNT\tPATH\t"); err != nil {
		return err
	}
	for _, info := range movieInfos {
		if _, err := fmt.Fprintf(
			w,
			"%s\t%s\t%t\t%s\t%d\t%s\t\n",
			info.Name,
			humanize.Bytes(uint64(info.Size)),
			info.Default,
			info.Duration.Round(time.Second),
			info.NumFrames,
			info.Path,
		); err != nil {
			return err
		}
	}
	return w.Flush()
}
