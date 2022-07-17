package gocognitive

import (
	"fmt"
	"go/ast"
	"log"
	"os"
	"reflect"

	"github.com/distroy/git-go-tool/core/filecore"
	"github.com/distroy/git-go-tool/core/iocore"
)

var (
	_debug bool
)

func SetDebug(enable bool) { _debug = enable }

// AnalyzeFileByPath builds the complexity statistics.
func AnalyzeFileByPath(filePath string) ([]*Complexity, error) {
	f := &filecore.File{
		Path: filePath,
		Name: filePath,
	}

	return AnalyzeFile(f)
}

// AnalyzeDirByPath builds the complexity statistics.
func AnalyzeDirByPath(dirPath string) ([]*Complexity, error) {
	complexites := make([]*Complexity, 0, 32)
	err := filecore.WalkFiles(dirPath, func(file *filecore.File) error {
		if !file.IsGo() {
			return nil
		}

		res, err := AnalyzeFile(file)
		complexites = append(complexites, res...)
		return err
	})

	return complexites, err
}

// AnalyzeFile builds the complexity statistics.
func AnalyzeFile(f *filecore.File) ([]*Complexity, error) {
	file, err := f.Parse()
	if err != nil {
		return nil, err
	}
	res := make([]*Complexity, 0, len(file.Decls))
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			res = append(res, AnalyzeFunction(f, fn))
		}
	}
	return res, nil
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
func AnalyzeFunction(file *filecore.File, fn *ast.FuncDecl) *Complexity {
	l := log.New(iocore.Discard(), "", 0)
	if _debug {
		l = log.New(os.Stdout, fmt.Sprintf("debug %s ", funcName(fn)),
			log.Lshortfile|log.LstdFlags)
	}
	v := visitor{
		log:  l,
		file: file,
		name: fn.Name,
	}

	f := file.MustParse()

	v.log.Printf("***** %s begin *****", v.name)

	ast.Walk(&v, fn)

	v.log.Printf("***** %s end *****", v.name)
	v.log.Print("")

	pos, end := file.Position(fn.Pos()), file.Position(fn.End())
	return &Complexity{
		PkgName:    f.Name.Name,
		FuncName:   funcName(fn),
		Filename:   pos.Filename,
		Complexity: v.complexity,
		BeginLine:  pos.Line,
		EndLine:    end.Line,
	}
}
