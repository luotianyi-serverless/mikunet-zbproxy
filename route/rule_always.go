package route

import (
	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/config"
)

type RuleAlways struct {
	config *config.Rule
}

var _ Rule = (*RuleAlways)(nil)

func (r *RuleAlways) Config() *config.Rule {
	return r.config
}

func (r *RuleAlways) Match(*adapter.Metadata) bool {
	if r.config.Invert {
		return false
	}
	return true
}
