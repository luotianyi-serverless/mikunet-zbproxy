package proxyprotocol

import (
	"io"
	"net"
	"time"
)

type connReader struct {
	io.Reader
}

var _ net.Conn = (*connReader)(nil)

func (c *connReader) Write(b []byte) (n int, err error) {
	panic("")
}

func (c *connReader) Close() error {
	return nil
}

func (c *connReader) LocalAddr() net.Addr {
	return nil
}

func (c *connReader) RemoteAddr() net.Addr {
	return nil
}

func (c *connReader) SetDeadline(t time.Time) error {
	return nil
}

func (c *connReader) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *connReader) SetWriteDeadline(t time.Time) error {
	return nil
}
