package proxyprotocol

import (
	"bytes"
	"net/netip"
	"testing"

	"github.com/layou233/zbproxy/v3/common/bufio"
)

func TestServerV1(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		source   netip.AddrPort
		protocol uint8
		unknown  bool
	}{
		{
			name:     "TCP4 Minimal",
			header:   "PROXY TCP4 1.1.1.5 1.1.1.6 2 3\r\n",
			source:   netip.MustParseAddrPort("1.1.1.5:2"),
			protocol: TransportProtocolStream | TransportProtocolIPv4,
		},
		{
			name:     "TCP4 Maximal",
			header:   "PROXY TCP4 255.255.255.255 255.255.255.254 65535 65535\r\n",
			source:   netip.MustParseAddrPort("255.255.255.255:65535"),
			protocol: TransportProtocolStream | TransportProtocolIPv4,
		},
		{
			name:     "TCP6 Minimal",
			header:   "PROXY TCP6 ::1 ::2 3 4\r\n",
			source:   netip.MustParseAddrPort("[::1]:3"),
			protocol: TransportProtocolStream | TransportProtocolIPv6,
		},
		{
			name:     "TCP6 Maximal",
			header:   "PROXY TCP6 0000:0000:0000:0000:0000:0000:0000:0002 0000:0000:0000:0000:0000:0000:0000:0001 65535 65535\r\n",
			source:   netip.MustParseAddrPort("[::2]:65535"),
			protocol: TransportProtocolStream | TransportProtocolIPv6,
		},
		{
			name:    "UNKNOWN Minimal",
			header:  "PROXY UNKNOWN\r\n",
			unknown: true,
		},
		{
			name:    "UNKNOWN Maximal",
			header:  "PROXY UNKNOWN 0000:0000:0000:0000:0000:0000:0000:0002 0000:0000:0000:0000:0000:0000:0000:0001 65535 65535\r\n",
			unknown: true,
			source:  netip.MustParseAddrPort("[::2]:65535"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bytes.NewReader([]byte(tt.header))
			header, err := ReadHeader(bufio.NewCachedConn(&connReader{reader}))
			if err != nil {
				t.Fatalf("failed to read header: %v", err)
			}
			if header.Version != Version1 {
				t.Fatalf("version mismatch: got=%v, expect=1", header.Version)
			}
			if tt.unknown {
				if header.TransportProtocol != transportProtocolUnspecified {
					t.Fatalf("unexpected transport protocol: got=%v, expect=unspecified", header.TransportProtocol)
				}
				if tt.source.IsValid() && header.SourceAddress != tt.source {
					t.Fatalf("unexpected source address: got=%v, expect=%v", header.SourceAddress, tt.source)
				}
			} else {
				if header.TransportProtocol != tt.protocol {
					t.Fatalf("unexpected transport protocol: got=%v, expect=%v", header.TransportProtocol, tt.protocol)
				}
				if header.SourceAddress != tt.source {
					t.Fatalf("unexpected source address: got=%v, expect=%v", header.SourceAddress, tt.source)
				}
			}
		})
	}
}
