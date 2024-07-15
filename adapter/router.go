package adapter

import (
	"net"

	"github.com/layou233/zbproxy/v3/common/set"
)

type Router interface {
	FindOutboundByName(name string) (Outbound, error)
	FindListsByTag(tags []string) ([]set.StringSet, error)
	HandleConnection(conn net.Conn, metadata *Metadata)
}
