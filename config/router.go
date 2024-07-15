package config

type Router struct {
	DefaultOutbound string  `json:",omitempty"`
	Rules           []*Rule `json:",omitempty"`
}

type Rule struct {
	Type      string
	Parameter RawJSON
	//SubRules []Rule `json:",omitempty"`
	Rewrite  RuleRewrite      `json:",omitempty"`
	Sniff    Listable[string] `json:",omitempty"`
	Outbound string           `json:",omitempty"`
	Invert   bool             `json:",omitempty"`
}

type RuleRewrite struct {
	TargetAddress string                `json:",omitempty"`
	TargetPort    uint16                `json:",omitempty"`
	Minecraft     *ruleRewriteMinecraft `json:",omitempty"`
}

type ruleRewriteMinecraft struct {
	Hostname string `json:",omitempty"`
	Port     uint16 `json:",omitempty"`
}

type RuleDomain struct {
	Domain       Listable[string] `json:",omitempty"`
	DomainSuffix Listable[string] `json:",omitempty"`
}
