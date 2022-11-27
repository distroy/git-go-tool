/*
 * Copyright (C) distroy
 */

package config

import (
	"github.com/distroy/git-go-tool/core/ptrcore"
	"github.com/distroy/git-go-tool/service/modeservice"
)

var DefaultGitDiff = &GitDiffConfig{
	Mode:   ptrcore.NewString(""),
	Branch: ptrcore.NewString(""),
}

type GitDiffConfig struct {
	Mode   *string `yaml:"mode"    flag:"meta:mode; usage:compare mode: default=show the result with git diff. all=show all the result"`
	Branch *string `yaml:"branch"  flag:"meta:branch; usage:view the changes you have in your working tree relative to the named <branch>"`
}

func (c *GitDiffConfig) ToConfig(filter func(file string) bool) *modeservice.Config {
	if filter == nil {
		filter = func(file string) bool { return true }
	}

	return &modeservice.Config{
		Mode:       ptrcore.GetString(c.Mode),
		Branch:     ptrcore.GetString(c.Branch),
		FileFilter: filter,
	}
}
