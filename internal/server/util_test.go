package server

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteIp(t *testing.T) {
	type args struct {
		n net.Addr
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"127.0.0.1:12345",
			args{&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 12345}},
			"127.0.0.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, RemoteIP(tt.args.n))
		})
	}
}
