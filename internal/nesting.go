package internal

import (
	"go/ast"
	"strings"
)

func calculateNesting(stats *FunctionStats, offset int, fun *ast.FuncDecl, contents string) {
	calculateMaxDepth(stats, fun, offset, contents)
}

func calculateMaxDepth(stats *FunctionStats, node ast.Node, offset int, contents string) {
	v := &blockNestingVisitor{
		offset:   offset,
		contents: contents,
	}
	ast.Walk(v, node)

	if v.maxNesting == 0 {
		v.maxNesting = 1
	}

	stats.Set(MaxNesting, float64(v.maxNesting))
	stats.Set(TotalNesting, float64(v.totalNesting))
}

type blockNestingVisitor struct {
	blocks       []*ast.BlockStmt
	maxNesting   int
	totalNesting int
	offset       int
	contents     string
}

func (v *blockNestingVisitor) Visit(node ast.Node) ast.Visitor {
	if v.blocks == nil {
		v.blocks = make([]*ast.BlockStmt, 0)
	}
	if node != nil {
		if b, is := node.(*ast.BlockStmt); is {
			v.calcMaxNesting(b)
			v.calcTotalNesting(b)
		}
	}
	return v
}

func (v *blockNestingVisitor) calcTotalNesting(b *ast.BlockStmt) {
	body := v.contents[int(b.Pos())-v.offset-1 : int(b.End())-v.offset]
	body = strings.TrimSpace(strings.Trim(strings.TrimSpace(body), "{}"))
	c := countLines(body)
	//fmt.Println("+", c, "for body:", body)
	//fmt.Println("----------------------------------------------------------------------------------------------------")
	v.totalNesting += c
}

func (v *blockNestingVisitor) calcMaxNesting(b *ast.BlockStmt) {
	depth := 1
	for _, previous := range v.blocks {
		if previous.Pos() < b.Pos() && b.End() < previous.End() {
			depth += 1
			if depth > v.maxNesting {
				v.maxNesting = depth
			}
		}
	}
	v.blocks = append(v.blocks, b)
}
