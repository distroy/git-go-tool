package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/flagcore"
	"github.com/distroy/git-go-tool/core/gocognitive"
	"github.com/distroy/git-go-tool/core/regexpcore"
)

const (
	defaultBufferSize = 10240
)

type Flags struct {
	Over  int  `flag:"name:over; meta:N; usage:show functions with complexity <N> only and return exit code 1 if the set is non-empty"`
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

func analyzePathes(pathes []string, filter *filter.Filter) []*gocognitive.Complexity {
	files := make([]*filecore.File, 0, defaultBufferSize)
	count := 0
	filecore.MustWalkPathes(pathes, func(f *filecore.File) error {
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

	return complexities
}

func writeResult(w io.Writer, res []*gocognitive.Complexity, flags *Flags) int {
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

func showAverage(w io.Writer, cplxes []*gocognitive.Complexity) {
	fmt.Fprintf(w, "Average: %.3g\n", average(cplxes))
}

func average(arr []*gocognitive.Complexity) float64 {
	total := 0
	for _, s := range arr {
		total += s.Complexity
	}
	return float64(total) / float64(len(arr))
}
