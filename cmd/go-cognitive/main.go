package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/distroy/git-go-tool/core/filter"
	"github.com/distroy/git-go-tool/core/gocognitive"
	"github.com/distroy/git-go-tool/core/regexpcore"
)

const usageDoc = `Calculate cognitive complexities of Go functions.
Usage:
        go-cognitive [flags] <Go file or directory> ...
<Go file or directory>:
        default current directory
Flags:
        -over <N>   show functions with complexity > N only and
                    return exit code 1 if the set is non-empty
        -top <N>    show the top N most complex functions only
        -avg        show the average complexity over all functions,
                    not depending on whether -over or -top are set
        -include <regexp>
                    the regexp for include pathes
        -exclude <regexp>
                    the regexp for exclude pathes
                    default:
                        ^vendor/
                        /vendor/
                        \.pb\.go$

The output fields for each line are:
<complexity> <package> <function> <file:begin_row,end_row>

The document of cognitive complexity:
https://sonarsource.com/docs/CognitiveComplexity.pdf
`

type Flags struct {
	Over  int
	Top   int
	Avg   bool
	Debug bool

	Filter *filter.Filter
	Pathes []string
}

func parseFlags() *Flags {
	f := &Flags{
		Filter: &filter.Filter{
			Includes: regexpcore.MustNewRegExps(nil),
			Excludes: regexpcore.MustNewRegExps(regexpcore.DefaultExcludes),
		},
	}

	flag.IntVar(&f.Over, "over", 0, "show functions with complexity > N only")
	flag.IntVar(&f.Top, "top", -1, "show the top N most complex functions only")
	flag.BoolVar(&f.Avg, "avg", false, "show the average complexity")

	flag.BoolVar(&f.Debug, "debug", false, "show the debug message")

	flag.Var(f.Filter.Includes, "include", "the regexp for include pathes")
	flag.Var(f.Filter.Excludes, "exclude", "the regexp for exclude pathes")

	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usageDoc)
		os.Exit(2)
	}

	flag.Parse()

	f.Pathes = flag.Args()
	if len(f.Pathes) == 0 {
		f.Pathes = []string{"."}
	}

	return f
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// log.SetPrefix("go-cognitive: ")

	f := parseFlags()

	gocognitive.SetDebug(f.Debug)

	res := analyzePathes(f.Pathes, f.Filter)

	writeResult(os.Stdout, res, f)

	if f.Avg {
		showAverage(res)
	}

	if f.Over > 0 && len(res) > f.Over {
		os.Exit(1)
	}
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func analyzePathes(pathes []string, filter *filter.Filter) []gocognitive.Complexity {
	var res []gocognitive.Complexity
	for _, path := range pathes {
		if isDir(path) {
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

func writeResult(w io.Writer, res []gocognitive.Complexity, flags *Flags) {
	top := flags.Top
	over := flags.Over
	if top < 0 {
		top = math.MaxInt64
	}

	sort.Sort(complexites(res))

	for i, stat := range res {
		if i >= top {
			break
		}
		if stat.Complexity <= over {
			break
		}
		fmt.Fprintln(w, stat)
	}
}

func showAverage(cplxes []gocognitive.Complexity) {
	fmt.Printf("Average: %.3g\n", average(cplxes))
}

func average(arr []gocognitive.Complexity) float64 {
	total := 0
	for _, s := range arr {
		total += s.Complexity
	}
	return float64(total) / float64(len(arr))
}

type complexites []gocognitive.Complexity

func (s complexites) Len() int      { return len(s) }
func (s complexites) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s complexites) Less(i, j int) bool {
	a, b := s[i], s[j]
	if a.Complexity != b.Complexity {
		return a.Complexity > b.Complexity
	}
	if a.Filename != b.Filename {
		return a.Filename < b.Filename
	}
	return a.BeginLine <= b.BeginLine
}
