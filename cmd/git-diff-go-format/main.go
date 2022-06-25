/*
 * Copyright (C) distroy
 */

package main

import (
	"log"
	"os"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/flagcore"
	"github.com/distroy/git-go-tool/core/goformat"
	"github.com/distroy/git-go-tool/core/regexpcore"
	"github.com/distroy/git-go-tool/service/modeservice"
)

type Flags struct {
	ModeConfig    modeservice.Config
	CheckerConfig goformat.Config

	Filter *filter.Filter
	// Pathes []string `flag:"args; meta:path; default:."`
}

func main() {
	flags := &Flags{
		Filter: &filter.Filter{
			Includes: regexpcore.MustNewRegExps(nil),
			Excludes: regexpcore.MustNewRegExps(regexpcore.DefaultExcludes),
		},
	}
	flagcore.MustParse(flags)

	mode := modeservice.New(&flags.ModeConfig)

	checker := goformat.BuildChecker(&flags.CheckerConfig)
	writer := goformat.NewIssueWriter(os.Stdout)

	filecore.MustWalkFiles(".", func(f *filecore.File) error {
		if !f.IsGo() || !flags.Filter.Check(f.Name) {
			return nil
		}

		x := goformat.NewContext(f)
		if err := checker.Check(x); err != nil {
			log.Fatalf("check file format fail. file:%s, err:%v", x.Name, err)
		}

		issues := x.Issues()
		n := filter.FilterSlice(issues, func(issue *goformat.Issue) bool {
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
