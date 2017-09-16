package internal

import (
	"go/ast"
	"go/token"
)

// calcFunc returns the number of input params, output params and variables defined in the function
func calculateVariables(stats *FunctionStats, fun *ast.FuncDecl) {
	v := &varCounterVisitor{
		assignIdents:       map[token.Pos]bool{},
		defineAssignIdents: map[token.Pos]bool{},
	}
	ast.Walk(v, fun)
	stats.Set(InputParams, float64(v.inParams))
	stats.Set(OutputParams, float64(v.outParams))
	stats.Set(Variables, float64(v.variables))
	stats.Set(Assignments, float64(v.assignments))
}

type varCounterVisitor struct {
	inParams, outParams, variables, assignments int

	defineAssignIdents map[token.Pos]bool
	assignIdents       map[token.Pos]bool
}

func (v *varCounterVisitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}
	switch t := node.(type) {
	case *ast.FuncDecl:
		if t.Recv != nil && len(t.Recv.List) > 0 {
			v.variables++
		}
	case *ast.AssignStmt:
		for _, l := range t.Lhs {
			if t.Tok == token.DEFINE {
				v.defineAssignIdents[l.Pos()] = true
			} else {
				v.assignIdents[l.Pos()] = true
			}
		}
	case *ast.Ident:
		if _, is := v.defineAssignIdents[t.Pos()]; is {
			if t.Name != "_" {
				//fmt.Printf("Assigned variable %s\n", t.Name)
				v.variables++
				v.assignments++
			}
		}
		if _, is := v.assignIdents[t.Pos()]; is {
			if t.Name != "_" {
				//fmt.Printf("Assigned variable %s\n", t.Name)
				v.assignments++
			}
		}
	case *ast.ValueSpec:
		for _, n := range t.Names {
			if n.Name != "_" {
				//fmt.Printf("Variable %s\n", n.Name)
				v.variables++
				v.assignments++
			}
		}
	case *ast.FuncType:
		for _, p := range t.Params.List {
			v.inParams += len(p.Names)
		}
		if t.Results != nil {
			for _, r := range t.Results.List {
				if r.Names == nil {
					v.outParams++
				} else {
					v.outParams += len(r.Names)
				}
			}
		}
	}
	return v
}
