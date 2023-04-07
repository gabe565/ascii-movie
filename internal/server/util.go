package server

import "net"

func RemoteIp(remoteIpPort string) string {
	remoteIP, _, err := net.SplitHostPort(remoteIpPort)
	if err != nil {
		remoteIP = remoteIpPort
	}
	return remoteIP
}
