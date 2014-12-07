package positioning

import (
	."math"
	"testing"
)

func TestSolve2d(t *testing.T) {
	pranges1, pranges2 := 460.0, 400.0
	ps  := [][]float64{{0.0, 0.0}, {500.0, 0.0}}
	rec := [][]float64{{0.0, 0.0}, {0.0, 0.0}}
	expected_rec := [][]float64{{301.6, 347.329}, {301.6, -347.329}} 
	
	Solve2d(rec, ps, pranges1, pranges2)
	err := Pow(rec[0][0] - expected_rec[0][0], 2) +
		Pow(rec[0][1] - expected_rec[0][1], 2)

	if err >= 1 {
		t.Errorf("Solve2d failed with err: %f", err)
	}
}
