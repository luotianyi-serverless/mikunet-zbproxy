package route

import (
	"encoding/json"
	"fmt"
	"net/netip"
	"strings"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common"
	"github.com/layou233/zbproxy/v3/common/jsonx"
	"github.com/layou233/zbproxy/v3/common/set"
	"github.com/layou233/zbproxy/v3/config"

	"go4.org/netipx"
)

type RuleSourceIP struct {
	set    netipx.IPSet
	config *config.Rule
}

var _ Rule = (*RuleSourceIP)(nil)

func NewSourceIPRule(newConfig *config.Rule, listMap map[string]set.StringSet) (*RuleSourceIP, error) {
	var cidrList jsonx.Listable[string]
	err := json.Unmarshal(newConfig.Parameter, &cidrList)
	if err != nil {
		return nil, fmt.Errorf("bad IP CIDR list [%v]: %w", newConfig.Parameter, err)
	}
	var builder netipx.IPSetBuilder
	for _, i := range cidrList {
		if strings.HasPrefix(i, parameterListPrefix) {
			i = strings.TrimPrefix(i, parameterListPrefix)
			list, found := listMap[i]
			if !found {
				return nil, fmt.Errorf("list [%s] is not found", i)
			}
			for k := range list {
				err = addIPToBuilder(&builder, k)
				if err != nil {
					return nil, fmt.Errorf("bad IP address or CIDR [%s]: %w", i, err)
				}
			}
		} else {
			err = addIPToBuilder(&builder, i)
			if err != nil {
				return nil, fmt.Errorf("bad IP address or CIDR [%s]: %w", i, err)
			}
		}
	}
	ipSet, err := builder.IPSet()
	if err != nil {
		return nil, common.Cause("build IP set: ", err)
	}
	return &RuleSourceIP{
		config: newConfig,
		set:    *ipSet,
	}, nil
}

func addIPToBuilder(builder *netipx.IPSetBuilder, s string) error {
	// modified from netipx.ParsePrefixOrAddr
	i := strings.LastIndexByte(s, '/')
	if i < 0 {
		addr, err := netip.ParseAddr(s)
		if err != nil {
			return err
		}
		builder.Add(addr)
	} else {
		prefix, err := netip.ParsePrefix(s)
		if err != nil {
			return err
		}
		builder.AddPrefix(prefix)
	}
	return nil
}

func (r *RuleSourceIP) Config() *config.Rule {
	return r.config
}

func (r *RuleSourceIP) Match(metadata *adapter.Metadata) (match bool) {
	match = r.set.Contains(metadata.SourceAddress.Addr().WithZone(""))
	if r.config.Invert {
		match = !match
	}
	return
}
