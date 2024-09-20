package proxyprotocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net/netip"

	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/buf"
)

func (h *Header) writeHeader2(w io.Writer, destination netip.AddrPort) error {
	buffer := buf.New()
	defer buffer.Release()
	buffer.Write(headerIdentityV2)

	switch h.Command {
	case CommandLocal, CommandProxy:
	default:
		return fmt.Errorf("unknown command: %v", h.Command)
	}
	buffer.WriteByte(h.Version<<4 | h.Command)
	buffer.WriteByte(h.TransportProtocol)
	if h.Command&maskCommand == CommandLocal {
		buffer.WriteZeroN(2) // uint16
		_, err := w.Write(buffer.Bytes())
		if err != nil {
			err = common.Cause("write v2 header: ", err)
		}
		return err
	} else {
		switch h.TransportProtocol & maskTransportProtocolAddressFamily {
		case TransportProtocolIPv4:
			if !h.SourceAddress.Addr().Is4() {
				return fmt.Errorf("invalid IPv4 source address: %v", h.SourceAddress.Addr())
			}
			if destination.IsValid() {
				if !destination.Addr().Is4() {
					return fmt.Errorf("invalid IPv4 destination address: %v", destination.Addr())
				}
			} else {
				destination = netip.AddrPortFrom(netip.IPv4Unspecified(), 0)
			}
			binary.BigEndian.PutUint16(buffer.Extend(2), 12)
			buffer.Write(h.SourceAddress.Addr().AsSlice())
			buffer.Write(destination.Addr().AsSlice())
			binary.BigEndian.PutUint16(buffer.Extend(2), h.SourceAddress.Port())
			binary.BigEndian.PutUint16(buffer.Extend(2), destination.Port())

		case TransportProtocolIPv6:
			if !h.SourceAddress.Addr().Is6() {
				return fmt.Errorf("invalid IPv6 source address: %v", h.SourceAddress.Addr())
			}
			if destination.IsValid() {
				if !destination.Addr().Is6() {
					return fmt.Errorf("invalid IPv6 destination address: %v", destination.Addr())
				}
			} else {
				destination = netip.AddrPortFrom(netip.IPv6Unspecified(), 0)
			}
			binary.BigEndian.PutUint16(buffer.Extend(2), 36)
			buffer.Write(h.SourceAddress.Addr().AsSlice())
			buffer.Write(destination.Addr().AsSlice())
			binary.BigEndian.PutUint16(buffer.Extend(2), h.SourceAddress.Port())
			binary.BigEndian.PutUint16(buffer.Extend(2), destination.Port())

		case TransportProtocolUnix:
			return errors.New("transport protocol unix socket is not supported for now")

		default:
			return fmt.Errorf("unrecognized transport protocol family: %v", h.TransportProtocol&maskTransportProtocolAddressFamily)
		}
	}

	_, err := w.Write(buffer.Bytes())
	if err != nil {
		err = common.Cause("write v2 header: ", err)
	}
	return err
}
