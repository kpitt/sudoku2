package sudoku

import "testing"

func TestLookups(t *testing.T) {
	initLookups()

	// Test Houses: Row 0 should contain indices 0-8
	if Houses[0][8] != 8 {
		t.Errorf("Expected Houses[0][8] to be 8, got %d", Houses[0][8])
	}

	// Test Peers: Peers of index 0 should contain index 1
	hasPeer := false
	for _, p := range Peers[0] {
		if p == 1 {
			hasPeer = true
			break
		}
	}
	if !hasPeer {
		t.Errorf("Expected Peers[0] to contain 1")
	}
}

func TestHouses_Structure(t *testing.T) {
	// 27 houses: 0-8 rows, 9-17 columns, 18-26 boxes
	if len(Houses) != 27 {
		t.Fatalf("Expected 27 houses, got %d", len(Houses))
	}

	// Check Rows (0-8)
	for r := range 9 {
		seen := make(map[int]bool)
		for c := range 9 {
			idx := Houses[r][c]
			if idx != r*9+c {
				t.Errorf("House[%d][%d] = %d, expected %d", r, c, idx, r*9+c)
			}
			seen[idx] = true
		}
		if len(seen) != 9 {
			t.Errorf("Row %d has duplicate cells", r)
		}
	}

	// Check Columns (9-17)
	for c := range 9 {
		houseIdx := 9 + c
		seen := make(map[int]bool)
		for r := range 9 {
			idx := Houses[houseIdx][r]
			expected := r*9 + c
			if idx != expected {
				t.Errorf("House[%d][%d] = %d, expected %d", houseIdx, r, idx, expected)
			}
			seen[idx] = true
		}
		if len(seen) != 9 {
			t.Errorf("Column %d has duplicate cells", c)
		}
	}

	// Check Boxes (18-26)
	for b := range 9 {
		houseIdx := 18 + b
		boxStartRow := (b / 3) * 3
		boxStartCol := (b % 3) * 3
		seen := make(map[int]bool)
		for pos := range 9 {
			idx := Houses[houseIdx][pos]
			r := idx / 9
			c := idx % 9
			actualBox := (r/3)*3 + (c / 3)
			if actualBox != b {
				t.Errorf("House[%d][%d] cell %d is in box %d, expected box %d", houseIdx, pos, idx, actualBox, b)
			}
			expectedRow := boxStartRow + pos/3
			expectedCol := boxStartCol + pos%3
			expectedIdx := expectedRow*9 + expectedCol
			if idx != expectedIdx {
				t.Errorf("House[%d][%d] = %d, expected %d", houseIdx, pos, idx, expectedIdx)
			}
			seen[idx] = true
		}
		if len(seen) != 9 {
			t.Errorf("Box %d has duplicate cells", b)
		}
	}
}

func TestPeers_Structure(t *testing.T) {
	for i := range 81 {
		r := i / 9
		c := i % 9
		b := (r/3)*3 + (c / 3)

		if len(Peers[i]) != 20 {
			t.Fatalf("Peers[%d] has length %d, expected 20", i, len(Peers[i]))
		}

		seen := make(map[int]bool)
		for _, peer := range Peers[i] {
			if peer < 0 || peer >= 81 {
				t.Errorf("Peers[%d] contains out-of-range index %d", i, peer)
			}
			if peer == i {
				t.Errorf("Peers[%d] contains itself", i)
			}
			if seen[peer] {
				t.Errorf("Peers[%d] contains duplicate peer %d", i, peer)
			}
			seen[peer] = true

			pr := peer / 9
			pc := peer % 9
			pb := (pr/3)*3 + (pc / 3)

			if pr != r && pc != c && pb != b {
				t.Errorf("Peers[%d] contains %d which does not share row, col, or box", i, peer)
			}
		}

		if len(seen) != 20 {
			t.Errorf("Peers[%d] should have 20 unique peers, got %d", i, len(seen))
		}

		// Verify every true peer is in Peers[i]
		for j := range 81 {
			if j == i {
				continue
			}
			jr := j / 9
			jc := j % 9
			jb := (jr/3)*3 + (jc / 3)
			if jr == r || jc == c || jb == b {
				if !seen[j] {
					t.Errorf("Peers[%d] missing peer %d", i, j)
				}
			}
		}
	}
}
