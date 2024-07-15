package adapter

import (
	"net/netip"
	"strconv"
	"strings"

	"github.com/layou233/zbproxy/v3/common/console/color"

	"github.com/zhangyunhao116/fastrand"
)

type Protocol = uint8

/*var (
	greenPlus = color.Apply(color.FgHiGreen, "[+]")
	redMinus  = color.Apply(color.FgHiRed, "[-]")
)*/

type Metadata struct {
	ConnectionID        string
	SniffedProtocol     Protocol
	SourceIP            netip.Addr
	DestinationHostname string
	DestinationPort     uint16
	Minecraft           *MinecraftMetadata
	TLS                 *TLSMetadata
	Custom              map[string]any
}

func (m *Metadata) GenerateID() {
	id := int64(fastrand.Int31())
	idColor := fastrand.Intn(len(color.List))
	m.ConnectionID = color.Apply(color.List[idColor], "["+strconv.FormatInt(id, 10)+"]")
}

type MinecraftMetadata struct {
	ProtocolVersion      uint
	PlayerName           string
	OriginDestination    string
	RewrittenDestination string
	OriginPort           uint16
	RewrittenPort        uint16
	UUID                 [16]byte
	NextState            int8
	SniffPosition        int
}

func (m *MinecraftMetadata) IsFML() bool {
	return strings.HasSuffix(m.OriginDestination, "\x00FML\x00")
}

func (m *MinecraftMetadata) CleanOriginDestination() string {
	return strings.TrimSuffix(m.OriginDestination, "\x00FML\x00")
}

type TLSMetadata struct {
	SNI string
}
