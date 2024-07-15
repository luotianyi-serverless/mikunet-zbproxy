package adapter

import (
	"context"
	"errors"
	"net"

	"github.com/layou233/zbproxy/v3/common/bufio"
	"github.com/layou233/zbproxy/v3/config"
)

type Outbound interface {
	Name() string
	PostInitialize(router Router) error
	Reload(newConfig *config.Outbound) error
	DialContext(ctx context.Context, network string, address string) (net.Conn, error)
}

type InjectOutbound interface {
	InjectConnection(ctx context.Context, conn *bufio.CachedConn, metadata *Metadata) error
}

var ErrInjectionRequired = errors.New("injection required")
