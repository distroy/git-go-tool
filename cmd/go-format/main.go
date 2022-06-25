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
)

type Flags struct {
	CheckerConfig goformat.Config

	Filter *filter.Filter
	Pathes []string `flag:"args"`
}

func parseFlags() *Flags {
	f := &Flags{
		Filter: &filter.Filter{
			Includes: regexpcore.MustNewRegExps(nil),
			Excludes: regexpcore.MustNewRegExps(regexpcore.DefaultExcludes),
		},
	}

	flagcore.Parse(f)
	if len(f.Pathes) == 0 {
		f.Pathes = []string{"."}
	}

	// log.Printf(" === %#v", f)
	return f
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	flags := parseFlags()

	checker := goformat.BuildChecker(&flags.CheckerConfig)
	writer := goformat.NewIssueWriter(os.Stdout)

	filecore.MustWalkPathes(flags.Pathes, func(f *filecore.File) error {
		if !f.IsGo() || !flags.Filter.Check(f.Name) {
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
