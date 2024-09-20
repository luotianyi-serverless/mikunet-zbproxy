package proxyprotocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net/netip"
	"time"

	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/buf"
	"github.com/layou233/zbproxy/v3/common/bufio"
)

func readHeader2(conn *bufio.CachedConn) (*Header, error) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	versionAndCommandBuf, err := conn.Peek(1)
	if err != nil {
		return nil, common.Cause("read v2 version and command: ", err)
	}
	versionAndCommand := versionAndCommandBuf[0]
	version := versionAndCommand & maskVersion
	command := versionAndCommand & maskCommand
	if version != maskVersion2 {
		return nil, fmt.Errorf("v2 version mismatch: %v", version)
	}
	switch command {
	case CommandLocal, CommandProxy:
	default:
		return nil, fmt.Errorf("bad v2 command: %v", command)
	}
	header := &Header{
		Version: Version2,
		Command: command,
	}

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	familyAndAddressBuf, err := conn.Peek(1)
	if err != nil {
		return nil, common.Cause("read v2 family and address: ", err)
	}
	header.TransportProtocol = familyAndAddressBuf[0]

	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	restLenBuf, err := conn.Peek(2)
	if err != nil {
		return nil, common.Cause("read v2 rest length: ", err)
	}
	restLen := binary.BigEndian.Uint16(restLenBuf)
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	restBuf, err := conn.Peek(int(restLen))
	if err != nil {
		return nil, common.Cause("read v2 rest data: ", err)
	}

	if command == CommandProxy {
		buffer := buf.As(restBuf)
		protocolFamily := header.TransportProtocol & maskTransportProtocolAddressFamily
		switch protocolFamily {
		case TransportProtocolIPv4:
			rawSourceAddress, err := buffer.Peek(4)
			if err != nil {
				return nil, common.Cause("read v2 source address: ", err)
			}
			_, err = buffer.Peek(4) // destination address
			if err != nil {
				return nil, common.Cause("read v2 destination address: ", err)
			}
			rawSourcePort, err := buffer.Peek(2)
			if err != nil {
				return nil, common.Cause("read v2 source port: ", err)
			}
			sourcePort := binary.BigEndian.Uint16(rawSourcePort)
			_, err = buffer.Peek(2) // destination port
			if err != nil {
				return nil, common.Cause("read v2 destination port: ", err)
			}
			header.SourceAddress = netip.AddrPortFrom(netip.AddrFrom4(*(*[4]byte)(rawSourceAddress)), sourcePort)

		case TransportProtocolIPv6:
			rawSourceAddress, err := buffer.Peek(16)
			if err != nil {
				return nil, common.Cause("read v2 source address: ", err)
			}
			_, err = buffer.Peek(16) // destination address
			if err != nil {
				return nil, common.Cause("read v2 destination address: ", err)
			}
			rawSourcePort, err := buffer.Peek(2)
			if err != nil {
				return nil, common.Cause("read v2 source port: ", err)
			}
			sourcePort := binary.BigEndian.Uint16(rawSourcePort)
			_, err = buffer.Peek(2) // destination port
			if err != nil {
				return nil, common.Cause("read v2 destination port: ", err)
			}
			header.SourceAddress = netip.AddrPortFrom(netip.AddrFrom16(*(*[16]byte)(rawSourceAddress)).Unmap(), sourcePort)

		case TransportProtocolUnix:
			return nil, errors.New("transport protocol unix socket is not supported for now")

		default:
			return nil, fmt.Errorf("unrecognized transport protocol family: %v", protocolFamily)
		}
	}

	// we do not handle TLVs for now
	return header, nil
}
