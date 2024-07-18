package protocol

import (
	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common/bufio"
	"github.com/layou233/zbproxy/v3/protocol/minecraft"

	"github.com/phuslu/log"
)

const (
	SniffTypeUndefined = ""
	SniffTypeAll       = "all"
	SniffTypeMinecraft = "minecraft"
	SniffTypeTLS       = "tls"
)

type SnifferFunc = func(conn bufio.PeekConn, metadata *adapter.Metadata) error

func Sniff(logger *log.Logger, conn bufio.PeekConn, metadata *adapter.Metadata, registry map[string]SnifferFunc, protocols ...string) {
	startPosition := conn.CurrentPosition()
	if startPosition < 0 {
		startPosition = 0
	}
	for _, protocol := range protocols {
		var err error
		sniffAll := protocol == SniffTypeAll
		switch protocol {
		case SniffTypeAll:
			fallthrough

		case SniffTypeMinecraft:
			if metadata.Minecraft == nil {
				err = minecraft.SniffClientHandshake(conn, metadata)
				if err != nil {
					logger.Trace().Str("protocol", protocol).Err(err).Msg("sniff error")
				}
			}
			if !sniffAll {
				break
			}
			fallthrough

		//case SniffTypeTLS: // TODO

		default:
			if sniffAll {
				for _, snifferFunc := range registry {
					err = snifferFunc(conn, metadata)
					if err != nil {
						logger.Trace().Str("protocol", protocol).Err(err).Msg("sniff error")
					}
				}
				return
			} else if len(registry) > 0 {
				if snifferFunc := registry[protocol]; snifferFunc != nil {
					err = snifferFunc(conn, metadata)
					if err != nil {
						logger.Trace().Str("protocol", protocol).Err(err).Msg("sniff error")
					}
				} else {
					logger.Fatal().Str("protocol", protocol).Msg("unsupported protocol")
				}
			} else {
				logger.Fatal().Str("protocol", protocol).Msg("unsupported protocol")
			}
		}
		conn.Rewind(startPosition)
	}
}
