package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/distroy/git-go-tool/core/gocognitive"
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

func usage() {
	fmt.Fprint(os.Stderr, usageDoc)
	os.Exit(2)
}

var (
	_defaultExcludes = []*regexp.Regexp{
		regexp.MustCompile(`^vendor/`),
		regexp.MustCompile(`/vendor/`),
		regexp.MustCompile(`\.pb\.go$`),
	}

	over = flag.Int("over", 0, "show functions with complexity > N only")
	top  = flag.Int("top", -1, "show the top N most complex functions only")
	avg  = flag.Bool("avg", false, "show the average complexity")

	includes = flagRegexps("include", nil, "the regexp for include pathes")
	excludes = flagRegexps("exclude", _defaultExcludes, "the regexp for exclude pathes")

	debug = flag.Bool("debug", false, "show the debug message")
)

func flagRegexps(name string, def []*regexp.Regexp, usage string) *flagRegexpsValue {
	val := flagRegexpsValue(def)
	flag.Var(&val, name, usage)
	return &val
}

type flagRegexpsValue []*regexp.Regexp

func (p *flagRegexpsValue) Set(s string) error {
	re, err := regexp.Compile(s)
	if err == nil {
		*p = append(*p, re)
	}
	return nil
}

func (p *flagRegexpsValue) String() string { return "" }

func main() {
	// log.SetFlags(log.Flags() | log.Lshortfile)
	log.SetFlags(0)
	log.SetPrefix("go-cognitive: ")

	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	gocognitive.SetDebug(*debug)

	res := analyze(args)
	sort.Sort(complexites(res))
	written := writeStats(os.Stdout, res)

	if *avg {
		showAverage(res)
	}

	if *over > 0 && written > 0 {
		os.Exit(1)
	}
}

func isPathIgnored(path string) bool {
	for _, re := range *includes {
		loc := re.FindStringIndex(path)
		if len(loc) == 2 {
			return false
		}
	}
	for _, re := range *excludes {
		loc := re.FindStringIndex(path)
		if len(loc) == 2 {
			return true
		}
	}
	return false
}

func analyze(paths []string) []gocognitive.Complexity {
	var res []gocognitive.Complexity
	for _, path := range paths {
		if isDir(path) {
			res = analyzeDir(path, res)
		} else {
			res = analyzeFile(path, res)
		}
	}

	return res
}

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func analyzeFile(filePath string, cplxes []gocognitive.Complexity) []gocognitive.Complexity {
	if isPathIgnored(filePath) {
		return nil
	}

	res, err := gocognitive.AnalyzeFileByPath(filePath)
	if err != nil {
		log.Fatalf("analyze file fail. err:%s", err)
	}

	cplxes = append(cplxes, res...)
	return cplxes
}

func analyzeDir(dirname string, cplxes []gocognitive.Complexity) []gocognitive.Complexity {
	if isPathIgnored(dirname) {
		return cplxes
	}

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") {
			cplxes = analyzeFile(path, cplxes)
		}
		return err
	})
	if err != nil {
		log.Fatalf("average directory fail. err:%s", err)
	}

	return cplxes
}

func writeStats(w io.Writer, sortedStats []gocognitive.Complexity) int {
	for i, stat := range sortedStats {
		if i == *top {
			return i
		}
		if stat.Complexity <= *over {
			return i
		}
		fmt.Fprintln(w, stat)
	}
	return len(sortedStats)
}

func showAverage(cplxes []gocognitive.Complexity) {
	fmt.Printf("Average: %.3g\n", average(cplxes))
}

func average(cplxes []gocognitive.Complexity) float64 {
	total := 0
	for _, s := range cplxes {
		total += s.Complexity
	}
	return float64(total) / float64(len(cplxes))
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
