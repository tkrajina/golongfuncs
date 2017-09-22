[![Build Status](https://travis-ci.org/tkrajina/golongfuncs.svg?branch=master)](https://travis-ci.org/tkrajina/golongfuncs)

# GoLongFuncs

`golongfuncs` is a tool for searching "long" functions by various measures (and combinations of measures).

This tool can help you to answer questions line:

* What are the longest functions with **total complexity** more than *x*?
* What are the longest functions with complexity **per line of code** more than *x*?
* What are functions with the biggest **number of nested blocks**?
* What are functions with the biggest number of **control flow statements**?
* What are the longest functions with **max block nesting** more than *n*?
* What are the functions with the most **variables defined**?
* What are the functions with the most **variables and assignments**?
* etc.

In other words, it will help you filter long **and complex** functions from those that are just **long**.

## Installation

     go get -u github.com/tkrajina/golongfuncs/...

## Measures

This tool can calculate the following function "length" and "complexity" measures:

* `lines`: Number of lines **without empty lines, lines ending blocks (containing only `}`), and comments**,
* `total_lines`: Number of lines **including empty lines and comments**,
* `len`: Number of characters in the function (without comments and empty lines).
* `total_len`: Number of characters in the function (with comments and empty lines).
* `comments`: Number of comments. Multiline comments are counted once,
* `comment_lines`: Number of comment lines,
* `complexity`: [Cyclomatic complexity](https://en.wikipedia.org/wiki/Cyclomatic_complexity) (from [gocyclo](https://github.com/fzipp/gocyclo)),
* `max_nesting`: Max nested blocks depth (note, struct initializations are not counted),
* `total_nesting`: Total nesting (in other words, if the code is formatted properly -- every indentation tab is counted once)
* `in_params`: The number of function input parameters
* `out_params`: The number of function output parameters
* `variables`: The number of variables in the scope of the function (without function arguments and function receivers)
* `assignments`: The number of assignments in the function (including variable declarations, `:=`, `=`, `+=`, `-=`...)
* `control`: The number of control flow statements (`if`, `else`, `switch`, `case`, `default`, `select` and `defer`)

In addition to those, you can combine measures. For example:

* `complexity/lines`: Calculates average complexity per line of code.
* `total_nesting/total_lines`: Calculates average nesting (indentation) per line.
* `comment_lines/total_lines`: Calculates lines of functions per line.
* etc.

To calculate any of those measure for your Golang code:

    $ golongfuncs -type variable
    $ golongfuncs -type total_nesting
    $ golongfuncs -type lines
    $ golongfuncs -type lines,complexity
    $ golongfuncs -type lines,complexity,complexity/lines

Calling just `golongfuncs` (without arguments) is an alias for `golongfuncs -type lines`.

## Usage examples

Find longest functions:

    $ golongfuncs <go_file>
    $ golongfuncs <directory>
    $ golongfuncs <directory>/...

If path is not specified, golongfuncs assumes it is `./...`.

Show multiple measures:

    $ golongfuncs -type lines
    $ golongfuncs -type lines,in_params
    $ golongfuncs -type lines,in_params,complexity/lines

If multiple measures are specified, the results are sorted by the first column (in this example `lines`):

By default the result shows only the top 20 results. Change that with `-top`:

    $ golongfuncs -top 50

By default, functions shorter than 10 lines are ignored. You can change that with `-min-lines <n>`.

The 100 most complex functions:

    $ golongfuncs -top 100 -type complexity ./...

The most complex functions longer than 50 lines:

    $ golongfuncs -min-lines 50 -type complexity ./...

Find long functions, but calculate also their complexity, avg complexity and avg nesting:

    $ golongfuncs -type lines,complexity,complexity/lines,total_nesting/total_lines .
          ExampleVeryLongfunction golongfuncs/runner_test.go:118:1       lines=305.0    complexity=1.0    complexity/lines=0.1    total_nesting/total_lines=1.0
       ExampleVeryComplexFunction golongfuncs/runner_test.go:10:1         lines=69.0   complexity=44.0    complexity/lines=0.6    total_nesting/total_lines=6.7
                       printStats main.go:54:1                            lines=21.0    complexity=9.0    complexity/lines=0.4    total_nesting/total_lines=2.5
                             main main.go:12:1                            lines=19.0    complexity=3.0    complexity/lines=0.2    total_nesting/total_lines=1.0
                        TestLines golongfuncs/runner_test.go:476:1        lines=15.0    complexity=2.0    complexity/lines=0.1    total_nesting/total_lines=0.9
                       NewVisitor golongfuncs/runner.go:94:1              lines=15.0    complexity=3.0    complexity/lines=0.2    total_nesting/total_lines=1.0
                              Get golongfuncs/models.go:34:1              lines=15.0    complexity=7.0    complexity/lines=0.5    total_nesting/total_lines=1.7
                      TestNesting golongfuncs/runner_test.go:438:1        lines=15.0    complexity=2.0    complexity/lines=0.1    total_nesting/total_lines=0.9

You can see that `ExampleVeryLongfunction` is long (305 lines), but average complexity is low (0.1) and avg nesting is 1.0.
Avg nesting 1.0 means that there are **no nested blocks** in this function. If half the lines were in a nested block (for example a big `if &lt;expr&gt; { ...code... }` block of code) the avg nesting would be 1.5.

The `ExampleVeryComplexFunction` is shorter (69 lines) but with an average complexity (per line of code) of 0.6 and avg nesting 6.7 and that is probably a good hint that the function needs refactoring.

Find functions longer than 5 lines with avg nesting (per line of code) bigger than 5 and include total lines count and lines count:

    $ golongfuncs -type total_nesting/total_lines,total_lines,lines -treshold 5 .
            ExampleVeryComplexFunction golongfuncs/runner_test.go:10:1             total_nesting/total_lines=6.7  total_lines=108.0   lines=69.0

Find functions with longest average line length:

    $ golongfuncs -type len/lines ./...
    $ golongfuncs -type total_len/total_lines ./...

Tests and vendored files are also ignored, use `-include-tests` and `-include-vendor` if you want to measure them.

Arbitrary files/directories can be ignored with `-ignore "<regexp>"`. For example, if you want to ignore Golang files containing `_generated.go`: `-ignore "^.*_generated.go$"`.

# License

This tool is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)
