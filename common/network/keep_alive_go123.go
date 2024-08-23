//go:build go1.23

package network

import "net"

type KeepAliveConfig = net.KeepAliveConfig

func SetDialerTCPKeepAlive(dialer *net.Dialer, config KeepAliveConfig) {
	if config.Idle < 0 && config.Interval < 0 && config.Count < 0 {
		dialer.KeepAlive = -1
	} else if config.Enable {
		dialer.KeepAliveConfig = config
	}
}

func SetListenerTCPKeepAlive(listenConfig *net.ListenConfig, config KeepAliveConfig) {
	if config.Idle < 0 && config.Interval < 0 && config.Count < 0 {
		listenConfig.KeepAlive = -1
	} else if config.Enable {
		listenConfig.KeepAliveConfig = config
	}
}
