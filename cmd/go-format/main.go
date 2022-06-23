/*
 * Copyright (C) distroy
 */

package main

import (
	"flag"

	"log"
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
	checkers := make([]goformat.Checker, 0, 8)

	checkers = append(checkers, goformat.FileLineChecker(flags.FileLine))
	checkers = append(checkers, goformat.PackageChecker(flags.Package))
	checkers = append(checkers, goformat.ImportChecker(flags.Import))
	checkers = append(checkers, goformat.FormatChecker(flags.Formated))

	return goformat.Checkers(checkers...)
}

func walkPathes(pathes []string, fn func(f *filecore.File) error) {
	for _, path := range pathes {
		if !filecore.IsDir(path) {
			f := &filecore.File{
				Path: path,
				Name: path,
			}

			fn(f)
			continue
		}

		err := filecore.WalkFiles(path, fn)
		if err != nil {
			log.Fatalf("walk dir fail. dir:%s, err:%v", path, err)
		}
	}
}

func main() {
	flags := parseFlags()

	checker := buildChecker(flags)
	writer := goformat.NewIssueWriter(os.Stdout)

	walkPathes(flags.Pathes, func(f *filecore.File) error {
		if !f.IsGo() || !flags.Filter.Check(f.Name) {
			return nil
		}

		issues := checker.Check(f)
		writer.WriteIssues(issues)
		return nil
	})

	if writer.Count().Error > 0 {
		os.Exit(1)
	}
}
