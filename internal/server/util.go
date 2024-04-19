package server

import (
	"net"
)

func RemoteIP(remoteIPPort string) string {
	remoteIP, _, err := net.SplitHostPort(remoteIPPort)
	if err != nil {
		remoteIP = remoteIPPort
	}
	return remoteIP
}
