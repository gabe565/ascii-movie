package idleconn

import (
	"context"
	"errors"
	"net"
	"os"
	"time"
)

func New(conn net.Conn, idleTimeout, maxTimeout time.Duration, cancel context.CancelFunc) net.Conn {
	return &Conn{
		Conn:        conn,
		idleTimeout: idleTimeout,
		maxDeadline: time.Now().Add(maxTimeout),
		cancel:      cancel,
	}
}

type Conn struct {
	net.Conn
	idleTimeout time.Duration
	maxDeadline time.Time
	cancel      context.CancelFunc
}

func (c *Conn) Write(p []byte) (int, error) {
	_ = c.updateDeadline()

	n, err := c.Conn.Write(p)
	if errors.Is(err, os.ErrDeadlineExceeded) {
		c.cancel()
	}
	return n, err
}

func (c *Conn) Read(b []byte) (int, error) {
	_ = c.updateDeadline()

	n, err := c.Conn.Read(b)
	if errors.Is(err, os.ErrDeadlineExceeded) {
		c.cancel()
	}
	return n, err
}

func (c *Conn) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	return c.Conn.Close()
}

func (c *Conn) updateDeadline() error {
	if c.idleTimeout != 0 {
		idleDeadline := time.Now().Add(c.idleTimeout)
		if idleDeadline.Before(c.maxDeadline) || c.maxDeadline.IsZero() {
			return c.Conn.SetDeadline(idleDeadline)
		}
	}
	return c.Conn.SetDeadline(c.maxDeadline)
}
