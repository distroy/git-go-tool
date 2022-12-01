/*
 * Copyright (C) distroy
 */

package config

import (
	"strings"

	"github.com/distroy/git-go-tool/core/filtercore"
	"github.com/distroy/git-go-tool/core/regexpcore"
)

var DefaultFilter = &FilterConfig{
	Includes: (*includes)(&regexpcore.DefaultIncludes),
	Excludes: (*excludes)(&regexpcore.DefaultExcludes),
}

type includes []string

func (s *includes) Default() string { return strings.Join(regexpcore.DefaultIncludes, "\n") }
func (p *includes) String() string  { return strings.Join(*p, "\n") }
func (p *includes) Set(s string) error {
	*p = append(*p, s)
	return nil
}

type excludes []string

func (s *excludes) Default() string { return strings.Join(regexpcore.DefaultExcludes, "\n") }
func (p *excludes) String() string  { return strings.Join(*p, "\n") }
func (p *excludes) Set(s string) error {
	*p = append(*p, s)
	return nil
}

type FilterConfig struct {
	Includes *includes `yaml:"include"  flag:"name:include; meta:regexp; usage:the regexp for include pathes"`
	Excludes *excludes `yaml:"exclude"  flag:"name:exclude; meta:regexp; usage:the regexp for exclude pathes"`
}

func (c *FilterConfig) ToFilter() *filtercore.Filter {
	return &filtercore.Filter{
		Includes: regexpcore.MustNewRegExps(c.unique(*c.Includes)),
		Excludes: regexpcore.MustNewRegExps(c.unique(*c.Excludes)),
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
