/*
 * Copyright (C) distroy
 */

package gocognitive

import (
	"bytes"
	"go/ast"
	"go/token"
	"log"

	"github.com/distroy/git-go-tool/core/filecore"
)

type visitor struct {
	log             *log.Logger
	file            *filecore.File
	name            *ast.Ident
	complexity      int
	nesting         int
	elseNodes       map[ast.Node]bool
	calculatedExprs map[ast.Expr]bool
	level           int
}

func (v *visitor) printBody(n ast.Node) {
	if v.file == nil {
		return
	}

	// 0  *ast.BlockStmt {
	// 1  .  Lbrace: cmd/git-diff-go-coverage/file.go:30:10
	// 2  .  List: []ast.Stmt (len = 1) {
	// 3  .  .  0: *ast.ExprStmt {
	// 4  .  .  .  X: *ast.CallExpr {
	// 5  .  .  .  .  Fun: *ast.SelectorExpr {
	// 6  .  .  .  .  .  X: *ast.Ident {
	// 7  .  .  .  .  .  .  NamePos: cmd/git-diff-go-coverage/file.go:31:4
	// 8  .  .  .  .  .  .  Name: "fmt"
	// 9  .  .  .  .  .  .  Obj: nil
	buffer := &bytes.Buffer{}
	v.file.WriteNode(buffer, n)
	v.log.Printf("print content: \n%s", buffer.Bytes())
}

func (v *visitor) incLevel() int {
	v.level++
	l := v.level
	return l
}

func (v *visitor) decLevel() int {
	l := v.level
	v.level--
	return l
}

func (v *visitor) incNesting() {
	v.nesting++
	v.log.Printf("nesting +1. after: %d", v.nesting)
}

func (v *visitor) decNesting() {
	v.log.Printf("nesting -1. before: %d", v.nesting)
	v.nesting--
}

func (v *visitor) incComplexity() {
	v.log.Printf("*** incComplexity +1")
	v.complexity++
}

func (v *visitor) nestIncComplexity() {
	v.log.Printf("*** nestIncComplexity +%d", v.nesting+1)
	v.complexity += (v.nesting + 1)
}

func (v *visitor) nestIncComplexityOnly() {
	v.log.Printf("*** nestIncComplexityOnly +%d", v.nesting)
	v.complexity += v.nesting
}

func (v *visitor) markAsElseNode(n ast.Node) {
	if v.elseNodes == nil {
		v.elseNodes = make(map[ast.Node]bool)
	}

	v.elseNodes[n] = true
}

func (v *visitor) markedAsElseNode(n ast.Node) bool {
	if v.elseNodes == nil {
		return false
	}

	return v.elseNodes[n]
}

func (v *visitor) markCalculated(e ast.Expr) {
	if v.calculatedExprs == nil {
		v.calculatedExprs = make(map[ast.Expr]bool)
	}

	v.calculatedExprs[e] = true
}

func (v *visitor) isCalculated(e ast.Expr) bool {
	if v.calculatedExprs == nil {
		return false
	}

	return v.calculatedExprs[e]
}

// Visit implements the ast.Visitor interface.
func (v *visitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return v
	}

	var fn func() ast.Visitor
	switch n := n.(type) {
	default:
		return v

	case *ast.IfStmt:
		fn = func() ast.Visitor { return v.visitIfStmt(n) }
	case *ast.SwitchStmt:
		fn = func() ast.Visitor { return v.visitSwitchStmt(n) }
	case *ast.TypeSwitchStmt:
		fn = func() ast.Visitor { return v.visitTypeSwitchStmt(n) }
	case *ast.SelectStmt:
		fn = func() ast.Visitor { return v.visitSelectStmt(n) }
	case *ast.ForStmt:
		fn = func() ast.Visitor { return v.visitForStmt(n) }
	case *ast.RangeStmt:
		fn = func() ast.Visitor { return v.visitRangeStmt(n) }
	case *ast.FuncLit:
		fn = func() ast.Visitor { return v.visitFuncLit(n) }
	case *ast.BranchStmt:
		fn = func() ast.Visitor { return v.visitBranchStmt(n) }
	case *ast.BinaryExpr:
		fn = func() ast.Visitor { return v.visitBinaryExpr(n) }
	case *ast.CallExpr:
		fn = func() ast.Visitor { return v.visitCallExpr(n) }
	}

	v.log.Printf("%s begin L%d", typeName(n), v.incLevel())
	defer func() {
		v.log.Printf("%s end L%d", typeName(n), v.decLevel())
	}()
	res := fn()
	return res
}

