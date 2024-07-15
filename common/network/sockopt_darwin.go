package network

import (
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	// TCP_FASTOPEN_SERVER is the value to enable TCP fast open on darwin for server connections.
	TCP_FASTOPEN_SERVER = 0x01
	// TCP_FASTOPEN_CLIENT is the value to enable TCP fast open on darwin for client connections.
	TCP_FASTOPEN_CLIENT = 0x02 //nolint: revive,stylecheck
)

func NewDialerControlFromOptions(option *OutboundSocketOptions) ControlFunc {
	if option == nil {
		return nil
	}
	return func(network string, address string, c syscall.RawConn) (err error) {
		err_ := c.Control(func(fd uintptr) {
			fdInt := int(fd)

			if strings.HasPrefix(network, "tcp") {
				if option.TCPFastOpen {
					err = syscall.SetsockoptInt(fdInt, syscall.IPPROTO_TCP, unix.TCP_FASTOPEN, TCP_FASTOPEN_CLIENT)
					if err != nil {
						return
					}
				}
			}
		})
		if err_ != nil {
			return err_
		}
		return err
	}
}

func NewListenerControlFromOptions(option *InboundSocketOptions) ControlFunc {
	if option == nil {
		return nil
	}
	return func(network string, address string, c syscall.RawConn) (err error) {
		err_ := c.Control(func(fd uintptr) {
			fdInt := int(fd)

			if strings.HasPrefix(network, "tcp") {
				if option.TCPFastOpen {
					err = syscall.SetsockoptInt(fdInt, syscall.IPPROTO_TCP, unix.TCP_FASTOPEN, TCP_FASTOPEN_SERVER)
					if err != nil {
						return
					}
				}
			}
		})
		if err_ != nil {
			return err_
		}
		return err
	}
}
