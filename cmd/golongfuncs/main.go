package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/tkrajina/golongfuncs/internal"
)

func main() {
	ty := make([]string, len(internal.AllTypes))
	for n := range internal.AllTypes {
		ty[n] = string(internal.AllTypes[n])
	}

	var ignoreFilesRegexp, ignoreFuncsRegexp, types string

	args := []string{}
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "+") {
			types += "," + arg[1:]
		} else {
			args = append(args, arg)
		}
	}
	types = strings.Trim(types, ",")
	os.Args = args

	var params internal.CmdParams
	flag.StringVar(&types, "type", types, "Type of stats, valid types are: "+strings.Join(ty, ", "))
	flag.Float64Var(&params.Threshold, "threshold", 0, "Min value, functions with value less than this will be ignored")
	flag.IntVar(&params.MinLines, "min-lines", 10, "Functions shorter than this will be ignored")
	flag.IntVar(&params.Top, "top", 25, "Show only top n functions")
	flag.BoolVar(&params.IncludeTests, "include-tests", false, "Include tests")
	flag.BoolVar(&params.IncludeVendor, "include-vendor", false, "Include vendored files")
	flag.StringVar(&ignoreFilesRegexp, "ignore", "", "Regexp for files/directories to ignore")
	flag.StringVar(&ignoreFuncsRegexp, "ignore-func", "", "Regexp for functions to ignore")
	flag.BoolVar(&params.Verbose, "verbose", false, "Verbose")
	flag.Parse()

	paths := flag.Args()

	if len(paths) == 0 {
		paths = append(paths, "./...")
	}
	if len(types) == 0 {
		types = fmt.Sprintf("%s,%s,%s", internal.Lines, internal.Complexity, internal.MaxNesting)
	}

	prepareParams(&params, types, ignoreFilesRegexp, ignoreFuncsRegexp)
	stats := internal.Do(params, paths)
	printStats(params, stats)
}

func prepareParams(params *internal.CmdParams, types, ignoreFilesRegexp, ignoreFuncsRegexp string) {
	var err error
	params.Types, err = internal.ParseTypes(types)
	if err != nil {
		internal.PrintUsage("Invalid type(s) '%s'", types)
	}
	if len(ignoreFilesRegexp) > 0 {
		r, err := regexp.Compile(ignoreFilesRegexp)
		if err != nil {
			internal.PrintUsage("Invalid ignore regexp '%s'", ignoreFilesRegexp)
		}
		params.Ignore = r
	}
	if len(ignoreFuncsRegexp) > 0 {
		r, err := regexp.Compile(ignoreFuncsRegexp)
		if err != nil {
			internal.PrintUsage("Invalid ignore regexp '%s'", ignoreFuncsRegexp)
		}
		params.IgnoreFuncs = r
	}
}

func printStats(params internal.CmdParams, stats []internal.FunctionStats) {
	count := 0
	for _, st := range stats {
		val, err := st.Get(params.Types[0])
		if err != nil {
			internal.PrintUsage("Invalid type %s\n", params.Types[0])
		}
		lines, _ := st.Get(internal.Lines)
		if val >= params.Threshold && int(lines) >= params.MinLines {
			fmt.Printf("%40s %-40s", shortenTo(st.FuncWithRecv(), 40), shortenTo(st.Location, 40))
			printSingleStat(params.Types[0], val)
			count += 1
			if len(params.Types) > 1 {
				for i := 1; i < len(params.Types); i++ {
					val, _ := st.Get(params.Types[i])
					printSingleStat(params.Types[i], val)
				}
			}
			fmt.Println()
		}
		if count >= params.Top {
			return
		}
	}
}

func shortenTo(str string, l int) string {
	if len(str) > l {
		return "..." + str[len(str)-l+5:]
	}
	return str
}

func printSingleStat(ty internal.FuncMeasurement, val float64) {
	format := fmt.Sprintf("%%%ds", len(string(ty))+8)
	fmt.Printf(format, fmt.Sprintf("%s=%.1f", ty, val))
}