func (v *visitor) visitIfStmt(n *ast.IfStmt) ast.Visitor {
	v.incIfComplexity(n)

	if t := n.Init; t != nil {
		v.log.Printf("if init begin L%d", v.incLevel())
		ast.Walk(v, t)
		v.log.Printf("if init end L%d", v.decLevel())
	}

	v.log.Printf("if cond begin L%d", v.incLevel())
	ast.Walk(v, n.Cond)
	v.log.Printf("if cond end L%d", v.decLevel())

	v.incNesting()
	v.log.Printf("if body begin L%d", v.incLevel())
	ast.Walk(v, n.Body)
	v.log.Printf("if body end L%d", v.decLevel())
	v.decNesting()

	switch t := n.Else.(type) {
	case *ast.BlockStmt:
		v.incComplexity()

		v.log.Printf("if else block begin L%d", v.incLevel())
		// v.printBody(t)
		ast.Walk(v, t)
		v.log.Printf("if else block end L%d", v.decLevel())

	case *ast.IfStmt:
		v.markAsElseNode(t)
		v.log.Printf("if else begin L%d", v.incLevel())
		ast.Walk(v, t)
		v.log.Printf("if else end L%d", v.decLevel())
	}

	return nil
}

func (v *visitor) visitSwitchTagStmt(n *ast.SwitchStmt) ast.Visitor {
	tag := n.Tag

	v.nestIncComplexity()
	v.log.Printf("switch tag begin %s", typeName(tag))
	ast.Walk(v, tag)
	v.log.Printf("switch tag end %s", typeName(tag))

	v.incNesting()
	v.log.Printf("switch tag body begin")
	for i, tmp := range n.Body.List {
		v.log.Printf("switch tag case begin. *%d*", i)

		n, _ := tmp.(*ast.CaseClause)
		for _, n := range n.Body {
			ast.Walk(v, n)
		}

		v.log.Printf("switch tag case end. *%d*", i)
	}
	v.log.Printf("switch tag body end")
	v.decNesting()
	return nil
}

func (v *visitor) visitSwitchStmt(n *ast.SwitchStmt) ast.Visitor {
	if n := n.Init; n != nil {
		v.log.Printf("switch init begin %s", typeName(n))
		ast.Walk(v, n)
		v.log.Printf("switch init end %s", typeName(n))
	}

	if n.Tag != nil {
		return v.visitSwitchTagStmt(n)
	}

	if len(n.Body.List) == 0 {
		v.log.Printf("switch body is empty")
		return nil
	}

	v.log.Printf("switch body begin")
	for i, tmp := range n.Body.List {
		v.log.Printf("switch case begin. *%d*", i)

		if i == 0 {
			v.nestIncComplexity()
		} else {
			v.incComplexity()
		}

		n, _ := tmp.(*ast.CaseClause)
		for _, expr := range n.List {
			ast.Walk(v, expr)
		}

		v.incNesting()
		for _, n := range n.Body {
			ast.Walk(v, n)
		}
		v.decNesting()

		v.log.Printf("switch case end. *%d*", i)
	}
	v.log.Printf("switch body end")
	return nil
}

func (v *visitor) visitTypeSwitchStmt(n *ast.TypeSwitchStmt) ast.Visitor {
	v.nestIncComplexity()

	if n := n.Init; n != nil {
		v.log.Printf("switch type init begin")
		ast.Walk(v, n)
		v.log.Printf("switch type init end")
	}

	if n := n.Assign; n != nil {
		v.log.Printf("switch type assign begin %s", typeName(n))
		ast.Walk(v, n)
		v.log.Printf("switch type assign end %s", typeName(n))
	}

	v.incNesting()
	v.log.Printf("switch type body begin")
	ast.Walk(v, n.Body)
	v.log.Printf("switch type body end")
	v.decNesting()
	return nil
}

