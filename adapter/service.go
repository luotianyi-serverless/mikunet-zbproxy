package adapter

import (
	"context"
	"io"

	"github.com/layou233/zbproxy/v3/config"
)

type Service interface {
	Start(ctx context.Context) error
	Reload(ctx context.Context, newConfig *config.Service) error
	UpdateRouter(router Router)
	io.Closer
}
