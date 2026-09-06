# Sudoku Solver Architecture Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a high-performance Go Sudoku solver with a Data-Oriented Design grid, deductive logic pipeline, brute-force fallback, and a `cobra` CLI.

**Architecture:** DOD with Bitmasks for grid state (`[81]byte` and `[81]uint16`), pre-computed lookups for houses and peers. A deductive pipeline looping over techniques, falling back to recursive backtracking if stalled.

**Tech Stack:** Go 1.24+, `spf13/cobra`, `fatih/color`

**Spec:** `docs/superpowers/specs/2026-09-05-sudoku-solver-design.md`

## Global Constraints

- Written in Go language (1.24+)
- Must use zero allocations inside hot search loops (Brute Force fallback)
- Utility functions must use semantic domain-specific names (e.g., `CandidatesCount`)
- Use modern Go 1.24 idioms including range-over-int (e.g., `for i := range 81`)

---

### Task 1: Project Setup and Core Grid Data Structures

**Files:**
- Create: `go.mod`
- Create: `internal/sudoku/grid.go`
- Test: `internal/sudoku/grid_test.go`

**Interfaces:**
- Produces: `type Grid struct { Values [81]byte; Candidates [81]uint16; Clues int }`
- Produces: `func (g *Grid) CandidatesCount(i int) int`
- Produces: `func (g *Grid) Display(style string)`

- [ ] **Step 1: Initialize Go module**

```bash
go mod init github.com/sudoku-solver
```

- [ ] **Step 2: Write failing test for Grid initialization**

```go
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
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./internal/sudoku -v`
Expected: FAIL (Grid undefined)

- [ ] **Step 4: Implement Grid struct and semantic utilities**

```go
package sudoku

import "math/bits"

// Grid represents the 9x9 Sudoku board
type Grid struct {
	Values     [81]byte
	Candidates [81]uint16
	Clues      int
}

// CandidatesCount returns the number of valid candidates for a cell
func (g *Grid) CandidatesCount(i int) int {
    return bits.OnesCount16(g.Candidates[i])
}
```

- [ ] **Step 5: Run test to verify it passes**

Run: `go test ./internal/sudoku -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add go.mod internal/sudoku/grid.go internal/sudoku/grid_test.go
git commit -m "feat: initialize module and Grid struct"
```

---

### Task 2: Pre-computed Lookups (Houses and Peers)

**Files:**
- Create: `internal/sudoku/lookups.go`
- Test: `internal/sudoku/lookups_test.go`

**Interfaces:**
- Produces: `var Houses [27][9]int`
- Produces: `var Peers [81][20]int`

- [ ] **Step 1: Write failing test for lookups**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/sudoku -v`
Expected: FAIL (initLookups, Houses, Peers undefined)

- [ ] **Step 3: Implement Lookups using range-over-int**

```go
package sudoku

var Houses [27][9]int
var Peers [81][20]int

func init() {
    initLookups()
}

