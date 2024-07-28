package route

import (
	"encoding/json"
	"fmt"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/config"
)

type RuleSourceIPVersion struct {
	config  *config.Rule
	version uint8
}

var _ Rule = (*RuleSourceIPVersion)(nil)

func NewSourceIPVersionRule(config *config.Rule) (*RuleSourceIPVersion, error) {
	var version uint8
	err := json.Unmarshal(config.Parameter, &version)
	if err != nil {
		return nil, fmt.Errorf("bad IP version [%v]: %w", config.Parameter, err)
	}
	return &RuleSourceIPVersion{
		config:  config,
		version: version,
	}, nil
}

func (r *RuleSourceIPVersion) Config() *config.Rule {
	return r.config
}

func (r *RuleSourceIPVersion) Match(metadata *adapter.Metadata) (match bool) {
	if metadata.SourceAddress.Addr().Is4() {
		match = r.version == 4
	} else {
		match = r.version == 6
	}
	if r.config.Invert {
		match = !match
	}
	return
}
