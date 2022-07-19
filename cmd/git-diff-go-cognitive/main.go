/*
 * Copyright (C) distroy
 */

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/flagcore"
	"github.com/distroy/git-go-tool/core/gocognitive"
	"github.com/distroy/git-go-tool/core/regexpcore"
	"github.com/distroy/git-go-tool/core/termcolor"
	"github.com/distroy/git-go-tool/service/modeservice"
)

const (
	defaultBufferSize = 10240
)

type Flags struct {
	ModeConfig modeservice.Config
	Over       int `flag:"name:over; meta:N; default:15; usage:show functions with complexity > <N> only and return exit code 1 if the set is non-empty"`
	Top        int `flag:"name:top; meta:N; default:10; usage:show the top <N> most complex functions only"`
	Filter     *filter.Filter
}

func filterComplexities(array []*gocognitive.Complexity, f func(*gocognitive.Complexity) bool) []*gocognitive.Complexity {
	n := filter.FilterSlice(array, f)
	return array[:n]
}

func analyzeCognitive(over int, filter *filter.Filter) []*gocognitive.Complexity {
	files := make([]*filecore.File, 0, defaultBufferSize)
	count := 0
	filecore.MustWalkFiles(".", func(f *filecore.File) error {
		if !f.IsGo() || !filter.Check(f.Name) {
			return nil
		}

		n, err := gocognitive.GetCount(f)
		if err != nil {
			log.Fatalf("analyze file cognitive complexities fail. file:%s, err:%s", f.Name, err)
		}

		count += n
		files = append(files, f)
		return nil
	})

	complexities := make([]*gocognitive.Complexity, 0, count)
	for _, f := range files {
		res, err := gocognitive.AnalyzeFile(complexities, f)
		if err != nil {
			log.Fatalf("analyze file cognitive complexities fail. file:%s, err:%s", f.Name, err)
		}

		complexities = res
	}

	return filterComplexities(complexities, func(c *gocognitive.Complexity) bool {
		return c.Complexity > over
	})
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

	// filters := getFilters(flags)
	flags.ModeConfig.FileFilter = flags.Filter.Check
	mode := modeservice.New(&flags.ModeConfig)

	complexities := analyzeCognitive(flags.Over, flags.Filter)

	overs := filterComplexities(complexities, func(c *gocognitive.Complexity) bool {
		return c.Complexity > flags.Over && flags.Filter.Check(c.Filename)
	})

	newOvers := filterComplexities(overs, func(c *gocognitive.Complexity) bool {
		return mode.IsIn(c.Filename, c.BeginLine, c.EndLine)
	})

	if len(newOvers) > 0 {
		printForGitNews(os.Stdout, flags, newOvers)
		os.Exit(1)
	}

	printOldOvers(os.Stdout, flags, overs)
}

func printForGitNews(w io.Writer, flags *Flags, cplxes []*gocognitive.Complexity) {
	sort.Sort(gocognitive.Complexites(cplxes))
	if len(cplxes) > flags.Top {
		cplxes = cplxes[:flags.Top]
	}

	fmt.Fprint(w, termcolor.Red)
	fmt.Fprintf(w, "The cognitive complexity of these *new* functions is too high (over %d): \n", flags.Over)

	for _, v := range cplxes {
		fmt.Fprintf(w, "%s\n", v.String())
	}

	fmt.Fprint(w, termcolor.Reset)
	fmt.Fprint(w, "\n")
}

func printOldOvers(w io.Writer, flags *Flags, cplxes []*gocognitive.Complexity) {
	sort.Sort(gocognitive.Complexites(cplxes))
	if len(cplxes) == 0 {
		fmt.Fprint(w, termcolor.Green)
		fmt.Fprintf(w, "there is no function's cognitive complexity over %d\n", flags.Over)
		fmt.Fprint(w, termcolor.Reset)
		return
	}

	if len(cplxes) > flags.Top {
		cplxes = cplxes[:flags.Top]
	}

	fmt.Fprint(w, termcolor.Green)
	fmt.Fprintf(w, "The cognitive complexity of these *old* functions is too high (over %d): \n", flags.Over)
	fmt.Fprint(w, termcolor.Reset)

	for _, v := range cplxes {
		fmt.Fprintf(w, "%s\n", v.String())
	}

	fmt.Fprint(w, "\n")
}
