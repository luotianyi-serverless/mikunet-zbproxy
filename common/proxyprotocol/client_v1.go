package proxyprotocol

import (
	"fmt"
	"io"
	"net/netip"
	"strconv"

	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/buf"
)

func (h *Header) writeHeader1(w io.Writer, destination netip.AddrPort) error {
	buffer := buf.New()
	defer buffer.Release()
	buffer.Write(headerIdentityV1)

	switch h.TransportProtocol & maskTransportProtocolAddressFamily {
	case TransportProtocolIPv4:
		buffer.WriteString("TCP4 ")
	case TransportProtocolIPv6:
		buffer.WriteString("TCP6 ")
	case transportProtocolUnspecified:
		buffer.WriteString("UNKNOWN\r\n")
		_, err := w.Write(buffer.Bytes())
		if err != nil {
			err = common.Cause("write v1 header: ", err)
		}
		return err
	default:
		return fmt.Errorf("unknown address family: %d", h.TransportProtocol&maskTransportProtocolAddressFamily)
	}

	buffer.WriteString(h.SourceAddress.Addr().Unmap().WithZone("").String())
	buffer.WriteByte(' ')
	buffer.WriteString(destination.Addr().Unmap().WithZone("").String())
	buffer.WriteByte(' ')
	buffer.WriteString(strconv.FormatUint(uint64(h.SourceAddress.Port()), 10))
	buffer.WriteByte(' ')
	buffer.WriteString(strconv.FormatUint(uint64(destination.Port()), 10))
	buffer.Write(common.CRLF)

	_, err := w.Write(buffer.Bytes())
	if err != nil {
		err = common.Cause("write v1 header: ", err)
	}
	return err
}
