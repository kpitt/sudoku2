package solver

// Global Look-Up Tables.
var (
	RowLUT        [81]int
	ColLUT        [81]int
	BoxLUT        [81]int
	PeersLUT      [81][20]int
	HouseLUT      [27][9]int
	CellHousesLUT [81][3]int
)

func init() {
	// Initialize Row, Col, Box LUTs
	for i := range 81 {
		r := i / 9
		c := i % 9
		b := (r/3)*3 + (c / 3)
		RowLUT[i] = r
		ColLUT[i] = c
		BoxLUT[i] = b

		CellHousesLUT[i][0] = r      // Row house index 0-8
		CellHousesLUT[i][1] = 9 + c  // Col house index 9-17
		CellHousesLUT[i][2] = 18 + b // Box house index 18-26
	}

	// Initialize HouseLUT
	// Rows 0-8
	for r := range 9 {
		for c := range 9 {
			HouseLUT[r][c] = r*9 + c
		}
	}
	// Cols 9-17
	for c := range 9 {
		for r := range 9 {
			HouseLUT[9+c][r] = r*9 + c
		}
	}
	// Boxes 18-26
	for b := range 9 {
		br := (b / 3) * 3
		bc := (b % 3) * 3
		idx := 0

		for r := range 3 {
			for c := range 3 {
				HouseLUT[18+b][idx] = (br+r)*9 + (bc + c)
				idx++
			}
		}
	}

	// Initialize PeersLUT
	for i := range 81 {
		peerIdx := 0
		r, c, b := RowLUT[i], ColLUT[i], BoxLUT[i]

		// Use a simple set to avoid duplicates
		seen := make(map[int]bool)
		seen[i] = true

		// Peers in same row
		for _, cell := range HouseLUT[r] {
			if !seen[cell] {
				PeersLUT[i][peerIdx] = cell
				peerIdx++
				seen[cell] = true
			}
		}
		// Peers in same col
		for _, cell := range HouseLUT[9+c] {
			if !seen[cell] {
				PeersLUT[i][peerIdx] = cell
				peerIdx++
				seen[cell] = true
			}
		}
		// Peers in same box
		for _, cell := range HouseLUT[18+b] {
			if !seen[cell] {
				PeersLUT[i][peerIdx] = cell
				peerIdx++
				seen[cell] = true
			}
		}
	}
}
