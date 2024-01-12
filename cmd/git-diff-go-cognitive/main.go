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

	"github.com/distroy/git-go-tool/config"
	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/filtercore"
	"github.com/distroy/git-go-tool/core/gocognitive"
	"github.com/distroy/git-go-tool/core/ptrcore"
	"github.com/distroy/git-go-tool/core/termcolor"
	"github.com/distroy/git-go-tool/obj/resultobj"
	"github.com/distroy/git-go-tool/service/configservice"
	"github.com/distroy/git-go-tool/service/modeservice"
	"github.com/distroy/git-go-tool/service/resultservice"
)

const (
	defaultBufferSize = 10240
)

type Flags struct {
	GitDiff     *config.GitDiffConfig     `yaml:"git-diff"`
	Filter      *config.FilterConfig      `yaml:",inline"`
	GoCognitive *config.GoCognitiveConfig `yaml:",inline"`
	Push        *config.PushConfig        `yaml:"push"`
}

func parseFlags() *Flags {
	cfg := &Flags{
		GitDiff:     config.DefaultGitDiff,
		Filter:      config.DefaultFilter,
		GoCognitive: config.DefaultGoCognitive,
		Push:        config.DefaultPush,
	}

	configservice.MustParse(cfg, "go-cognitive")
	return cfg
}

func filterComplexities(array []*gocognitive.Complexity, f func(*gocognitive.Complexity) bool) []*gocognitive.Complexity {
	n := filtercore.FilterSlice(array, f)
	// log.Printf(" === %d", n)
	return array[:n]
}

func analyzeCognitive(over int, filter *filtercore.Filter) []*gocognitive.Complexity {
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

	flags := parseFlags()
	filter := flags.Filter.ToFilter()

	// filters := getFilters(flags)
	mode := modeservice.New(flags.GitDiff.ToConfig(filter.Check))

	complexities := analyzeCognitive(*flags.GoCognitive.Over, filter)

	overs := filterComplexities(complexities, func(c *gocognitive.Complexity) bool {
		return c.Complexity > *flags.GoCognitive.Over && !mode.IsGitSub(c.Filename)
	})

	// log.Printf(" === ")
	newOvers := filterComplexities(overs, func(c *gocognitive.Complexity) bool {
		// log.Printf(" === ")
		return mode.IsIn(c.Filename, c.BeginLine, c.EndLine)
	})

	if len(newOvers) > 0 {
		printForGitNews(os.Stdout, flags, newOvers)
		pushResult(flags, newOvers)
		os.Exit(1)
	}

	printOldOvers(os.Stdout, flags, overs)
	pushResult(flags, overs)
}

func printForGitNews(w io.Writer, flags *Flags, cplxes []*gocognitive.Complexity) {
	top := *flags.GoCognitive.Top
	over := *flags.GoCognitive.Over

	sort.Sort(gocognitive.Complexities(cplxes))
	if len(cplxes) > top {
		cplxes = cplxes[:top]
	}

	fmt.Fprint(w, termcolor.Red)
	fmt.Fprintf(w, "The cognitive complexity of these *new* functions is too high (over %d): \n", over)

	for _, v := range cplxes {
		fmt.Fprintf(w, "%s\n", v.String())
	}

	fmt.Fprint(w, termcolor.Reset)
	fmt.Fprint(w, "\n")
}

func printOldOvers(w io.Writer, flags *Flags, cplxes []*gocognitive.Complexity) {
	top := *flags.GoCognitive.Top
	over := *flags.GoCognitive.Over

	sort.Sort(gocognitive.Complexities(cplxes))
	if len(cplxes) == 0 {
		fmt.Fprint(w, termcolor.Green)
		fmt.Fprintf(w, "there is no function's cognitive complexity over %d\n", over)
		fmt.Fprint(w, termcolor.Reset)
		return
	}

	if len(cplxes) > top {
		cplxes = cplxes[:top]
	}

	fmt.Fprint(w, termcolor.Green)
	fmt.Fprintf(w, "The cognitive complexity of these *old* functions is too high (over %d): \n", over)
	fmt.Fprint(w, termcolor.Reset)

	for _, v := range cplxes {
		fmt.Fprintf(w, "%s\n", v.String())
	}

	fmt.Fprint(w, "\n")
}

func pushResult(flags *Flags, overs []*gocognitive.Complexity) {
	push := flags.Push
	if push == nil {
		return
	}
	resultservice.Push(push.PushUrl, &resultobj.Result{
		Mode:         ptrcore.GetString(flags.GitDiff.Mode),
		Type:         resultobj.TypeGoCognitive,
		ProjectUrl:   push.ProjectUrl,
		TargetBranch: push.TargetBranch,
		SourceBranch: push.SourceBranch,
		Data: &resultobj.GoComplexityData{
			Threshold:              ptrcore.GetInt(flags.GoCognitive.Over),
			FunctionsOverThreshold: overs,
		},
	})
}
