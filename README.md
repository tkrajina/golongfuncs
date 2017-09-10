[![Build Status](https://api.travis-ci.org/tkrajina/golongfuncs)](https://travis-ci.org/tkrajina/golongfuncs)

# GoLongFuncs

`golongfuncs` is a tool for searching "long" functions by various measures (and combinations of measures).

## Installation

     go get -u github.com/tkrajina/golongfuncs/cmd/golongfuncs

## Measures

This tool can calculate the following function "length" measures:

* `lines`: Number of lines **without empty lines, lines ending blocks (containing only `}`), and comments**,
* `total-lines`: Number of lines **including empty lines and comments**,
* `complexity`: Cyclomatic complexity (from [gocyclo](https://github.com/fzipp/gocyclo)),
* `max-nesting`: Max nested blocks depth (note, struct initializations are not counted),
* `total-nesting`: Total nesting (in other words, if the code is formatted properly -- every indentation tab is counted once)
* `len`: Number of characters in the function (without comments and empty lines).
* `total-len`: Number of characters in the function (with comments and empty lines).
* `comments`: Number of comments. Multiline comments are counted once,
* `comment-lines`: Number of comment lines,

In addition to those, you can combine measures. Examples:

* `complexity/lines`: Calculates average complexity per line of code.
* `total-nesting/total-lines`: Calculates average nesting (indentation) per line.
* `comment-lines/total-lines`: Calculates lines of functions per line.
* etc.

## Usage

Find longest functions:

    $ golongfuncs <path>

The 50 longest functions in the code:

    $ golongfuncs -top 50 <path>

The most complex functions:

    $ golongfuncs -type complexity <path>

The tool can output multiple measures, but the result is always ordered by the first column (in this case `complexity`):

    $ golongfuncs -type complexity,lines,len <path>

Find long functions, but calculate also their complexity, avg complexity and avg nesting:

    $ golongfuncs -type lines,complexity,complexity/lines,total-nesting/total-lines .
          ExampleVeryLongfunction golongfuncs/runner_test.go:118:1       lines=305.0    complexity=1.0    complexity/lines=0.1    total-nesting/total-lines=1.0
       ExampleVeryComplexFunction golongfuncs/runner_test.go:10:1         lines=69.0   complexity=44.0    complexity/lines=0.6    total-nesting/total-lines=6.7
                       printStats main.go:54:1                            lines=21.0    complexity=9.0    complexity/lines=0.4    total-nesting/total-lines=2.5
                             main main.go:12:1                            lines=19.0    complexity=3.0    complexity/lines=0.2    total-nesting/total-lines=1.0
                        TestLines golongfuncs/runner_test.go:476:1        lines=15.0    complexity=2.0    complexity/lines=0.1    total-nesting/total-lines=0.9
                       NewVisitor golongfuncs/runner.go:94:1              lines=15.0    complexity=3.0    complexity/lines=0.2    total-nesting/total-lines=1.0
                              Get golongfuncs/models.go:34:1              lines=15.0    complexity=7.0    complexity/lines=0.5    total-nesting/total-lines=1.7
                      TestNesting golongfuncs/runner_test.go:438:1        lines=15.0    complexity=2.0    complexity/lines=0.1    total-nesting/total-lines=0.9

You can see that `ExampleVeryLongfunction` is long (305 lines), but average complexity is low (0.1) and avg nesting is 1.0.
The `ExampleVeryComplexFunction` is shorter (69 lines) but with an average complexity (per line of code) of 0.6 and avg nesting 6.7 and that is probably a good hint that the function needs refactoring.

Find functions longer than 5 lines with avg nesting (per line of code) bigger than 5 and include total lines count (including comments and empty lines) and lines count (without comments and empty lines):

    $ golongfuncs -type total-nesting/total-lines,total-lines -treshold 5 .
            ExampleVeryComplexFunction golongfuncs/runner_test.go:10:1             total-nesting/total-lines=6.7  total-lines=108.0   lines=69.0

Find functions with longest average line length:

    $ golongfuncs -type len/lines
    $ golongfuncs -type total-len/total-lines

By default, functions shorter than 10 lines are ignored. You can change that with `-min-lines <n>`.

Tests and vendored files are also ignored, use `-include-tests` and `-include-vendor` if you want to measure them.

Arbitrary files/directories can be ignored with `-ignore "<regexp>"`. For example, if you want to ignore Golang files containing `_generated.go`: `-ignore "^.*_generated.go$"`.

# License

This tool is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)
