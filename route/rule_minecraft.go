package route

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/layou233/zbproxy/v3/adapter"
	"github.com/layou233/zbproxy/v3/common/jsonx"
	"github.com/layou233/zbproxy/v3/common/set"
	"github.com/layou233/zbproxy/v3/config"
)

type RuleMinecraftPlayerName struct {
	sets   []set.StringSet
	config *config.Rule
}

var _ Rule = (*RuleMinecraftPlayerName)(nil)

func NewMinecraftPlayerNameRule(newConfig *config.Rule, listMap map[string]set.StringSet) (Rule, error) {
	var playerList jsonx.Listable[string]
	err := json.Unmarshal(newConfig.Parameter, &playerList)
	if err != nil {
		return nil, fmt.Errorf("bad player name list %v: %w", newConfig.Parameter, err)
	}
	sets := []set.StringSet{
		{}, // new set for individual names
	}
	for _, i := range playerList {
		if strings.HasPrefix(i, parameterListPrefix) {
			i = strings.TrimPrefix(i, parameterListPrefix)
			nameSet, found := listMap[i]
			if !found {
				return nil, fmt.Errorf("list [%v] is not found", i)
			}
			sets = append(sets, nameSet)
		} else {
			sets[0].Add(i)
		}
	}
	return &RuleMinecraftPlayerName{
		sets:   sets,
		config: newConfig,
	}, nil
}

func (r RuleMinecraftPlayerName) Config() *config.Rule {
	return r.config
}

func (r RuleMinecraftPlayerName) Match(metadata *adapter.Metadata) (match bool) {
	if metadata.Minecraft != nil {
		for _, nameSet := range r.sets {
			match = nameSet.Has(metadata.Minecraft.PlayerName)
			if match {
				break
			}
		}
	}
	if r.config.Invert {
		match = !match
	}
	return
}
