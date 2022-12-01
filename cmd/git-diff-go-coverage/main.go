/*
 * Copyright (C) distroy
 */

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/distroy/git-go-tool/config"
	"github.com/distroy/git-go-tool/core/gocoverage"
	"github.com/distroy/git-go-tool/core/termcolor"
	"github.com/distroy/git-go-tool/service/configservice"
	"github.com/distroy/git-go-tool/service/modeservice"
)

type Flags struct {
	GitDiff  *config.GitDiffConfig  `yaml:"git-diff"`
	Filter   *config.FilterConfig   `yaml:",inline"`
	Coverage *config.CoverageConfig `yaml:",inline"`
}

func parseFlags() *Flags {
	cfg := &Flags{
		GitDiff:  config.DefaultGitDiff,
		Filter:   config.DefaultFilter,
		Coverage: config.DefaultCoverage,
	}

	configservice.MustParse(cfg, "go-coverage")
	return cfg
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

	flags := parseFlags()

	filter := flags.Filter.ToFilter()
	mode := modeservice.New(flags.GitDiff.ToConfig(filter.Check))

	coverages := analyzeCoverages(*flags.Coverage.File, func(file string, lineNo int) bool {
		return filter.Check(file) && mode.IsIn(file, lineNo, lineNo)
	})

	mode.Walk(func(file string, begin, end int) {
		if strings.HasSuffix(file, "_test.go") {
			return
		}
		coverages.Add(&gocoverage.Coverage{
			Filename:  file,
			BeginLine: begin,
			EndLine:   end,
			Count:     0,
		})
	})

	printResult(os.Stdout, flags, coverages)
}

func printResult(w io.Writer, flags *Flags, coverages gocoverage.Files) {
	count := coverages.GetCount()
	if count.IsZero() {
		fmt.Fprintf(w, "coverage rate: -, coverages:0, non coverages:0\n")
		return
	}

	rate := count.GetRate()
	if rate >= *flags.Coverage.Rate {
		fmt.Fprintf(w, "coverage rate: %.04g, coverages:%d, non coverages:%d\n",
			rate, count.Coverages, count.NonCoverages)
		return
	}

	files := getTopNonCoverageFiles(coverages, *flags.Coverage.Top)

	fmt.Fprintf(w, "%smust improve coverage rate. rate:%.04g, threshold:%g coverages:%d, non coverages:%d%s\n",
		termcolor.Red, rate, *flags.Coverage.Rate, count.Coverages, count.NonCoverages, termcolor.Reset)

	fmt.Fprint(w, termcolor.Red)
	fmt.Fprint(w, "\n")
	if top := *flags.Coverage.Top; top > 0 {
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
	os.Exit(1)
}
