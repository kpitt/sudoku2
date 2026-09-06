package sudoku

import "testing"

func TestGrid_Initialization(t *testing.T) {
	g := &Grid{}
	if len(g.Values) != 81 || len(g.Candidates) != 81 {
		t.Errorf("Grid must have 81 cells")
	}

	g.Candidates[0] = 0b1111111110
	if g.CandidatesCount(0) != 9 {
		t.Errorf("Expected 9 candidates")
	}
}
