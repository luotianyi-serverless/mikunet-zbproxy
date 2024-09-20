package proxyprotocol

import (
	"fmt"
	"io"
	"net/netip"
)

func (h *Header) WriteHeader(w io.Writer, destination netip.AddrPort) error {
	switch h.Version {
	case Version1:
		return h.writeHeader1(w, destination)
	case Version2:
		return h.writeHeader2(w, destination)
	default:
		return fmt.Errorf("unknown PROXY protocol version: %d", h.Version)
	}
}
