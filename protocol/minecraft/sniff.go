package minecraft

import (
	"errors"
	"time"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/buf"
	"github.com/layou233/zbproxy/v3/common/bufio"
	"github.com/layou233/zbproxy/v3/common/mcprotocol"
)

var (
	ErrBadPacket = errors.New("bad Minecraft handshake packet")
)

func SniffClientHandshake(conn bufio.PeekConn, metadata *adapter.Metadata) error {
	if metadata.Minecraft == nil {
		metadata.Minecraft = &adapter.MinecraftMetadata{
			NextState: -1,
		}
	}
	defer conn.SetReadDeadline(time.Time{}) // clear deadline

	// handshake packet
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	packetSize, _, err := mcprotocol.ReadVarIntFrom(conn)
	if err != nil {
		return common.Cause("read packet size: ", err)
	}
	if packetSize > 264 { // maximum possible size of this kind of packet
		return ErrBadPacket
	}
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	packetContent, err := conn.Peek(int(packetSize))
	if err != nil {
		return common.Cause("read handshake packet: ", err)
	}
	buffer := buf.As(packetContent)

	var packetID byte
	packetID, err = buffer.ReadByte()
	if err != nil {
		return common.Cause("read packet ID: ", err)
	}
	if packetID != 0 { // Server bound : Handshake
		return ErrBadPacket
	}
	protocolVersion, _, err := mcprotocol.ReadVarIntFrom(buffer)
	if err != nil {
		return common.Cause("read protocol version: ", err)
	}
	if protocolVersion <= 0 {
		return ErrBadPacket
	}
	metadata.Minecraft.ProtocolVersion = uint(protocolVersion)

	metadata.Minecraft.OriginDestination, err = mcprotocol.ReadString(buffer)
	if err != nil {
		return common.Cause("read destination: ", err)
	}
	if metadata.Minecraft.OriginDestination == "" {
		return ErrBadPacket
	}

	metadata.Minecraft.OriginPort, err = mcprotocol.ReadUint16(buffer)
	if err != nil {
		return common.Cause("read port: ", err)
	}
	if metadata.Minecraft.OriginPort == 0 {
		return ErrBadPacket
	}

	nextState, err := buffer.ReadByte()
	if err != nil {
		return common.Cause("read next state: ", err)
	}
	switch nextState {
	case mcprotocol.NextStateLogin,
		mcprotocol.NextStateStatus,
		mcprotocol.NextStateTransfer:
	default:
		return ErrBadPacket
	}
	metadata.Minecraft.NextState = int8(nextState)

	if nextState == mcprotocol.NextStateStatus {
		// status packet
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		_, err = conn.Peek(2)
		if err != nil {
			return common.Cause("read status request: ", err)
		}
	} else {
		// login packet or transfer packet
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		packetSize, _, err = mcprotocol.ReadVarIntFrom(conn)
		if err != nil {
			return common.Cause("read packet size: ", err)
		}
		if packetSize > 33 { // maximum possible size of this kind of packet
			return ErrBadPacket
		}
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		packetContent, err = conn.Peek(int(packetSize))
		if err != nil {
			return common.Cause("read login packet: ", err)
		}
		buffer = buf.As(packetContent)

		packetID, err = buffer.ReadByte()
		if err != nil {
			return common.Cause("read packet ID: ", err)
		}
		if packetID != 0 { // Server bound : Login Start
			return ErrBadPacket
		}
		metadata.Minecraft.PlayerName, err = mcprotocol.ReadLimitedString(buffer, 16)
		if err != nil {
			return common.Cause("read player name: ", err)
		}
		if buffer.Len() == 16 { // UUID exists
			copy(metadata.Minecraft.UUID[:], buffer.Bytes())
		}
	}

	metadata.Minecraft.SniffPosition = conn.CurrentPosition()
	return nil
}
