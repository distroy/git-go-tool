/*
 * Copyright (C) distroy
 */

package config

import (
	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/regexpcore"
)

var DefaultFilter = &FilterConfig{
	Includes: regexpcore.DefaultIncludes,
	Excludes: regexpcore.DefaultExcludes,
}

type FilterConfig struct {
	Includes []string `yaml:"include"  flag:"name:include; meta:regexp; usage:the regexp for include pathes"`
	Excludes []string `yaml:"exclude"  flag:"name:exclude; meta:regexp; usage:the regexp for exclude pathes"`
}

func (c *FilterConfig) ToFilter() *filter.Filter {
	return &filter.Filter{
		Includes: regexpcore.MustNewRegExps(c.unique(c.Includes)),
		Excludes: regexpcore.MustNewRegExps(c.unique(c.Excludes)),
	}
}

func (c *FilterConfig) unique(s []string) []string {
	l := len(s)
	m := make(map[string]struct{}, l)
	i := 0
	for j := 0; j < l; j++ {
		v := s[j]
		if _, ok := m[v]; ok {
			continue
		}

		m[v] = struct{}{}
		s[i] = v
		i++
	}

	return s[:i]
}
