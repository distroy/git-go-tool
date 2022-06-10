/*
 * Copyright (C) distroy
 */

package modeservice

const (
	ModeAll = "all"
)

type WalkFunc = func(file string, begin, end int)

type Mode interface {
	mustInit(c *Config)

	IsIn(file string, begin, end int) bool

	Walk(fn WalkFunc)
}

type Config struct {
	Mode   string
	Branch string
}

func New(c *Config) Mode {
	var m Mode
	switch c.Mode {
	default:
		m = &modeDelta{}

	case ModeAll:
		m = &modeAll{}
	}

	if m == nil {
		m = &modeDelta{}
	}

	m.mustInit(c)
	return m
}
