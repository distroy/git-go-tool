/*
 * Copyright (C) distroy
 */

package main

import (
	"flag"
	"log"

	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/git"
	"github.com/distroy/git-go-tool/core/gocoverage"
	"github.com/distroy/git-go-tool/core/regexpcore"
)

type Flags struct {
	Rate   float64
	Top    int
	Branch string
	File   string
	Filter *filter.Filter
}

func parseFlags() *Flags {
	f := &Flags{
		Filter: &filter.Filter{
			Includes: regexpcore.MustNewRegExps(nil),
			Excludes: regexpcore.MustNewRegExps(regexpcore.DefaultExcludes),
		},
	}

	flag.Float64Var(&f.Rate, "rate", 0.65, "the lowest coverage rate")
	flag.IntVar(&f.Top, "top", 10, "show the top N most complex functions only")
	flag.StringVar(&f.Branch, "branch", "", "view the changes you have in your working tree relative to the named <branch>")
	flag.StringVar(&f.File, "file", "", "the coverage file path, cannot be empty")

	flag.Var(f.Filter.Includes, "include", "the regexp for include pathes")
	flag.Var(f.Filter.Excludes, "exclude", "the regexp for exclude pathes")

	flag.Parse()

	if f.File == "" {
		flag.Usage()
		log.Fatalf("-file must be empty")
	}
	if f.Branch == "" {
		f.Branch = git.GetBranch()
	}
	if f.Top <= 0 {
		f.Top = 10
	}

	return f
}

func analyzeGitNews(branch string) git.Files {
	s, err := git.ParseNewLines(branch)
	if err != nil {
		log.Fatalf("parse the git different relative to the branch:%s. err:%v", branch, err)
	}

	return git.NewFileDifferents(s)
}

func analyzeCoverages(file string, filters ...func(file string, begin, end int) bool) gocoverage.Files {
	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	f := parseFlags()

	differents := analyzeGitNews(f.Branch)

	filters := make([]func(file string, begin, end int) bool, 0, 2)
	filters = append(filters, func(file string, begin, end int) bool {
		return f.Filter.Check(file)
	})
	filters = append(filters, func(file string, begin, end int) bool {
		return differents.IsIn(file, begin, end)
	})

	// coverages := analyzeCoverages(f.File, filters...)
}
