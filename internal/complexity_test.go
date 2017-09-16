package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func controlStmts() {
	defer fmt.Println("OK")
	if 1 == 2 {
		for false {
			println(1)
		}
	}
	if 2 == 3 {
	} else {
		println("jok")
	}
	var a string
	switch a {
	case "a":
	case "b":
	case "x":
	default:
	}
	// for, if, else, switch and defer - A Tour of Go
}

func TestControlFlows(t *testing.T) {
	params := CmdParams{
		Types:        []FuncMeasurement{MaxNesting},
		IncludeTests: true,
		Verbose:      true,
	}
	stats := Do(params, []string{"."})
	assert.NotEmpty(t, stats)

	s, found := findResultForFunc("controlStmts", stats)
	if !assert.True(t, found) {
		return
	}

	v, err := s.Get(Control)
	assert.Nil(t, err)
	assert.Equal(t, 9.0, v)
}
