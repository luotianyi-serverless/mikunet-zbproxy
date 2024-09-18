// Package proxyprotocol partially implements PROXY protocol proposed by HAProxy,
// aiming at high-performance integration with the rest of ZBProxy.
// The full protocol specification can be found at https://www.haproxy.org/download/1.8/doc/proxy-protocol.txt
package proxyprotocol

import (
	"errors"
	"net/netip"
)

const (
	VersionUnspecified uint8 = iota
	Version1
	Version2
)

const (
	maskVersion  = 0xF0
	maskVersion2 = 0x20

	maskCommand      = 0x0F
	maskCommandLocal = 0x0
	maskCommandProxy = 0x1

	transportProtocolUnspecified       = 0x00
	maskTransportProtocolAddressFamily = 0xF0
	maskTransportProtocolType          = 0xF
	transportProtocolIPv4              = 0x10
	transportProtocolIPv6              = 0x20
	transportProtocolUnix              = 0x30
	transportProtocolStream            = 0x1
	transportProtocolDatagram          = 0x2
)

var (
	ErrNotProxyProtocol = errors.New("not PROXY protocol")
)

type Header struct {
	Version           uint8
	Command           uint8
	TransportProtocol uint8
	SourceAddress     netip.AddrPort
	//DestinationAddress netip.AddrPort
}
