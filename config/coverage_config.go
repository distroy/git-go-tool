/*
 * Copyright (C) distroy
 */

package config

import "github.com/distroy/git-go-tool/core/ptrcore"

var DefaultCoverage = &CoverageConfig{
	Rate: ptrcore.NewFloat64(0.65),
	Top:  ptrcore.NewInt(10),
	File: ptrcore.NewString(""),
}

type CoverageConfig struct {
	Rate *float64 `yaml:"rate"  flag:"default:0.65; usage:the lowest coverage rate. range: [0, 1.0)"`
	Top  *int     `yaml:"top"   flag:"meta:N; default:10; usage:show the top <N> least coverage rage file only"`
	File *string  `yaml:"file"  flag:"meta:file; usage:the coverage file path, cannot be empty"`
}
