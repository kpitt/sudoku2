package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	input := "310004069000000200008005040000000005006000017807030000590700006600003050000100002"

	b, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if b.Resolved < 25 { // Just a rough check for this specific puzzle
		t.Errorf("Expected at least 25 resolved cells, got %d", b.Resolved)
	}

	// Check a specific cell (0,0) which is '3'
	if b.Cells[0] != (1 << 2) {
		t.Errorf("Expected cell 0 to be '3' (mask 4), got %d", b.Cells[0])
	}
}
