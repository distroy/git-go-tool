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

	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/flagcore"
	"github.com/distroy/git-go-tool/core/gocoverage"
	"github.com/distroy/git-go-tool/core/regexpcore"
	"github.com/distroy/git-go-tool/core/termcolor"
	"github.com/distroy/git-go-tool/service/modeservice"
)

type Flags struct {
	ModeConfig modeservice.Config
	Rate       float64 `flag:"default:0.65; usage:the lowest coverage rate. range: [0, 1.0)"`
	Top        int     `flag:"meta:N; default:10; usage:show the top <N> least coverage rage file only"`
	File       string  `flag:"meta:file; usage:the coverage file path, cannot be empty"`
	Filter     *filter.Filter
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

	flags := &Flags{
		Filter: &filter.Filter{
			Includes: regexpcore.MustNewRegExps(nil),
			Excludes: regexpcore.MustNewRegExps(regexpcore.DefaultExcludes),
		},
	}
	flagcore.MustParse(flags)

	flags.ModeConfig.FileFilter = flags.Filter.Check
	mode := modeservice.New(&flags.ModeConfig)

	coverages := analyzeCoverages(flags.File, func(file string, lineNo int) bool {
		return flags.Filter.Check(file) && mode.IsIn(file, lineNo, lineNo)
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
	if rate >= flags.Rate {
		fmt.Fprintf(w, "coverage rate: %.04g, coverages:%d, non coverages:%d\n",
			rate, count.Coverages, count.NonCoverages)
		return
	}

	files := getTopNonCoverageFiles(coverages, flags.Top)

	fmt.Fprintf(w, "%smust improve coverage rate. rate:%.04g, threshold:%g coverages:%d, non coverages:%d%s\n",
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
	os.Exit(1)
}
