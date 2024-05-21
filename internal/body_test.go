package internal

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// WTF <-- this will not be counted because it's not part of the function definition

// WTF BUG NOTE
func sampleFunc() {
	// TODO[TK]
	// todo
	// FIXME
}

func TestCountCommentTags(t *testing.T) {
	params := CmdParams{
		Types:        []FuncMeasurement{MaxNesting},
		IncludeTests: true,
		Verbose:      true,
	}
	stats := Do(params, []string{"."})
	assert.NotEmpty(t, stats)

	s, found := findResultForFunc("sampleFunc", stats)
	if !assert.True(t, found) {
		return
	}

	v, err := s.Get(Todos)
	assert.Nil(t, err)
	assert.Equal(t, 5.0, v)
}

type GenericArgs[T int] struct {
	Value T
}

func GenericFunc[T any](t T, arg GenericArgs[int]) string {
	if 1 == 2 {
		return ""
	}
	// TODO
	return fmt.Sprint("a")
}

func TestGeneric(t *testing.T) {
	t.Parallel()

	params := CmdParams{
		Types:        []FuncMeasurement{MaxNesting},
		IncludeTests: true,
		Verbose:      true,
	}
	stats := Do(params, []string{"."})
	s, found := findResultForFunc("GenericFunc", stats)
	if !assert.True(t, found) {
		return
	}

	{
		v, err := s.Get(Lines)
		assert.Nil(t, err)
		assert.Equal(t, 4.0, v)
	}
	{
		v, err := s.Get(Control)
		assert.Nil(t, err)
		assert.Equal(t, 1.0, v, "%#v", s)
	}
	{
		v, err := s.Get(Todos)
		assert.Nil(t, err)
		assert.Equal(t, 1.0, v)
	}
}
