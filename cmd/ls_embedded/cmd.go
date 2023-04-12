package ls_embedded

import (
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gabe565/ascii-movie/internal/movie"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ls-embedded",
		Short: "Lists embedded movies.",
		RunE:  run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	movieInfos, err := movie.ListEmbedded()
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
	if _, err := fmt.Fprintln(w, "NAME\tSIZE\tDEFAULT\tDURATION\tFRAME COUNT\t"); err != nil {
		return err
	}
	for _, info := range movieInfos {
		if _, err := fmt.Fprintf(
			w,
			"%s\t%s\t%t\t%s\t%d\t\n",
			info.Name,
			humanize.Bytes(uint64(info.Size)),
			info.Default,
			info.Duration.Round(time.Second),
			info.NumFrames,
		); err != nil {
			return err
		}
	}
	return w.Flush()
}
