package network

import (
	"strings"
	"syscall"

	"golang.org/x/sys/unix"
)

func NewDialerControlFromOptions(option *OutboundSocketOptions) ControlFunc {
	if option == nil {
		return nil
	}
	return func(network string, address string, c syscall.RawConn) (err error) {
		err_ := c.Control(func(fd uintptr) {
			fdInt := int(fd)

			if option.Mark != 0 {
				err = syscall.SetsockoptInt(fdInt, syscall.SOL_SOCKET, syscall.SO_MARK, option.Mark)
				if err != nil {
					return
				}
			}

			if option.Interface != "" {
				err = syscall.BindToDevice(fdInt, option.Interface)
				if err != nil {
					return
				}
			}

			if strings.HasPrefix(network, "tcp") {
				if option.TCPFastOpen {
					err = syscall.SetsockoptInt(fdInt, syscall.SOL_TCP, unix.TCP_FASTOPEN_CONNECT, 1)
					if err != nil {
						return
					}
				}

				if option.TCPCongestion != "" {
					err = syscall.SetsockoptString(fdInt, syscall.SOL_TCP, syscall.TCP_CONGESTION, option.TCPCongestion)
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

			if option.Mark != 0 {
				err = syscall.SetsockoptInt(fdInt, syscall.SOL_SOCKET, syscall.SO_MARK, option.Mark)
				if err != nil {
					return
				}
			}

			if strings.HasPrefix(network, "tcp") {
				if option.TCPFastOpen {
					err = syscall.SetsockoptInt(fdInt, syscall.SOL_TCP, unix.TCP_FASTOPEN, 1)
					if err != nil {
						return
					}
				}

				if option.TCPCongestion != "" {
					err = syscall.SetsockoptString(fdInt, syscall.SOL_TCP, syscall.TCP_CONGESTION, option.TCPCongestion)
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
