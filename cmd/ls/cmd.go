package ls

import (
	"fmt"
	"io/fs"
	"log/slog"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	cmdutil "github.com/gabe565/ascii-movie/cmd/util"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/gabe565/ascii-movie/internal/util"
	"github.com/spf13/cobra"
)

func NewCommand(opts ...cmdutil.Option) *cobra.Command {
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
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	movieInfos := make([]movie.Info, 0, len(args))

	if len(args) != 0 {
		for _, arg := range args {
			if err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
				if err != nil || d.IsDir() || !util.HasMovieExt(path) {
					return err
				}

				movieInfo, err := movie.GetInfo(nil, path)
				if err != nil {
					slog.Warn("Failed to get movie info", "path", arg)
					return err
				}
				movieInfos = append(movieInfos, movieInfo)
				return nil
			}); err != nil {
				return err
			}
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
