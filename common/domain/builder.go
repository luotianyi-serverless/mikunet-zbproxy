package domain

import "sort"

type MatcherBuilder struct {
	seen       map[string]bool
	domainList []string
}

func NewMatcherBuilder(size int) MatcherBuilder {
	return MatcherBuilder{
		seen:       make(map[string]bool, size),
		domainList: make([]string, 0, size),
	}
}

func (b *MatcherBuilder) AddDomainSuffix(domain string) {
	if b.seen[domain] {
		return
	}
	b.seen[domain] = true
	if domain[0] == '.' {
		b.domainList = append(b.domainList, reverseDomainSuffix(domain))
	} else {
		b.domainList = append(b.domainList, reverseDomainRoot(domain))
	}
}

func (b *MatcherBuilder) AddDomain(domain string) {
	if b.seen[domain] {
		return
	}
	b.seen[domain] = true
	b.domainList = append(b.domainList, reverseDomain(domain))
}

func (b *MatcherBuilder) Build() *Matcher {
	sort.Strings(b.domainList)
	return &Matcher{newSuccinctSet(b.domainList)}
}
