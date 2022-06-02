/*
 * Copyright (C) distroy
 */

package main

import (
	"flag"

	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/regexpcore"
)

type Flags struct {
	Over   int
	Top    int
	Filter *filter.Filter
}

func paresFlags() *Flags {
	f := &Flags{
		Filter: &filter.Filter{
			Includes: regexpcore.MustNewRegExps(nil),
			Excludes: regexpcore.MustNewRegExps(regexpcore.DefaultExcludes),
		},
	}

	flag.IntVar(&f.Over, "over", 0, "show functions with complexity > N only")
	flag.IntVar(&f.Top, "top", 10, "show the top N most complex functions only")

	flag.Var(f.Filter.Includes, "include", "the regexp for include pathes")
	flag.Var(f.Filter.Excludes, "exclude", "the regexp for exclude pathes")

	flag.Parse()

	return f
}

func main() {
}
