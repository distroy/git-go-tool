/*
 * Copyright (C) distroy
 */

package config

import (
	"github.com/distroy/git-go-tool/core/ptrcore"
	"github.com/distroy/git-go-tool/service/modeservice"
)

type GitDiffConfig struct {
	Mode   *string `yaml:"mode"`
	Branch *string `yaml:"branch"`
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
