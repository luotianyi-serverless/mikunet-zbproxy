package proxyprotocol

import (
	"bytes"
	"net/netip"
	"testing"

	"github.com/layou233/zbproxy/v3/common/buf"
)

func TestClientV1(t *testing.T) {
	tests := []struct {
		name        string
		header      *Header
		destination netip.AddrPort
		expect      []byte
	}{
		{
			name: "TCP4 Minimal",
			header: &Header{
				Version:           Version1,
				TransportProtocol: TransportProtocolStream | TransportProtocolIPv4,
				SourceAddress:     netip.MustParseAddrPort("1.1.1.5:2"),
			},
			destination: netip.MustParseAddrPort("1.1.1.6:3"),
			expect:      []byte("PROXY TCP4 1.1.1.5 1.1.1.6 2 3\r\n"),
		},
		{
			name: "TCP4 Maximal",
			header: &Header{
				Version:           Version1,
				TransportProtocol: TransportProtocolStream | TransportProtocolIPv4,
				SourceAddress:     netip.MustParseAddrPort("255.255.255.255:65535"),
			},
			destination: netip.MustParseAddrPort("255.255.255.254:65535"),
			expect:      []byte("PROXY TCP4 255.255.255.255 255.255.255.254 65535 65535\r\n"),
		},
		{
			name: "TCP6",
			header: &Header{
				Version:           Version1,
				TransportProtocol: TransportProtocolStream | TransportProtocolIPv6,
				SourceAddress:     netip.MustParseAddrPort("[::1]:3"),
			},
			destination: netip.MustParseAddrPort("[::2]:4"),
			expect:      []byte("PROXY TCP6 ::1 ::2 3 4\r\n"),
		},
		{
			name: "UNKNOWN",
			header: &Header{
				Version:           Version1,
				TransportProtocol: transportProtocolUnspecified,
			},
			expect: []byte("PROXY UNKNOWN\r\n"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := buf.New()
			defer buffer.Release()
			err := tt.header.WriteHeader(buffer, tt.destination)
			if err != nil {
				t.Fatalf("failed to write v1 header: %v", err)
			}
			if !bytes.Equal(buffer.Bytes(), tt.expect) {
				t.Errorf("got=%v, expect=%v", string(buffer.Bytes()), string(tt.expect))
			}
		})
	}
}
