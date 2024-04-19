package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteIp(t *testing.T) {
	type args struct {
		remoteIPPort string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"127.0.0.1", args{"127.0.0.1"}, "127.0.0.1"},
		{"127.0.0.1:12345", args{"127.0.0.1:12345"}, "127.0.0.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, RemoteIP(tt.args.remoteIPPort), "RemoteIP(%v)", tt.args.remoteIPPort)
		})
	}
}
