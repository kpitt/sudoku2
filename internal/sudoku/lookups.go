package sudoku

// Houses contains the 27 units of 9 cells each:
// Indices 0-8: Rows 0-8
// Indices 9-17: Columns 0-8
// Indices 18-26: 3x3 Boxes 0-8
var Houses [27][9]int

// Peers contains the 20 peers for each of the 81 cells:
// A peer is any cell in the same row, column, or 3x3 box (excluding the cell itself).
var Peers [81][20]int

func init() {
	initLookups()
}

func initLookups() {
	// Initialize Houses (0-8 Rows, 9-17 Cols, 18-26 Boxes)
	for i := range 81 {
		r := i / 9
		c := i % 9
		b := (r/3)*3 + (c / 3)

		Houses[r][c] = i
		Houses[9+c][r] = i
		// Box indexing logic
		boxPos := (r%3)*3 + (c % 3)
		Houses[18+b][boxPos] = i
	}

	// Initialize Peers
	for i := range 81 {
		r := i / 9
		c := i % 9
		b := (r/3)*3 + (c / 3)

		peerIdx := 0
		for j := range 81 {
			if i == j {
				continue
			}
			jr := j / 9
			jc := j % 9
			jb := (jr/3)*3 + (jc / 3)
			if r == jr || c == jc || b == jb {
				Peers[i][peerIdx] = j
				peerIdx++
			}
		}
	}
}
