/*
 * Copyright (C) distroy
 */

package main

import (
	"flag"
	"log"

	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/git"
	"github.com/distroy/git-go-tool/core/regexpcore"
)

type Flags struct {
	Threshold float64
	Top       int
	Branch    string
	Filter    *filter.Filter
}

func parseFlags() *Flags {
	f := &Flags{
		Filter: &filter.Filter{
			Includes: regexpcore.MustNewRegExps(nil),
			Excludes: regexpcore.MustNewRegExps(regexpcore.DefaultExcludes),
		},
	}

	flag.Float64Var(&f.Threshold, "top", 10, "show the top N most complex functions only")
	flag.StringVar(&f.Branch, "branch", "", "view the changes you have in your working tree relative to the named <branch>")

	flag.Var(f.Filter.Includes, "include", "the regexp for include pathes")
	flag.Var(f.Filter.Excludes, "exclude", "the regexp for exclude pathes")

	flag.Parse()

	if f.Branch == "" {
		f.Branch = git.GetBranch()
	}
	if f.Top <= 0 {
		f.Top = 10
	}

	return f
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
