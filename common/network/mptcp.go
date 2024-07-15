//go:build go1.21

package network

import "net"

func SetDialerMultiPathTCP(dialer *net.Dialer, use bool) {
	dialer.SetMultipathTCP(use)
}

func SetListenerMultiPathTCP(listenConfig *net.ListenConfig, use bool) {
	listenConfig.SetMultipathTCP(use)
}
