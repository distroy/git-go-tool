/*
 * Copyright (C) distroy
 */

package config

import "github.com/distroy/git-go-tool/core/ptrcore"

var DefaultGoCognitive = &GoCognitiveConfig{
	Over: ptrcore.NewInt(15),
	Top:  ptrcore.NewInt(10),
}

type GoCognitiveConfig struct {
	Over *int `yaml:"over"  flag:"name:over; meta:N; default:15; usage:show functions with complexity <N> only and return exit code 1 if the set is non-empty"`
	Top  *int `yaml:"top"   flag:"name:top; meta:N; default:10; usage:show the top <N> most complex functions only"`
}
