package internal

import (
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
