package route

import (
	"encoding/json"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/set"
	"github.com/layou233/zbproxy/v3/config"

	"github.com/phuslu/log"
)

type ruleLogic struct {
	rules  []Rule
	config *config.Rule
}

func newLogicalRule(logger *log.Logger, newConfig *config.Rule, listMap map[string]set.StringSet, ruleRegistry map[string]CustomRuleInitializer) (ruleLogic, error) {
	var ruleConfig []config.Rule
	err := json.Unmarshal(newConfig.Parameter, &ruleConfig)
	if err != nil {
		return ruleLogic{}, common.Cause("bad rule in logic rule parameter: ", err)
	}
	rules := make([]Rule, 0, len(ruleConfig))
	for i := range ruleConfig {
		var newRule Rule
		newRule, err = NewRule(logger, &ruleConfig[i], listMap, ruleRegistry)
		if err != nil {
			return ruleLogic{}, common.Cause("initialize rule in logic rule parameter: ", err)
		}
		rules = append(rules, newRule)
	}
	return ruleLogic{
		rules:  rules,
		config: newConfig,
	}, nil
}

func (r *ruleLogic) Config() *config.Rule {
	return r.config
}

type RuleLogicalAnd struct {
	ruleLogic
}

var _ Rule = (*RuleLogicalAnd)(nil)

func NewLogicalAndRule(logger *log.Logger, newConfig *config.Rule, listMap map[string]set.StringSet, ruleRegistry map[string]CustomRuleInitializer) (Rule, error) {
	logicRule, err := newLogicalRule(logger, newConfig, listMap, ruleRegistry)
	if err != nil {
		return nil, err
	}
	return &RuleLogicalAnd{logicRule}, nil
}

func (r RuleLogicalAnd) Match(metadata *adapter.Metadata) (match bool) {
	match = true
	for _, rule := range r.rules {
		if !rule.Match(metadata) {
			match = false
			break
		}
	}
	if r.config.Invert {
		match = !match
	}
	return
}

type RuleLogicalOr struct {
	ruleLogic
}

var _ Rule = (*RuleLogicalOr)(nil)

func NewLogicalOrRule(logger *log.Logger, newConfig *config.Rule, listMap map[string]set.StringSet, ruleRegistry map[string]CustomRuleInitializer) (Rule, error) {
	logicRule, err := newLogicalRule(logger, newConfig, listMap, ruleRegistry)
	if err != nil {
		return nil, err
	}
	return &RuleLogicalOr{logicRule}, nil
}

func (r RuleLogicalOr) Match(metadata *adapter.Metadata) (match bool) {
	for _, rule := range r.rules {
		if rule.Match(metadata) {
			match = true
			break
		}
	}
	if r.config.Invert {
		match = !match
	}
	return
}
