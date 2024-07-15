package network

import "time"

type InboundSocketOptions struct {
	KeepAlivePeriod time.Duration `json:",omitempty"`
	Mark            int           `json:",omitempty"`
	TCPCongestion   string        `json:",omitempty"`
	TCPFastOpen     bool          `json:",omitempty"`
	MultiPathTCP    bool          `json:",omitempty"`
}

type OutboundSocketOptions struct {
	KeepAlivePeriod time.Duration `json:",omitempty"`
	Mark            int           `json:",omitempty"`
	Interface       string        `json:",omitempty"`
	TCPCongestion   string        `json:",omitempty"`
	TCPFastOpen     bool          `json:",omitempty"`
	MultiPathTCP    bool          `json:",omitempty"`
}

func ConvertLegacyOutboundOptions(inbound *InboundSocketOptions) *OutboundSocketOptions {
	if inbound == nil {
		return nil
	}
	return &OutboundSocketOptions{
		KeepAlivePeriod: inbound.KeepAlivePeriod,
		TCPFastOpen:     inbound.TCPFastOpen,
		MultiPathTCP:    inbound.MultiPathTCP,
	}
}
