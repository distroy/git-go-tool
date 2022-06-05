package gocognitive

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/distroy/git-go-tool/core/iocore"
)

var (
	_debug bool
)

func SetDebug(enable bool) { _debug = enable }

// Complexity is statistic of the complexity.
type Complexity struct {
	PkgName    string
	FuncName   string
	Filename   string
	Complexity int
	BeginLine  int
	EndLine    int
}

func (s Complexity) String() string {
	filePos := fmt.Sprintf("%s:%d,%d", s.Filename, s.BeginLine, s.EndLine)
	return fmt.Sprintf("%d %s %s %s", s.Complexity, s.PkgName, s.FuncName, filePos)
}

// AnalyzeFileByPath builds the complexity statistics.
func AnalyzeFileByPath(filePath string) ([]Complexity, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filePath, nil, 0)
	if err != nil {
		return nil, err
	}

	return AnalyzeFile(fset, f), nil
}

// AnalyzeDirByPath builds the complexity statistics.
func AnalyzeDirByPath(dirPath string) ([]Complexity, error) {
	complexites := make([]Complexity, 0, 32)
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") {
			var res []Complexity
			res, err = AnalyzeFileByPath(path)
			complexites = append(complexites, res...)
		}
		return err
	})

	return complexites, err
}

// AnalyzeFile builds the complexity statistics.
func AnalyzeFile(fset *token.FileSet, f *ast.File) []Complexity {
	res := make([]Complexity, 0, len(f.Decls))
	for _, decl := range f.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			res = append(res, AnalyzeFunction(fset, f, fn))
		}
	}
	return res
}

// funcName returns the name representation of a function or method:
// "(Type).Name" for methods or simply "Name" for functions.
func funcName(fn *ast.FuncDecl) string {
	if fn.Recv != nil {
		if fn.Recv.NumFields() > 0 {
			typ := fn.Recv.List[0].Type
			return fmt.Sprintf("(%s).%s", recvString(typ), fn.Name)
		}
	}
	return fn.Name.Name
}

// recvString returns a string representation of recv of the
// form "T", "*T", or "BADRECV" (if not a proper receiver type).
func recvString(recv ast.Expr) string {
	switch t := recv.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + recvString(t.X)
	}
	return "BADRECV"
}

func typeName(i interface{}) string {
	return reflect.TypeOf(i).String()
}

// AnalyzeFunction calculates the cognitive complexity of a function.
func AnalyzeFunction(fset *token.FileSet, f *ast.File, fn *ast.FuncDecl) Complexity {
	l := log.New(iocore.Discard(), "", 0)
	if _debug {
		l = log.New(os.Stdout, fmt.Sprintf("debug %s ", funcName(fn)),
			log.Lshortfile|log.LstdFlags)
	}
	v := visitor{
		log:  l,
		fset: fset,
		name: fn.Name,
	}

	v.log.Printf("***** %s begin *****", v.name)

	ast.Walk(&v, fn)

	v.log.Printf("***** %s end *****", v.name)
	v.log.Print("")

	pos, end := fset.Position(fn.Pos()), fset.Position(fn.End())
	return Complexity{
		PkgName:    f.Name.Name,
		FuncName:   funcName(fn),
		Filename:   pos.Filename,
		Complexity: v.complexity,
		BeginLine:  pos.Line,
		EndLine:    end.Line,
	}
}
