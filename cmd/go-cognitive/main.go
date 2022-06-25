package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/flagcore"
	"github.com/distroy/git-go-tool/core/gocognitive"
	"github.com/distroy/git-go-tool/core/regexpcore"
)

type Flags struct {
	Over  int  `flag:"name:over; meta:N; usage:show functions with complexity > <N> only and return exit code 1 if the set is non-empty"`
	Top   int  `flag:"name:top; meta:N; usage:show the top <N> most complex functions only"`
	Avg   bool `flag:"usage:show the average complexity over all functions, not depending on whether -over or -top are set"`
	Debug bool `flag:"usage:print debug log"`

	Filter *filter.Filter
	Pathes []string `flag:"args; meta:path; default:."`
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// log.SetPrefix("go-cognitive: ")

	flags := &Flags{
		Filter: &filter.Filter{
			Includes: regexpcore.MustNewRegExps(nil),
			Excludes: regexpcore.MustNewRegExps(regexpcore.DefaultExcludes),
		},
	}

	flagcore.MustParse(flags)
	// log.Printf(" === %#v", flags)

	gocognitive.SetDebug(flags.Debug)

	res := analyzePathes(flags.Pathes, flags.Filter)
	// log.Printf(" === %#v", res)

	out := os.Stdout

	sort.Sort(gocognitive.Complexites(res))
	written := writeResult(out, res, flags)

	if flags.Avg {
		showAverage(out, res)
	}

	if flags.Over > 0 && written > 0 {
		os.Exit(1)
	}
}

func analyzePathes(pathes []string, filter *filter.Filter) []gocognitive.Complexity {
	var res []gocognitive.Complexity
	for _, path := range pathes {
		if filecore.IsDir(path) {
			res = analyzeDir(path, filter, res)
		} else {
			res = analyzeFile(path, filter, res)
		}
	}
	return res
}

func analyzeFile(filePath string, filter *filter.Filter, res []gocognitive.Complexity) []gocognitive.Complexity {
	if !strings.HasSuffix(filePath, ".go") {
		return res
	}
	if !filter.Check(filePath) {
		return res
	}

	r, err := gocognitive.AnalyzeFileByPath(filePath)
	if err != nil {
		log.Fatalf("analyze file fail. err:%s", err)
	}

	res = append(res, r...)
	return res
}

func analyzeDir(dirPath string, filter *filter.Filter, res []gocognitive.Complexity) []gocognitive.Complexity {
	if !filter.Check(dirPath) {
		return res
	}

	tmpRes, err := gocognitive.AnalyzeDirByPath(dirPath)
	if err != nil {
		log.Fatalf("analyze directory fail. err:%s", err)
	}

	for _, v := range tmpRes {
		if !filter.Check(v.Filename) {
			continue
		}
		res = append(res, v)
	}

	return res
}

func writeResult(w io.Writer, res []gocognitive.Complexity, flags *Flags) int {
	top := flags.Top
	over := flags.Over
	if top <= 0 {
		top = math.MaxInt32
	}

	for i, stat := range res {
		if i >= top {
			return i
		}
		if stat.Complexity <= over {
			return i
		}
		fmt.Fprintln(w, stat)
	}
	return len(res)
}

func showAverage(w io.Writer, cplxes []gocognitive.Complexity) {
	fmt.Fprintf(w, "Average: %.3g\n", average(cplxes))
}

func average(arr []gocognitive.Complexity) float64 {
	total := 0
	for _, s := range arr {
		total += s.Complexity
	}
	return float64(total) / float64(len(arr))
}
