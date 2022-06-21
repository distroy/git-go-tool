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
	FileLine int

	Filter *filter.Filter
	Pathes []string
}

func parseFlags() *Flags {
	f := &Flags{
		Filter: &filter.Filter{
			Includes: regexpcore.MustNewRegExps(nil),
			Excludes: regexpcore.MustNewRegExps(regexpcore.DefaultExcludes),
		},
	}

	flag.IntVar(&f.FileLine, "file-line", 1000, "file line")

	flag.Var(f.Filter.Includes, "include", "the regexp for include pathes")
	flag.Var(f.Filter.Excludes, "exclude", "the regexp for exclude pathes")

	flag.Parse()

	f.Pathes = flag.Args()
	if len(f.Pathes) == 0 {
		f.Pathes = []string{"."}
	}

	return f
}
func main() {
	// flags := parseFlags()
}
