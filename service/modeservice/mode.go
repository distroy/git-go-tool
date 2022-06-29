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

	// if begin == 0 && end == 0, check the whole file
	IsIn(file string, begin, end int) bool

	Walk(fn WalkFunc)
}

type Config struct {
	Mode       string                 `flag:"meta:mode; usage:compare mode: default=show the result with git diff. all=show all the coverage"`
	Branch     string                 `flag:"meta:branch; usage:view the changes you have in your working tree relative to the named <branch>"`
	FileFilter func(file string) bool `flag:"-"`
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
