/*
 * Copyright (C) distroy
 */

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/git"
	"github.com/distroy/git-go-tool/core/gocoverage"
	"github.com/distroy/git-go-tool/core/regexpcore"
	"github.com/distroy/git-go-tool/core/termcolor"
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

	return f
}

func analyzeGitNews(branch string) git.Files {
	s, err := git.ParseNewLines(branch)
	if err != nil {
		log.Fatalf("parse the git different relative to the branch:%s. err:%v", branch, err)
	}

	return git.NewFileDifferents(s)
}

func analyzeCoverages(file string, filters ...func(file string, lineNo int) bool) gocoverage.Files {
	coverages, err := gocoverage.ParseFile(file)
	if err != nil {
		log.Fatalf("parse coverage file fail. file:%s, err:%v", file, err)
	}

	return gocoverage.NewFileCoverages(coverages, filters...)
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	f := parseFlags()

	differents := analyzeGitNews(f.Branch)

	filters := make([]func(file string, lineNo int) bool, 0, 1)
	filters = append(filters, func(file string, lineNo int) bool {
		return f.Filter.Check(file) && differents.IsIn(file, lineNo, lineNo)
	})

	coverages := analyzeCoverages(f.File, filters...)

	printResult(os.Stderr, f, coverages)
}

func printResult(w io.Writer, flags *Flags, coverages gocoverage.Files) {
	count := coverages.GetCount()
	if count.IsZero() {
		log.Printf("coverage rate: -, coverages:0, non coverages:0")
		return
	}

	rate := count.GetRate()
	if rate >= flags.Rate {
		log.Printf("coverage rate: %.04g, coverages:%d, non coverages:%d",
			rate, count.Coverages, count.NonCoverages)
		return
	}

	files := getTopNonCoverageFiles(coverages, flags.Top)

	log.Printf("%smust improve coverage rate. rate:%.04g, threshold:%g coverages:%d, non coverages:%d%s",
		termcolor.Red, rate, flags.Rate, count.Coverages, count.NonCoverages, termcolor.Reset)

	fmt.Fprint(w, termcolor.Red)
	fmt.Fprint(w, "\n")
	if top := flags.Top; top > 0 {
		fmt.Fprintf(w, "Non coverage files (top %d):\n", top)
	} else {
		fmt.Fprintf(w, "Non coverage files (all):\n")
	}

	for _, f := range files {
		fmt.Fprintf(w, "%s:\n", f.File.Filename)
		fmt.Fprintf(w, "coverages: %d, non coverages: %d, coverage rate: %.04g\n",
			f.Count.Coverages, f.Count.NonCoverages, f.Count.GetRate())
		fmt.Fprintf(w, "non coverage lines:\n")
		fmt.Fprint(w, f.GetNonCoverageLineNos())
		fmt.Fprintf(w, "\n")
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprint(w, termcolor.Reset)
}
