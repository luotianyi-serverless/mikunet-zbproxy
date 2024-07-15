package bufio

import "net"

type PeekConn interface {
	net.Conn
	Peek(n int) ([]byte, error)
	Rewind(position int)
	CurrentPosition() int
}
