package internal

import (
	"fmt"
	"strings"

	"math"
	"regexp"
)

type CmdParams struct {
	Types         []FuncMeasurement
	Treshold      float64
	MinLines      int
	Top           int
	IncludeTests  bool
	IncludeVendor bool
	Ignore        *regexp.Regexp
}

type FunctionStats struct {
	Receiver, Name, Location string
	stats                    map[FuncMeasurement]float64
}

func newFunctionStats(name, location string) *FunctionStats {
	return &FunctionStats{
		Name:     name,
		Location: location,
		stats:    map[FuncMeasurement]float64{},
	}
}

func (fs FunctionStats) FuncWithRecv() string {
	if fs.Receiver == "" {
		return fs.Name + ""
	}
	return fmt.Sprintf("(%s) %s", fs.Receiver, fs.Name)
}

func (fs FunctionStats) Get(ty FuncMeasurement) (float64, error) {
	if strings.Index(string(ty), "/") > 0 {
		parts := strings.Split(string(ty), "/")
		if len(parts) != 2 {
			return 0, fmt.Errorf("Invalit type %s", ty)
		}
		a, b := FuncMeasurement(parts[0]), FuncMeasurement(parts[1])
		if !isValidBasicType(a) || !isValidBasicType(b) {
			return 0, fmt.Errorf("Invalit type %s", ty)
		}
		val1, val2 := fs.stats[a], fs.stats[b]

		if val2 == 0 {
			return math.NaN(), nil
		}
		return val1 / val2, nil
	} else if !isValidBasicType(ty) {
		return 0, fmt.Errorf("Invalit type %s", ty)
	}
	return fs.stats[ty], nil
}

func (fs FunctionStats) Set(ty FuncMeasurement, value float64) {
	fs.stats[ty] = value
}

func (fs FunctionStats) Incr(ty FuncMeasurement, value float64) {
	fs.stats[ty] += value
}

type FunctionStatsList struct {
	SortType FuncMeasurement
	Stats    []FunctionStats
}

func (s FunctionStatsList) Len() int      { return len(s.Stats) }
func (s FunctionStatsList) Swap(i, j int) { s.Stats[i], s.Stats[j] = s.Stats[j], s.Stats[i] }
func (s FunctionStatsList) Less(i, j int) bool {
	val1, _ := s.Stats[i].Get(s.SortType)
	val2, _ := s.Stats[j].Get(s.SortType)
	return val1 >= val2
}
