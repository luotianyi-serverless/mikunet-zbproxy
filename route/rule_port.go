package route

import (
	"encoding/json"
	"fmt"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/config"
)

type RuleSourcePort struct {
	config *config.Rule
	ports  map[uint16]struct{}
}

var _ Rule = (*RuleSourcePort)(nil)

func NewSourcePortRule(config *config.Rule) (*RuleSourcePort, error) {
	var ports []uint16
	err := json.Unmarshal(config.Parameter, &ports)
	if err != nil {
		return nil, fmt.Errorf("bad port list [%v]: %w", config.Parameter, err)
	}
	portsMap := make(map[uint16]struct{}, len(ports))
	for _, port := range ports {
		portsMap[port] = struct{}{}
	}
	return &RuleSourcePort{
		config: config,
		ports:  portsMap,
	}, nil
}

func (r *RuleSourcePort) Config() *config.Rule {
	return r.config
}

func (r *RuleSourcePort) Match(metadata *adapter.Metadata) (match bool) {
	_, match = r.ports[metadata.SourceAddress.Port()]
	if r.config.Invert {
		match = !match
	}
	return
}
