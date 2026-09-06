package sudoku

import (
	"strings"
	"testing"
)

func TestParseAndInitialize(t *testing.T) {
	// 81 char string with one '1' and '.' for the rest
	input := "1" + strings.Repeat(".", 80)

	g, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if g.Values[0] != 1 {
		t.Errorf("Expected Values[0] to be 1, got %d", g.Values[0])
	}
	if g.Clues != 1 {
		t.Errorf("Expected Clues to be 1, got %d", g.Clues)
	}

	g.InitializeCandidates()

	// Cell 0 is a given, its candidate mask should be 0
	if g.Candidates[0] != 0 {
		t.Errorf("Expected Candidates[0] to be 0 for given cell, got %b", g.Candidates[0])
	}

	// Cell 1 is in same row (peer), shouldn't have '1' as candidate
	if g.Candidates[1]&(1<<1) != 0 {
		t.Errorf("Candidate 1 should be eliminated from peer cell 1")
	}

	// Cell 80 (bottom-right) is not a peer of cell 0 (top-left), should still have '1' as candidate
	if g.Candidates[80]&(1<<1) == 0 {
		t.Errorf("Candidate 1 should remain valid for non-peer cell 80")
	}

	// All 20 peers of cell 0 must have candidate 1 eliminated
	for _, p := range Peers[0] {
		if g.Candidates[p]&(1<<1) != 0 {
			t.Errorf("Peer %d still has candidate 1 set", p)
		}
	}
}

func TestParse_FormattedInput(t *testing.T) {
	// Grid formatted with borders, whitespace, and '0' for blanks
	input := `
		+---+---+---+
		|100|000|000|
		|020|000|000|
		|003|000|000|
		+---+---+---+
		|000|400|000|
		|000|050|000|
		|000|006|000|
		+---+---+---+
		|000|000|700|
		|000|000|080|
		|000|000|009|
		+---+---+---+
	`

	g, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error on formatted input: %v", err)
	}

	if g.Clues != 9 {
		t.Errorf("Expected 9 clues, got %d", g.Clues)
	}
	if g.Values[0] != 1 || g.Values[10] != 2 || g.Values[20] != 3 {
		t.Errorf("Values not parsed correctly: [0]=%d, [10]=%d, [20]=%d", g.Values[0], g.Values[10], g.Values[20])
	}
	if g.Values[30] != 4 || g.Values[40] != 5 || g.Values[50] != 6 {
		t.Errorf("Values not parsed correctly: [30]=%d, [40]=%d, [50]=%d", g.Values[30], g.Values[40], g.Values[50])
	}
	if g.Values[60] != 7 || g.Values[70] != 8 || g.Values[80] != 9 {
		t.Errorf("Values not parsed correctly: [60]=%d, [70]=%d, [80]=%d", g.Values[60], g.Values[70], g.Values[80])
	}
}

func TestParse_TooShortInput(t *testing.T) {
	input := "12345"
	_, err := Parse(input)
	if err == nil {
		t.Errorf("Expected error for input with fewer than 81 cells, got nil")
	}
}

func TestParse_ExtraCharactersIgnored(t *testing.T) {
	// 81 dots followed by extra characters
	input := strings.Repeat(".", 81) + "extra trailing text 123456"
	g, err := Parse(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if g.Clues != 0 {
		t.Errorf("Expected 0 clues, got %d", g.Clues)
	}
}

func TestInitializeCandidates_EliminationAcrossHouses(t *testing.T) {
	// Cell 0 is '5'
	// Peers in row 0, col 0, and box 0 should all have candidate 5 removed
	input := "5" + strings.Repeat(".", 80)
	g, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	g.InitializeCandidates()

	// Row 0 peer: cell 8
	if g.Candidates[8]&(1<<5) != 0 {
		t.Errorf("Candidate 5 should be eliminated in row peer cell 8")
	}
	// Col 0 peer: cell 72 (row 8, col 0)
	if g.Candidates[72]&(1<<5) != 0 {
		t.Errorf("Candidate 5 should be eliminated in col peer cell 72")
	}
	// Box 0 peer: cell 20 (row 2, col 2)
	if g.Candidates[20]&(1<<5) != 0 {
		t.Errorf("Candidate 5 should be eliminated in box peer cell 20")
	}
	// Non-peer: cell 12 (row 1, col 3)
	if g.Candidates[12]&(1<<5) == 0 {
		t.Errorf("Candidate 5 should not be eliminated in non-peer cell 12")
	}
}
