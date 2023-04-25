package list

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"text/tabwriter"
	"time"

	"github.com/gabe565/ascii-movie/internal/config"
	"github.com/gabe565/ascii-movie/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "Fetches a list of active streams",
		RunE:    run,
	}
	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	u, ok := cmd.Context().Value(config.UrlContextKey).(*url.URL)
	if !ok {
		panic(fmt.Errorf("invalid URL from context: %v", u))
	}

	u, err := u.Parse("/streams?fields=streams")
	if err != nil {
		return err
	}

	resp, err := http.Get(u.String())
	if err != nil {
		log.Error("Failed to connect to API. Is the server running?")
		return fmt.Errorf("failed to connect to API: %w", err)
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response status: %s", resp.Status)
	}

	var decoded server.StreamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 3, ' ', 0)
	if _, err := fmt.Fprintln(w, "IP\tCONNECTED\t"); err != nil {
		return err
	}

	if decoded.Streams == nil {
		_ = w.Flush()
		return fmt.Errorf("unexpected nil value: streams")
	}

	for _, stream := range *decoded.Streams {
		if _, err := fmt.Fprintf(
			w,
			"%s\t%s\t\n",
			stream.RemoteIp,
			stream.Connected.Truncate(time.Second),
		); err != nil {
			return err
		}
	}
	return w.Flush()
}
