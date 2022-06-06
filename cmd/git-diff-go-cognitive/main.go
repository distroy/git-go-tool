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
	"sort"

	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/git"
	"github.com/distroy/git-go-tool/core/gocognitive"
	"github.com/distroy/git-go-tool/core/regexpcore"
	"github.com/distroy/git-go-tool/core/termcolor"
)

type Flags struct {
	Over   int
	Top    int
	Branch string
	Filter *filter.Filter
}

func parseFlags() *Flags {
	f := &Flags{
		Filter: &filter.Filter{
			Includes: regexpcore.MustNewRegExps(nil),
			Excludes: regexpcore.MustNewRegExps(regexpcore.DefaultExcludes),
		},
	}

	flag.IntVar(&f.Over, "over", 15, "show functions with complexity > N only")
	flag.IntVar(&f.Top, "top", 10, "show the top N most complex functions only")
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

func analyzeGitNews(branch string) git.Files {
	s, err := git.ParseNewLines(branch)
	if err != nil {
		log.Fatalf("parse the git different relative to the branch:%s. err:%v", branch, err)
	}

	return git.NewFileDifferents(s)
}

func filterComplexities(array []gocognitive.Complexity, f func(gocognitive.Complexity) bool) []gocognitive.Complexity {
	n := filter.FilterSlice(array, f)
	return array[:n]
}

func analyzeCognitive(over int, filter *filter.Filter) []gocognitive.Complexity {
	complexities, err := gocognitive.AnalyzeDirByPath(".")
	if err != nil {
		log.Fatalf("analyze cognitive complexities fail. err:%s", err)
	}

	return filterComplexities(complexities, func(c gocognitive.Complexity) bool {
		return c.Complexity > over && filter.Check(c.Filename)
	})
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	f := parseFlags()

	differents := analyzeGitNews(f.Branch)

	complexities := analyzeCognitive(f.Over, f.Filter)

	newCplxes := filterComplexities(complexities, func(c gocognitive.Complexity) bool {
		return f.Filter.Check(c.Filename) && differents.IsIn(c.Filename, c.BeginLine, c.EndLine)
	})

	if len(newCplxes) > 0 {
		printForGitNews(os.Stdout, f, newCplxes)
		os.Exit(1)
	}

	printOldOvers(os.Stdout, f, complexities)
}

func printForGitNews(w io.Writer, flags *Flags, cplxes []gocognitive.Complexity) {
	sort.Sort(gocognitive.Complexites(cplxes))
	if len(cplxes) > flags.Top {
		cplxes = cplxes[:flags.Top]
	}

	fmt.Fprint(w, termcolor.Red)
	fmt.Fprintf(w, "The cogntive complexity of these *new* functions is too high (over %d): \n", flags.Over)

	for _, v := range cplxes {
		fmt.Fprintf(w, "%s\n", v.String())
	}

	fmt.Fprint(w, termcolor.Reset)
	fmt.Fprint(w, "\n")
}

func printOldOvers(w io.Writer, flags *Flags, cplxes []gocognitive.Complexity) {
	sort.Sort(gocognitive.Complexites(cplxes))
	if len(cplxes) == 0 {
		fmt.Fprint(w, termcolor.Green)
		fmt.Fprintf(w, "there is no function's cogntive complexity over %d\n", flags.Over)
		fmt.Fprint(w, termcolor.Reset)
		return
	}

	if len(cplxes) > flags.Top {
		cplxes = cplxes[:flags.Top]
	}

	fmt.Fprint(w, termcolor.Green)
	fmt.Fprintf(w, "The cogntive complexity of these *old* functions is too high (over %d): \n", flags.Over)
	fmt.Fprint(w, termcolor.Reset)

	for _, v := range cplxes {
		fmt.Fprintf(w, "%s\n", v.String())
	}

	fmt.Fprint(w, "\n")
}
