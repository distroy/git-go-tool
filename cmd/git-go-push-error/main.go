/*
 * Copyright (C) distroy
 */

package main

import (
	"os"

	"github.com/distroy/git-go-tool/config"
	"github.com/distroy/git-go-tool/core/flagcore"
	"github.com/distroy/git-go-tool/core/ptrcore"
	"github.com/distroy/git-go-tool/obj/resultobj"
	"github.com/distroy/git-go-tool/service/configservice"
	"github.com/distroy/git-go-tool/service/resultservice"
)

type Flags struct {
	Type    string                `yaml:"-" flag:"name:type; usage:required"`
	Error   string                `yaml:"-" flag:"name:error; usage:required"`
	GitDiff *config.GitDiffConfig `yaml:"git-diff" flag:"-"`
	Push    *config.PushConfig    `yaml:"push"`
}

func parseFlags() *Flags {
	cfg := &Flags{
		Push: config.DefaultPush,
	}

	cfgTmp := &Flags{}
	flagcore.MustParse(cfgTmp)

	configservice.MustParse(cfg, cfgTmp.Type)
	return cfg
}

func main() {
	flags := parseFlags()
	if flags.Type == "" || flags.Error == "" {
		flagcore.PrintUsage()
		os.Exit(1)
	}

	push := flags.Push
	err := resultservice.Push(push.PushUrl, &resultobj.Result{
		Mode:         ptrcore.GetString(flags.GitDiff.Mode),
		Type:         flags.Type,
		ProjectUrl:   push.ProjectUrl,
		TargetBranch: push.TargetBranch,
		SourceBranch: push.SourceBranch,
		Data: &resultobj.GoBaseData{
			ExecError: flags.Error,
		},
	})
	if err != nil {
		os.Exit(1)
	}
}
