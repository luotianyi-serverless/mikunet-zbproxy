package protocol

import (
	"context"
	"errors"
	"net"
	"os"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common/network"
	"github.com/layou233/zbproxy/v3/common/network/socks"
	"github.com/layou233/zbproxy/v3/config"
	"github.com/layou233/zbproxy/v3/protocol/minecraft"

	"github.com/phuslu/log"
)

func NewOutbound(logger *log.Logger, newConfig *config.Outbound) (adapter.Outbound, error) {
	if newConfig == nil {
		return nil, os.ErrInvalid
	}
	switch {
	case newConfig.Minecraft != nil:
		return minecraft.NewOutbound(logger, newConfig)
	}
	return &Plain{
		logger: logger,
		config: newConfig,
	}, nil
}

type Plain struct {
	logger *log.Logger
	config *config.Outbound
	router adapter.Router
	dialer network.Dialer
}

var (
	_ adapter.Outbound = (*Plain)(nil)
	_ network.Dialer   = (*Plain)(nil)
)

func (o *Plain) Name() string {
	if o.config != nil {
		return o.config.Name
	}
	return ""
}

func (o *Plain) PostInitialize(router adapter.Router) error {
	var err error
	if o.config.Dialer != "" {
		if o.config.SocketOptions != nil {
			return errors.New("socket options are not available when dialer is specified")
		}
		o.dialer, err = router.FindOutboundByName(o.config.Dialer)
		if err != nil {
			return err
		}
	} else {
		o.dialer = network.NewSystemDialer(o.config.SocketOptions)
	}
	switch o.config.ProxyOptions.Type {
	case "socks", "socks5", "socks4a", "socks4":
		o.dialer = &socks.Client{
			Dialer:  o.dialer,
			Version: o.config.ProxyOptions.Type,
			Network: o.config.ProxyOptions.Network,
			Address: o.config.ProxyOptions.Address,
		}
	}
	o.router = router
	return nil
}

func (o *Plain) Reload(newConfig *config.Outbound) error {
	o.config = newConfig
	return o.PostInitialize(o.router)
}

func (o *Plain) DialContext(ctx context.Context, network string, address string) (net.Conn, error) {
	return o.dialer.DialContext(ctx, network, address)
}
