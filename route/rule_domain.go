package route

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/domain"
	"github.com/layou233/zbproxy/v3/common/set"
	"github.com/layou233/zbproxy/v3/config"
)

type ruleDomain struct {
	matcher *domain.Matcher
	config  *config.Rule
}

func (r *ruleDomain) Config() *config.Rule {
	return r.config
}

func newDomainRule(newConfig *config.Rule, listMap map[string]set.StringSet) (ruleDomain, error) {
	var domainConfig config.RuleDomain
	err := json.Unmarshal(newConfig.Parameter, &domainConfig)
	if err != nil {
		return ruleDomain{}, common.Cause("bad domain rule parameter: ", err)
	}
	builder := domain.NewMatcherBuilder(len(domainConfig.Domain) + 2*len(domainConfig.DomainSuffix))
	for _, i := range domainConfig.Domain {
		if strings.HasPrefix(i, parameterListPrefix) {
			i = strings.TrimPrefix(i, parameterListPrefix)
			list, found := listMap[i]
			if !found {
				return ruleDomain{}, fmt.Errorf("list [%s] is not found", i)
			}
			for k := range list {
				builder.AddDomain(k)
			}
		} else {
			builder.AddDomain(i)
		}
	}
	for _, i := range domainConfig.DomainSuffix {
		if strings.HasPrefix(i, parameterListPrefix) {
			i = strings.TrimPrefix(i, parameterListPrefix)
			list, found := listMap[i]
			if !found {
				return ruleDomain{}, fmt.Errorf("list [%s] is not found", i)
			}
			for k := range list {
				builder.AddDomainSuffix(k)
			}
		} else {
			builder.AddDomainSuffix(i)
		}
	}
	return ruleDomain{
		matcher: builder.Build(),
		config:  newConfig,
	}, nil
}

type RuleMinecraftHostname struct {
	ruleDomain
}

func NewMinecraftHostnameRule(newConfig *config.Rule, listMap map[string]set.StringSet) (Rule, error) {
	domainRule, err := newDomainRule(newConfig, listMap)
	if err != nil {
		return nil, err
	}
	return &RuleMinecraftHostname{domainRule}, nil
}

var _ Rule = (*RuleMinecraftHostname)(nil)

func (r *RuleMinecraftHostname) Match(metadata *adapter.Metadata) (match bool) {
	if metadata.Minecraft != nil {
		match = r.matcher.Match(metadata.Minecraft.CleanOriginDestination())
	}
	if r.config.Invert {
		match = !match
	}
	return
}
