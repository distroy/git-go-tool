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

type differents []git.Different

func (s differents) Len() int      { return len(s) }
func (s differents) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s differents) Less(i, j int) bool {
	a, b := s[i], s[j]
	if a.Filename != b.Filename {
		return a.Filename < b.Filename
	}
	return a.BeginLine <= b.BeginLine
}

func (s differents) toMap() map[string][]git.Different {
	sort.Sort(s)

	m := make(map[string][]git.Different)
	lastIdx := 0
	for i, v1 := range s {
		if v1.Filename == s[lastIdx].Filename {
			continue
		}

		v0 := s[lastIdx]
		m[v0.Filename] = s[lastIdx:i]
		lastIdx = i
	}

	if lastIdx < len(s) {
		v0 := s[lastIdx]
		m[v0.Filename] = s[lastIdx:]
	}

	return m
}

func analyzeGitNews(branch string) map[string][]git.Different {
	s, err := git.ParseNewLines(branch)
	if err != nil {
		log.Fatalf("parse the git different relative to the branch:%s. err:%v", branch, err)
	}

	return differents(s).toMap()
}

func filterComplexities(array []gocognitive.Complexity, f func(gocognitive.Complexity) bool) []gocognitive.Complexity {
	i := 0
	l := len(array)
	for i < l {
		var vl, vr *gocognitive.Complexity

		for i < l {
			vl = &array[i]
			if !f(*vl) {
				break
			}
			i++
		}

		for i < l {
			vr = &array[l-1]
			if f(*vr) {
				break
			}
			l--
		}

		*vl, *vr = *vr, *vl
	}
	return array[:l]
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

func isComplexityInGitNews(c gocognitive.Complexity, diffMap map[string][]git.Different) bool {
	ds := diffMap[c.Filename]
	if len(ds) == 0 {
		return false
	}

	idx := sort.Search(len(ds), func(i int) bool {
		return ds[i].BeginLine >= c.BeginLine
	})
	if idx >= len(ds) {
		return false
	}

	d := ds[idx]
	if d.BeginLine > c.EndLine || d.EndLine < c.BeginLine {
		return false
	}

	return true
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	f := parseFlags()

	differents := analyzeGitNews(f.Branch)

	complexities := analyzeCognitive(f.Over, f.Filter)

	newCplxes := filterComplexities(complexities, func(c gocognitive.Complexity) bool {
		return f.Filter.Check(c.Filename) && isComplexityInGitNews(c, differents)
	})

	if len(newCplxes) > 0 {
		printForGitNews(os.Stdout, f, newCplxes)
		os.Exit(1)
	}

	printOldOvers(os.Stdout, f, complexities)
}

func printForGitNews(w io.Writer, flags *Flags, cplxes []gocognitive.Complexity) {
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
	if len(cplxes) == 0 {
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
