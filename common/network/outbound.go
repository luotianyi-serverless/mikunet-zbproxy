// Package network contains all protocols for remote connection used by zbproxy.
//
// To implement an outbound protocol, one needs to do the following:
// 1. Implement the interface(s) below.
package network

import (
	"context"
	"io"
	"net"
)

type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
	//	DialTCPContext(ctx context.Context, network, address string) (*net.TCPConn, error)
}

type Handshake interface {
	Handshake(r io.Reader, w io.Writer, network, address string) error
}
