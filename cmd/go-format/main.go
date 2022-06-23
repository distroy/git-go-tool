/*
 * Copyright (C) distroy
 */

package main

import (
	"flag"
	"os"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/goformat"
	"github.com/distroy/git-go-tool/core/regexpcore"
)

type Flags struct {
	FileLine int
	Import   bool
	Formated bool
	Package  bool

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

	flag.IntVar(&f.FileLine, "file-line", 1000, "check file line")
	flag.BoolVar(&f.Import, "import", true, "enable/disable check import")
	flag.BoolVar(&f.Formated, "formated", true, "enable/disable check file formated")
	flag.BoolVar(&f.Package, "package", true, "enable/disable check package name")

	flag.Var(f.Filter.Includes, "include", "the regexp for include pathes")
	flag.Var(f.Filter.Excludes, "exclude", "the regexp for exclude pathes")

	flag.Parse()

	f.Pathes = flag.Args()
	if len(f.Pathes) == 0 {
		f.Pathes = []string{"."}
	}

	return f
}

func buildChecker(flags *Flags) goformat.Checker {
	checkers := []goformat.Checker{
		goformat.FileLineChecker(flags.FileLine),
		goformat.PackageChecker(flags.Package),
		goformat.ImportChecker(flags.Import),
		goformat.FormatChecker(flags.Formated),
	}

	return goformat.AddChecker(checkers...)
}

func main() {
	flags := parseFlags()

	checker := buildChecker(flags)

	issues := checker.Check(&filecore.File{
		Name: "cmd/go-format/main.go",
		Path: "cmd/go-format/main.go",
	})

	goformat.NewIssueWriter(os.Stdout).WriteIssues(issues)
}
