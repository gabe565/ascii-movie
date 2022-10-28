package serve

import (
	"github.com/gabe565/ascii-telnet-go/internal/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve movie to telnet clients",
		RunE:  run,
	}

	cmd.Flags().StringP("address", "a", ":23", "Listen address")

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	addr, err := cmd.Flags().GetString("address")
	if err != nil {
		return err
	}

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer func(listen net.Listener) {
		_ = listen.Close()
	}(listen)

	log.WithField("address", addr).Info("listening for connections")

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.WithError(err).Error("failed to accept connection")
			continue
		}

		go server.Serve(conn)
	}
}
