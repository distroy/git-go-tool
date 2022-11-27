package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"

	"github.com/distroy/git-go-tool/config"
	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/gocognitive"
	"github.com/distroy/git-go-tool/service/configservice"
)

const (
	defaultBufferSize = 10240
)

type Flags struct {
	Filter      *config.FilterConfig      `yaml:",inline"`
	GoCognitive *config.GoCognitiveConfig `yaml:",inline"`

	Avg    bool     `yaml:"-"  flag:"usage:show the average complexity over all functions, not depending on whether -over or -top are set"`
	Debug  bool     `yaml:"-"  flag:"usage:print debug log"`
	Pathes []string `yaml:"-"  flag:"args; meta:path; default:."`
}

func parseFlags() *Flags {
	cfg := &Flags{
		Filter:      config.DefaultFilter,
		GoCognitive: config.DefaultGoCognitive,
	}

	flags := &Flags{
		Filter: config.DefaultFilter,
	}

	configservice.MustParse(cfg, flags, "go-cognitive")

	if len(cfg.Pathes) == 0 {
		cfg.Pathes = []string{"."}
	}
	return cfg
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// log.SetPrefix("go-cognitive: ")

	flags := parseFlags()
	filter := flags.Filter.ToFilter()

	gocognitive.SetDebug(flags.Debug)

	res := analyzePathes(flags.Pathes, filter)
	// log.Printf(" === %#v", res)

	out := os.Stdout

	sort.Sort(gocognitive.Complexites(res))
	written := writeResult(out, res, flags)

	if flags.Avg {
		showAverage(out, res)
	}

	if *flags.GoCognitive.Over > 0 && written > 0 {
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
	top := *flags.GoCognitive.Top
	over := *flags.GoCognitive.Over
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
