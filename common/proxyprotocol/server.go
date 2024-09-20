package proxyprotocol

import (
	"bytes"
	"time"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/bufio"
)

// HandleConnection reads PROXY protocol header from the connection, returning whether
// the source address is changed and the error (if present).
func HandleConnection(conn *bufio.CachedConn, metadata *adapter.Metadata) (bool, error) {
	header, err := ReadHeader(conn)
	if err != nil {
		return false, common.Cause("read PROXY protocol header: ", err)
	}
	if header.Command == CommandLocal || header.TransportProtocol == transportProtocolUnspecified {
		return false, nil
	}
	if header.SourceAddress.IsValid() {
		metadata.SourceAddress = header.SourceAddress
	}
	return true, nil
}

var (
	headerIdentityV1 = []byte("PROXY ")                 // len=6
	headerIdentityV2 = []byte("\r\n\r\n\x00\r\nQUIT\n") // len=12
)

func ReadHeader(conn *bufio.CachedConn) (*Header, error) {
	defer conn.SetReadDeadline(time.Time{})
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	identity, err := conn.Peek(6)
	if err != nil {
		return nil, common.Cause("read identity: ", err)
	}
	if bytes.Equal(identity, headerIdentityV1) {
		return readHeader1(conn)
	} else if bytes.Equal(identity, headerIdentityV2[:6]) {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		identity, err = conn.Peek(6) // 12-6=6
		if err != nil {
			return nil, common.Cause("read identity: ", err)
		}
		if bytes.Equal(identity, headerIdentityV2[6:]) {
			return readHeader2(conn)
		}
	}
	return nil, ErrNotProxyProtocol
}
