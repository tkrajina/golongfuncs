package internal

import (
	"flag"
	"fmt"
	"go/ast"
	"os"
	"strings"
)

func PrintUsage(msg string, params ...interface{}) {
	if len(msg) > 0 {
		fmt.Fprintf(os.Stderr, msg+"\n\n", params...)
	}
	flag.Usage()
	os.Exit(1)
}

func ParseTypes(types string) ([]FuncMeasurement, error) {
	var fs FunctionStats

	var res []FuncMeasurement
	parts := strings.Split(types, ",")
	for _, p := range parts {
		ty := FuncMeasurement(strings.TrimSpace(p))
		if _, err := fs.Get(ty); err != nil {
			return nil, err
		}
		res = append(res, ty)
	}

	return res, nil
}

func calculateLines(stats *FunctionStats, offset int, fun *ast.FuncDecl, contents string, comments []*ast.CommentGroup) {
	funcBody := contents[int(fun.Pos())-offset-1 : int(fun.End())-offset]
	withoutComments := funcBody
	onlyComments := ""
	stats.Set(Comments, float64(len(comments)))
	for i := len(comments) - 1; i >= 0; i-- {
		c := comments[i]
		if c.Pos() >= fun.Pos() && c.End() <= fun.End() {
			withoutComments = withoutComments[0:c.Pos()-fun.Pos()+1] + withoutComments[c.End()-fun.Pos()+1:]
			onlyComments = contents[int(c.Pos())-offset:int(c.End())-offset] + "\n" + onlyComments
		}
	}

	stats.Set(Len, float64(len([]rune(withoutComments))))
	stats.Set(TotalLen, float64(len([]rune(funcBody))))
	stats.Set(TotalLines, float64(countLines(funcBody)))
	stats.Set(Lines, float64(countLines(withoutComments, "", "}")))
	stats.Set(CommentLines, float64(countLines(onlyComments, "", "//", "/*", "*/")))
}

func countLines(str string, ignoreLines ...string) int {
	ignore := map[string]bool{}
	for _, il := range ignoreLines {
		ignore[strings.TrimSpace(il)] = true
	}

	count := 0
	for _, line := range strings.Split(str, "\n") {
		line = strings.TrimSpace(line)
		if _, ignored := ignore[line]; ignored {
			continue
		}
		count++
	}

	return count
}
