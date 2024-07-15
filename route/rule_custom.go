package route

import (
	"github.com/layou233/zbproxy/v3/common/set"
	"github.com/layou233/zbproxy/v3/config"
)

type CustomRuleInitializer = func(config *config.Rule, listMap map[string]set.StringSet) (Rule, error)

func (r *Router) RegisterCustomRule(name string, initializer CustomRuleInitializer) error {
	if r.ruleRegistry == nil {
		r.ruleRegistry = make(map[string]CustomRuleInitializer)
	}
	if initializer == nil {
		delete(r.ruleRegistry, name)
	} else {
		r.ruleRegistry[name] = initializer
	}
	return nil
}
