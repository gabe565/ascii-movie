package ls_embedded

import (
	"fmt"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/gabe565/ascii-movie/movies"
	"github.com/spf13/cobra"
	"io/fs"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls-embedded",
		Short: "Lists embedded movies.",
		RunE:  run,
	}
	return cmd
}

type MovieInfo struct {
	Name      string
	Duration  time.Duration
	Default   bool
	NumFrames int
}

func run(cmd *cobra.Command, args []string) error {
	var movieInfos []MovieInfo

	if err := fs.WalkDir(
		movies.Movies,
		".",
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			f, err := movies.Movies.Open(path)
			if err != nil {
				return err
			}

			movie, err := movie.NewFromFile(path, f, 14, movie.Padding{}, movie.Padding{})
			if err != nil {
				return err
			}

			movieInfos = append(movieInfos, MovieInfo{
				Name:      strings.TrimSuffix(path, filepath.Ext(path)),
				Duration:  movie.Duration(),
				Default:   path == movies.Default,
				NumFrames: len(movie.Frames),
			})
			return nil
		},
	); err != nil {
		return err
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
	if _, err := fmt.Fprintln(w, "NAME\tDURATION\tDEFAULT\tFRAME COUNT\t"); err != nil {
		return err
	}
	for _, info := range movieInfos {
		if _, err := fmt.Fprintf(
			w,
			"%s\t%t\t%s\t%d\t\n",
			info.Name,
			info.Default,
			info.Duration.Round(time.Second),
			info.NumFrames,
		); err != nil {
			return err
		}
	}
	return w.Flush()
}
