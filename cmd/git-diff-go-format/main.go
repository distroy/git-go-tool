/*
 * Copyright (C) distroy
 */

package main

import (
	"log"
	"os"

	"github.com/distroy/git-go-tool/config"
	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/filtercore"
	"github.com/distroy/git-go-tool/core/goformat"
	"github.com/distroy/git-go-tool/service/configservice"
	"github.com/distroy/git-go-tool/service/modeservice"
)

type Flags struct {
	GitDiff  *config.GitDiffConfig  `yaml:"git-diff"`
	Filter   *config.FilterConfig   `yaml:",inline"`
	GoFormat *config.GoFormatConfig `yaml:",inline"`

	// Pathes []string `flag:"args; meta:path; default:."`
}

func parseFlags() *Flags {
	cfg := &Flags{
		GitDiff:  config.DefaultGitDiff,
		Filter:   config.DefaultFilter,
		GoFormat: config.DefaultGoFormat,
	}

	configservice.MustParse(cfg, "go-format")
	return cfg
}

func main() {
	flags := parseFlags()

	filter := flags.Filter.ToFilter()

	mode := modeservice.New(flags.GitDiff.ToConfig(filter.Check))

	checker := goformat.BuildChecker(flags.GoFormat.ToConfig())
	writer := goformat.NewIssueWriter(os.Stdout)

	filecore.MustWalkFiles(".", func(f *filecore.File) error {
		if !f.IsGo() || !filter.Check(f.Name) {
			return nil
		}

		x := goformat.NewContext(f)
		if err := checker.Check(x); err != nil {
			log.Fatalf("check file format fail. file:%s, err:%v", x.Name, err)
		}

		issues := x.Issues()
		n := filtercore.FilterSlice(issues, func(issue *goformat.Issue) bool {
			return mode.IsIn(issue.Filename, issue.BeginLine, issue.EndLine)
		})
		issues = issues[:n]

		writer.WriteIssues(issues)
		return nil
	})

	if writer.Count().Error > 0 {
		os.Exit(1)
	}
}
