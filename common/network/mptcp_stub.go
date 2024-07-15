//go:build !go1.21

package network

import "net"

const mptcpRequire121 = "MultiPath TCP requires go1.21, please recompile your binary."

func SetDialerMultiPathTCP(*net.Dialer, bool) {
	panic(mptcpRequire121)
}

func SetListenerMultiPathTCP(*net.ListenConfig, bool) {
	panic(mptcpRequire121)
}