func (v *visitor) visitSelectStmt(n *ast.SelectStmt) ast.Visitor {
	v.nestIncComplexity()

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *visitor) visitForStmt(n *ast.ForStmt) ast.Visitor {
	v.nestIncComplexity()

	if n := n.Init; n != nil {
		ast.Walk(v, n)
	}

	if n := n.Cond; n != nil {
		ast.Walk(v, n)
	}

	if n := n.Post; n != nil {
		ast.Walk(v, n)
	}

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *visitor) visitRangeStmt(n *ast.RangeStmt) ast.Visitor {
	v.nestIncComplexity()

	if n := n.Key; n != nil {
		ast.Walk(v, n)
	}

	if n := n.Value; n != nil {
		ast.Walk(v, n)
	}

	ast.Walk(v, n.X)

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *visitor) visitFuncLit(n *ast.FuncLit) ast.Visitor {
	ast.Walk(v, n.Type)

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *visitor) visitBranchStmt(n *ast.BranchStmt) ast.Visitor {
	if n.Label != nil {
		v.incComplexity()
	}
	return v
}

func (v *visitor) visitBinaryExpr(n *ast.BinaryExpr) ast.Visitor {
	if v.isCalculated(n) {
		ast.Walk(v, n.X)
		ast.Walk(v, n.Y)
		return nil
	}

	// v.printBody(n)
	ops := v.collectBinaryOps(n)

	var lastOp token.Token
	cache := make([]token.Token, 0, len(ops)/3)
	for _, op := range ops {
		v.log.Printf("op: %s", op.String())
		switch op {
		default:
			// v.log.Printf("xxx op skip: %s", op.String())

		case token.LPAREN:
			// v.log.Printf("xxx op paren %d op: %s, last: %s", len(cache), op.String(), lastOp.String())
			cache = append(cache, lastOp)
			lastOp = op

		case token.RPAREN:
			lastOp = cache[len(cache)-1]
			cache = cache[:len(cache)-1]
			// v.log.Printf("xxx op paren %d op: %s, last: %s", len(cache), op.String(), lastOp.String())

		case token.LAND, token.LOR:
			// v.log.Printf("xxx op: %s, last: %s", op.String(), lastOp.String())
			if lastOp != op {
				v.incComplexity()
			}
			lastOp = op
		}
	}

	ast.Walk(v, n.X)
	ast.Walk(v, n.Y)
	return nil
}

func (v *visitor) visitCallExpr(n *ast.CallExpr) ast.Visitor {
	if callIdent, ok := n.Fun.(*ast.Ident); ok {
		obj, name := callIdent.Obj, callIdent.Name
		if obj == v.name.Obj && name == v.name.Name {
			// called by same function directly (direct recursion)
			v.incComplexity()
		}
	}
	return v
}

func (v *visitor) collectBinaryOps(exp ast.Expr) []token.Token {
	v.markCalculated(exp)
	switch exp := exp.(type) {
	case *ast.BinaryExpr:
		return v.mergeBinaryOps(v.collectBinaryOps(exp.X), exp.Op, v.collectBinaryOps(exp.Y))

	case *ast.ParenExpr:
		// interest only on what inside paranthese
		ops := v.collectBinaryOps(exp.X)
		res := make([]token.Token, 0, len(ops)+2)
		res = append(res, token.LPAREN)
		res = append(res, ops...)
		res = append(res, token.RPAREN)
		return res

	case *ast.UnaryExpr:
		return v.collectBinaryOps(exp.X)

	default:
		return []token.Token{}
	}
}

func (v *visitor) mergeBinaryOps(x []token.Token, op token.Token, y []token.Token) []token.Token {
	var res []token.Token
	if len(x) != 0 {
		res = append(res, x...)
	}
	res = append(res, op)
	if len(y) != 0 {
		res = append(res, y...)
	}
	return res
}

func (v *visitor) incIfComplexity(n *ast.IfStmt) {
	if v.markedAsElseNode(n) {
		v.incComplexity()
	} else {
		v.nestIncComplexity()
	}
}