func initLookups() {
    // Initialize Houses (0-8 Rows, 9-17 Cols, 18-26 Boxes)
    for i := range 81 {
        r := i / 9
        c := i % 9
        b := (r/3)*3 + (c/3)
        
        Houses[r][c] = i
        Houses[9+c][r] = i
        // Box indexing logic
        boxPos := (r%3)*3 + (c%3)
        Houses[18+b][boxPos] = i
    }
    
    // Initialize Peers
    for i := range 81 {
        r := i / 9
        c := i % 9
        b := (r/3)*3 + (c/3)
        
        peerIdx := 0
        for j := range 81 {
            if i == j { continue }
            jr := j / 9
            jc := j % 9
            jb := (jr/3)*3 + (jc/3)
            if r == jr || c == jc || b == jb {
                Peers[i][peerIdx] = j
                peerIdx++
            }
        }
    }
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/sudoku -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/sudoku/lookups.go internal/sudoku/lookups_test.go
git commit -m "feat: implement pre-computed Houses and Peers lookups"
```

---

### Task 3: I/O Parsing and Initialization

**Files:**
- Create: `internal/sudoku/parser.go`
- Test: `internal/sudoku/parser_test.go`

**Interfaces:**
- Consumes: `Grid`
- Produces: `func Parse(input string) (*Grid, error)`
- Produces: `func (g *Grid) InitializeCandidates()`

- [ ] **Step 1: Write failing test for Parsing**

```go
package sudoku

import "testing"

func TestParseAndInitialize(t *testing.T) {
    // 81 char string with one '1'
    input := "1" + string(make([]byte, 80))
    for i := 1; i < 81; i++ {
        input = input[:i] + "." + input[i+1:]
    }
    
    g, err := Parse(input)
    if err != nil {
        t.Fatalf("Parse error: %v", err)
    }
    if g.Values[0] != 1 {
        t.Errorf("Expected 1, got %d", g.Values[0])
    }
    
    g.InitializeCandidates()
    // Cell 1 is in same row, shouldn't have '1' as candidate
    if g.Candidates[1] & (1 << 1) != 0 {
         t.Errorf("Candidate 1 should be eliminated from peer")
    }
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./internal/sudoku -v`
Expected: FAIL

- [ ] **Step 3: Implement Parser and Initializer**

```go
package sudoku

import "strings"

func Parse(input string) (*Grid, error) {
    g := &Grid{}
    clean := strings.Map(func(r rune) rune {
        if r >= '0' && r <= '9' || r == '.' { return r }
        return -1
    }, input)
    
    idx := 0
    for _, ch := range clean {
        if idx >= 81 { break }
        if ch >= '1' && ch <= '9' {
            g.Values[idx] = byte(ch - '0')
            g.Clues++
        }
        idx++
    }
    return g, nil
}

func (g *Grid) InitializeCandidates() {
    // Start with all candidates (bits 1-9)
    for i := range 81 {
        if g.Values[i] == 0 {
            g.Candidates[i] = 0b1111111110 // Bits 1-9 set
        }
    }
    
    // Eliminate based on givens
    for i := range 81 {
        val := g.Values[i]
        if val > 0 {
            mask := uint16(1 << val)
            for _, p := range Peers[i] {
                g.Candidates[p] &^= mask
            }
        }
    }
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/sudoku -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/sudoku/parser.go internal/sudoku/parser_test.go
git commit -m "feat: add parser and candidate initialization"
```

---

### Task 4: CLI Architecture (cobra)

**Files:**
- Create: `cmd/sudoku/main.go`
- Create: `cmd/sudoku/root.go`
- Create: `cmd/sudoku/solve.go`

- [ ] **Step 1: Install cobra**

```bash
go get -u github.com/spf13/cobra@latest
```

- [ ] **Step 2: Write basic CLI structure**

```go
// cmd/sudoku/main.go
package main

func main() {
    Execute()
}

// cmd/sudoku/root.go
package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
    Use:   "sudoku",
    Short: "A high-performance Sudoku solver",
}

func Execute() {
    rootCmd.Execute()
}

// cmd/sudoku/solve.go
package main

import (
    "fmt"
    "github.com/spf13/cobra"
)

var solveCmd = &cobra.Command{
    Use:   "solve [puzzle]",
    Short: "Solve a sudoku puzzle",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Solving:", args[0])
    },
}

func init() {
    rootCmd.AddCommand(solveCmd)
}
```

- [ ] **Step 3: Run the CLI**

Run: `go run ./cmd/sudoku solve "..."`
Expected: Outputs "Solving: ..."

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum cmd/sudoku/
git commit -m "feat: scaffold cobra CLI"
```

---

### Task 5: Brute-Force Fallback Solver

**Files:**
- Create: `internal/sudoku/bruteforce.go`
- Test: `internal/sudoku/bruteforce_test.go`

**Interfaces:**
- Consumes: `Grid`
- Produces: `func (g *Grid) SolveBruteForce() bool`

- [ ] **Step 1: Write failing test**

```go
package sudoku

import "testing"

func TestBruteForce(t *testing.T) {
    input := ".....6....59.....82....8....45........3........6..3.54...325..6.................."
    g, _ := Parse(input)
    g.InitializeCandidates()
    solved := g.SolveBruteForce()
    if !solved || g.Clues != 81 {
        t.Errorf("Expected board to be solved")
    }
}
```

- [ ] **Step 2: Run test**

Run: `go test ./internal/sudoku -v`
Expected: FAIL (SolveBruteForce undefined)

- [ ] **Step 3: Implement recursive backtracking (zero alloc)**

```go
package sudoku

func (g *Grid) SolveBruteForce() bool {
    if g.Clues == 81 {
        return true
    }
    
    // Find cell with fewest candidates
    minCand := 10
    bestCell := -1
    for i := range 81 {
        if g.Values[i] == 0 {
            c := g.CandidatesCount(i)
            if c < minCand {
                minCand = c
                bestCell = i
            }
        }
    }
    
    if bestCell == -1 || minCand == 0 {
        return false // Unsolvable
    }
    
    // Backup state (stored on stack, zero heap alloc)
    var backup [81]uint16
    copy(backup[:], g.Candidates[:])
    
    cands := g.Candidates[bestCell]
    for d := 1; d <= 9; d++ { // Using traditional loop here since we need exactly 1-9 for mask shifting
        if cands & (1 << d) != 0 {
            g.Values[bestCell] = byte(d)
            g.Clues++
            
            // Eliminate from peers
            mask := uint16(1 << d)
            for _, p := range Peers[bestCell] {
                g.Candidates[p] &^= mask
            }
            
            if g.SolveBruteForce() {
                return true
            }
            
            // Backtrack
            g.Values[bestCell] = 0
            g.Clues--
            copy(g.Candidates[:], backup[:])
        }
    }
    return false
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./internal/sudoku -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/sudoku/bruteforce.go internal/sudoku/bruteforce_test.go
git commit -m "feat: add zero-allocation brute force solver"
```
