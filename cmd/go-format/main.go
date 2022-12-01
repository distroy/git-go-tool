/*
 * Copyright (C) distroy
 */

package main

import (
	"log"
	"os"

	"github.com/distroy/git-go-tool/config"
	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/goformat"
	"github.com/distroy/git-go-tool/service/configservice"
)

type Flags struct {
	Filter   *config.FilterConfig   `yaml:",inline"`
	GoFormat *config.GoFormatConfig `yaml:",inline"`
	Pathes   []string               `yaml:"-"    flag:"args; meta:path; default:."`
}

func parseFlags() *Flags {
	cfg := &Flags{
		Filter:   config.DefaultFilter,
		GoFormat: config.DefaultGoFormat,
	}

	configservice.MustParse(cfg, "go-format")
	if len(cfg.Pathes) == 0 {
		cfg.Pathes = []string{"."}
	}
	return cfg
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	flags := parseFlags()

	filter := flags.Filter.ToFilter()
	checker := goformat.BuildChecker(flags.GoFormat.ToConfig())
	writer := goformat.NewIssueWriter(os.Stdout)

	filecore.MustWalkPathes(flags.Pathes, func(f *filecore.File) error {
		if !f.IsGo() || !filter.Check(f.Name) {
			return nil
		}

		x := goformat.NewContext(f)
		if err := checker.Check(x); err != nil {
			log.Fatalf("check file format fail. file:%s, err:%v", x.Name, err)
		}

		writer.WriteIssues(x.Issues())
		return nil
	})

	if writer.Count().Error > 0 {
		os.Exit(1)
	}
}
