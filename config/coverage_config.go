/*
 * Copyright (C) distroy
 */

package config

import (
	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/regexpcore"
)

type CoverageConfig struct {
	Rate     *float64 `yaml:"rate"`
	Top      *int     `yaml:"top"`
	File     *string  `yaml:"file"`
	Includes []string `yaml:"include"`
	Excludes []string `yaml:"exclude"`
}

func (c *CoverageConfig) ToFilter() *filter.Filter {
	return &filter.Filter{
		Includes: regexpcore.MustNewRegExps(c.Includes),
		Excludes: regexpcore.MustNewRegExps(c.Excludes),
	}
}
