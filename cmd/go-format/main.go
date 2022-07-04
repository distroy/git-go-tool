/*
 * Copyright (C) distroy
 */

package main

import (
	"log"
	"os"

	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/flagcore"
	"github.com/distroy/git-go-tool/core/goformat"
	"github.com/distroy/git-go-tool/core/regexpcore"
)

type Flags struct {
	CheckerConfig goformat.Config

	Filter *filter.Filter
	Pathes []string `flag:"args; meta:path; default:."`
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	flags := &Flags{
		Filter: &filter.Filter{
			Includes: regexpcore.MustNewRegExps(nil),
			Excludes: regexpcore.MustNewRegExps(regexpcore.DefaultExcludes),
		},
	}
	flagcore.MustParse(flags)

	checker := goformat.BuildChecker(&flags.CheckerConfig)
	writer := goformat.NewIssueWriter(os.Stdout)

	goformat.NewCache().MustWalkPathes(flags.Pathes, func(x *goformat.Context) goformat.Error {
		if !flags.Filter.Check(x.Name) {
			return nil
		}

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
