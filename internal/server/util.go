package server

import (
	"net"
)

func RemoteIP(addr net.Addr) string {
	switch addr := addr.(type) {
	case *net.TCPAddr:
		return addr.IP.String()
	default:
		ipPort := addr.String()
		ip, _, err := net.SplitHostPort(ipPort)
		if err != nil {
			ip = ipPort
		}
		return ip
	}
}
