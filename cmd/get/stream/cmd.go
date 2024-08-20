package stream

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/gabe565/ascii-movie/cmd/util"
	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/gabe565/ascii-movie/internal/server"
	"github.com/spf13/cobra"
)

const (
	CountActive = "active"
	CountTotal  = "total"
)

func NewCommand(opts ...util.Option) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stream",
		Aliases: []string{"streams", "connection", "connections", "client", "clients"},
		Short:   "Fetches stream metrics from a running server.",

		PreRunE: preRun,
		RunE:    run,

		Args: cobra.MaximumNArgs(1),
		ValidArgsFunction: func(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
			var completion []string
			if len(args) == 0 {
				completion = []string{"count", CountActive, CountTotal}
			}
			return completion, cobra.ShellCompDirectiveNoFileComp
		},
	}

	cmd.Flags().StringP("count", "c", "", "Gets stream count (one of "+CountActive+", "+CountTotal+")")

	for _, opt := range opts {
		opt(cmd)
	}

	return cmd
}

func preRun(cmd *cobra.Command, args []string) error {
	if len(args) != 0 {
		var count string
		switch args[0] {
		case "count", CountActive:
			count = "active"
		case CountTotal:
			count = "total"
		}
		if count != "" {
			if err := cmd.Flags().Set("count", count); err != nil {
				panic(err)
			}
		}
	}
	return nil
}

var (
	ErrInvalidURL      = errors.New("invalid URL from context")
	ErrInvalidResponse = errors.New("invalid response status")
	ErrEmptyValue      = errors.New("unexpected nil value")
)

func run(cmd *cobra.Command, _ []string) error {
	countFlag, err := cmd.Flags().GetString("count")
	if err != nil {
		panic(err)
	}

	u, ok := cmd.Context().Value(config.URLContextKey).(*url.URL)
	if !ok {
		panic(fmt.Errorf("%w: %v", ErrInvalidURL, u))
	}

	u, err = u.Parse("/streams")
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(cmd.Context(), http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", ErrInvalidResponse, resp.Status)
	}

	var decoded server.StreamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	switch countFlag {
	case CountActive:
		// Print active count
		if decoded.Active == nil {
			return fmt.Errorf("%w: count", ErrEmptyValue)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), *decoded.Active); err != nil {
			return err
		}
	case CountTotal:
		// Print total count
		if decoded.Total == nil {
			return fmt.Errorf("%w: count", ErrEmptyValue)
		}

		if _, err := fmt.Fprintln(cmd.OutOrStdout(), *decoded.Total); err != nil {
			return err
		}
	default:
		// Print list
		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
		if _, err := fmt.Fprintln(w, "SERVER\tIP\tCONNECTED\tDURATION\t"); err != nil {
			return err
		}

		if decoded.Streams == nil {
			_ = w.Flush()
			return fmt.Errorf("%w: streams", ErrEmptyValue)
		}

		streams := *decoded.Streams
		sort.Slice(streams, func(i, j int) bool {
			return streams[i].Connected.Compare(streams[j].Connected) < 0
		})

		for _, stream := range streams {
			if _, err := fmt.Fprintf(
				w,
				"%s\t%s\t%s\t%s\t\n",
				stream.Server,
				stream.RemoteIP,
				stream.Connected.Truncate(time.Second),
				time.Since(stream.Connected).Truncate(time.Second),
			); err != nil {
				return err
			}
		}
		return w.Flush()
	}
	return nil
}
