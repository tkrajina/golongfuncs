package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func calcVarsExample(s, a2 string, b2 string) (string, float64) {
	a := "jkljkL"
	var b string
	var z, p int
	x, y, _ := 1, 2, 3
	a = "aaa"
	a += "aaa"
	if a = "b"; a != b {
	}
	println(a, b, x, y, z, p)
	return "", 0.0
}

func calcVarsExampleWithNamesOutput(s2, a2 string, b2 string) (w1, w2 string, w3 int, w4 float64) {
	a := "jkljkL"
	var b string
	var z, p int
	x, y, _ := 1, 2, 3
	x = 7
	y += 13
	println(a, b, x, y, z, p)
	return "", "", 0, 0.0
}

func TestCalcVariables(t *testing.T) {
	testCalcVariables(t, "calcVarsExample", 3.0, 2.0, 6.0, 9.0)
}

func TestCalcVariablesWithNamedOut(t *testing.T) {
	testCalcVariables(t, "calcVarsExampleWithNamesOutput", 3.0, 4.0, 6.0, 8.0)
}

func testCalcVariables(t *testing.T, funcName string, in, out, vars, assignments float64) {
	params := CmdParams{
		Types:        []FuncMeasurement{MaxNesting},
		IncludeTests: true,
	}
	stats := Do(params, []string{"."})
	assert.NotEmpty(t, stats)

	s, found := findResultForFunc(funcName, stats)
	if !assert.True(t, found) {
		return
	}

	v, err := s.Get(InputParams)
	assert.Nil(t, err)
	assert.Equal(t, in, v)

	v, err = s.Get(OutputParams)
	assert.Nil(t, err)
	assert.Equal(t, out, v)

	v, err = s.Get(Variables)
	assert.Nil(t, err)
	assert.Equal(t, vars, v)

	v, err = s.Get(Assignments)
	assert.Nil(t, err)
	assert.Equal(t, assignments, v)
}
