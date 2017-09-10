package internal

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func Do(params CmdParams, paths []string) []FunctionStats {
	var stats []FunctionStats
	for _, path := range paths {
		if isDir(path) {
			stats = append(stats, analyzeDir(params, path)...)
		} else {
			stats = append(stats, analyzeFile(params, path)...)
		}
	}

	l := FunctionStatsList{
		SortType: params.Types[0],
		Stats:    stats,
	}
	sort.Sort(l)

	return l.Stats
}

func analyzeFile(params CmdParams, fname string) []FunctionStats {
	stats := []FunctionStats{}

	if params.Ignore != nil && params.Ignore.MatchString(fname) {
		//fmt.Println("Ignored file", fname)
		return stats
	}

	isTest := strings.HasSuffix(fname, "_test.go")
	//fmt.Println(params.IncludeTests, isTest, fname)
	if isTest && !params.IncludeTests {
		return stats
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		PrintUsage("Error opening parsing %s: %s", fname, err.Error())
	}

	//fmt.Println("file=", "pos=", f.Pos())
	v := NewVisitor(params, f, fset, stats)
	ast.Walk(v, f)

	return v.stats
}

func analyzeDir(params CmdParams, dirname string) []FunctionStats {
	stats := []FunctionStats{}

	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if !params.IncludeVendor {
			if strings.Contains(path, "vendor") { // TODO
				return err
			}
		}
		if err == nil && !info.IsDir() && strings.HasSuffix(path, ".go") {
			stats = append(stats, analyzeFile(params, path)...)
		}
		return err
	})
	if err != nil {
		PrintUsage("Error walking through files %s", err.Error())
	}
	return stats
}

type Visitor struct {
	file     *ast.File
	contents string
	fset     *token.FileSet
	offset   int
	stats    []FunctionStats
	params   CmdParams
}

func NewVisitor(params CmdParams, file *ast.File, fset *token.FileSet, stats []FunctionStats) *Visitor {
	v := Visitor{
		file:   file,
		fset:   fset,
		offset: int(file.Pos()),
		stats:  stats,
		params: params,
	}

	f := fset.File(file.Pos())
	if f == nil {
		panic("No file found for " + f.Name())
	}

	bytes, err := ioutil.ReadFile(f.Name())
	if err != nil {
		PrintUsage("Error reading %s: %s", f.Name(), err.Error())
	}

	v.contents = string(bytes)

	return &v
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		//fmt.Printf("%#v\n", node)
		if fun, is := node.(*ast.FuncDecl); is {
			stats := newFunctionStats(fun.Name.Name, v.fset.Position(fun.Pos()).String())
			calculateLines(stats, v.offset, fun, v.contents, v.file.Comments)
			calculateComplexity(stats, fun)
			calculateNesting(stats, v.offset, fun, v.contents)
			v.stats = append(v.stats, *stats)
			//fmt.Printf("stats=%d\n", len(v.stats))
		}
	}
	return v
}
