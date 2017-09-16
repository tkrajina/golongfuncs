package internal

import (
	"go/ast"
	"go/token"
)

func calculateComplexity(stats *FunctionStats, fun *ast.FuncDecl) {
	v := complexityVisitor{}
	ast.Walk(&v, fun)
	stats.Set(Complexity, float64(v.complexity))
	stats.Set(Control, float64(v.controlFlows))
}

type complexityVisitor struct {
	complexity   int
	controlFlows int
}

func (v *complexityVisitor) Visit(n ast.Node) ast.Visitor {
	switch t := n.(type) {
	case *ast.FuncDecl, *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause:
		v.complexity++
	case *ast.BinaryExpr:
		if t.Op == token.LAND || t.Op == token.LOR {
			v.complexity++
		}
	}

	switch n := n.(type) {
	case *ast.IfStmt:
		v.controlFlows++
		if n.Else != nil {
			v.controlFlows++
		}
	case *ast.ForStmt, *ast.RangeStmt, *ast.CaseClause, *ast.CommClause, *ast.DeferStmt, *ast.SelectStmt:
		v.controlFlows++
	}
	return v
}
