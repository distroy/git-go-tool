/*
 * Copyright (C) distroy
 */

package modeservice

const (
	ModeDefault = "diff"
	ModeAll     = "all"
)

type WalkFunc = func(file string, begin, end int)

type Mode interface {
	mustInit(c *Config)

	// if begin == 0 && end == 0, check the whole file
	IsIn(file string, begin, end int) bool

	IsGitSub(file string) bool

	Walk(fn WalkFunc)
}

type Config struct {
	Mode       string
	Branch     string
	FileFilter func(file string) bool
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
