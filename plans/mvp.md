# Implementation Plan - Sudoku CLI MVP

## 1. 🔍 Analysis & Context
*   **Objective:** Implement a Minimum Viable Product (MVP) of a zero-allocation, high-performance Sudoku CLI solver handling 81-character input strings, "Simple Sudoku" ASCII output, and a deductive loop with Singles and Pairs techniques.
*   **Affected Files:** 
    *   `go.mod` & `go.sum`
    *   `cmd/sudoku/main.go`, `cmd/sudoku/root.go`, `cmd/sudoku/solve.go`
    *   `internal/solver/board.go`, `internal/solver/lut.go`, `internal/solver/step.go`
    *   `internal/parser/parser.go`
    *   `internal/formatter/formatter.go`
    *   `internal/solver/singles.go`, `internal/solver/subsets.go`, `internal/solver/loop.go`
*   **Key Dependencies:** `github.com/spf13/cobra` (CLI framework).
*   **Risks/Unknowns:** Correctly managing the bitwise logic for subset duals (Pairs) without heap allocations. Ensuring strict adherence to `go build -gcflags="-m"` zero-allocation constraints.

## 2. 📋 Checklist
- [x] Step 1: Initialize CLI Structure
    Status: ✅ Implemented in cmd/sudoku/
- [x] Step 2: Implement Core Board State & LUTs
    Status: ✅ Implemented in internal/solver/
- [x] Step 3: Implement Basic Parser
    Status: ✅ Implemented in internal/parser/
- [x] Step 4: Implement `ss` Formatter
    Status: ✅ Implemented in internal/formatter/
- [x] Step 5: Implement Naked & Hidden Singles
    Status: ✅ Implemented in internal/solver/singles.go
- [x] Step 6: Implement Naked & Hidden Pairs
    Status: ✅ Implemented in internal/solver/subsets.go
- [x] Step 7: Implement Solver Priority Loop & CLI Integration
    Status: ✅ Implemented in internal/solver/loop.go and cmd/sudoku/solve.go
- [x] Verification
    Status: ✅ All tests passed, zero allocations verified, lint clean.

## 3. 📝 Step-by-Step Implementation Details

### Step 1: Initialize CLI Structure
*   **Goal:** Setup `spf13/cobra` and basic subcommands.
*   **Action:**
    *   Run `go get github.com/spf13/cobra`.
    *   Modify `cmd/sudoku/main.go` to execute the root command: `rootCmd.Execute()`.
    *   Create `cmd/sudoku/root.go` defining the `rootCmd`.
    *   Create `cmd/sudoku/solve.go` defining the `solveCmd`. Accept an 81-character positional argument.

### Step 2: Implement Core Board State & LUTs
*   **Goal:** Create the zero-allocation data structures.
*   **Action:**
    *   Create `internal/solver/board.go`. Define `Board` struct with `Cells [81]uint16`, `RowState [9]uint16`, `ColState [9]uint16`, `BoxState [9]uint16`, and `Resolved uint8`. Add a method `func (b *Board) RemoveCandidate(cellIdx int, val uint16) bool` to handle propagation.
    *   Create `internal/solver/step.go`. Define `ActionType` enum and `Step` struct `type Step struct { Technique string; Target []int; Values []int; Action ActionType }`.
    *   Create `internal/solver/lut.go`. Precompute global tables in `init()`: `RowLUT`, `ColLUT`, `BoxLUT`, `PeersLUT` (20 peers per cell), `HouseLUT` (27 houses: 9 rows, 9 cols, 9 boxes).

### Step 3: Implement Basic Parser
*   **Goal:** Convert 81-character input strings to a `Board` state.
*   **Action:**
    *   Create `internal/parser/parser.go`. Add `func Parse(input string) (*solver.Board, error)`.
    *   Implement regex or string trimming to reduce input to exactly 81 characters.
    *   Loop over characters. Initialize empty cells with `0x01FF`. Initialize given digits with `1 << (digit - 1)`.

### Step 4: Implement `ss` Formatter
*   **Goal:** Output the grid in the required "Simple Sudoku" ASCII format.
*   **Action:**
    *   Create `internal/formatter/formatter.go`. Add `func FormatSS(b *solver.Board) string`.
    *   Iterate over the 81 cells. Identify solved cells (`bits.OnesCount16 == 1`) and extract the digit using `bits.TrailingZeros16`. Use `.` for unsolved cells.
    *   Wrap rows in ASCII borders (e.g., `*-----------*`, `|---+---+---|`).

### Step 5: Implement Naked & Hidden Singles
*   **Goal:** Detect the simplest deductive resolutions.
*   **Action:**
    *   Create `internal/solver/singles.go`.
    *   Add `func FindNakedSingles(b *Board) *Step`. Loop `Cells`, find cells with `bits.OnesCount16 == 1` that aren't marked as resolved yet, and propagate eliminations to peers.
    *   Add `func FindHiddenSingles(b *Board) *Step`. Loop the 27 houses in `HouseLUT`. Use the bitwise frequency accumulator (`seenOnce`, `seenMultiple`) to find unique candidates in a house. Return the placement.

### Step 6: Implement Naked & Hidden Pairs
*   **Goal:** Detect pairs using subset logic.
*   **Action:**
    *   Create `internal/solver/subsets.go`.
    *   Add `func FindNakedPairs(b *Board) *Step`. For each house, evaluate combinations of 2 cells where the union of their masks has exactly 2 bits active. Eliminate those candidates from the other 7 cells in the house.
    *   Add `func FindHiddenPairs(b *Board) *Step`. For each house, evaluate combinations of 2 candidates. Check if those candidates appear in exactly the same 2 cells. If so, eliminate all other candidates from those 2 cells.

### Step 7: Implement Solver Priority Loop & CLI Integration
*   **Goal:** Coordinate the solver execution and print outputs.
*   **Action:**
    *   Create `internal/solver/loop.go`. Add `func Solve(b *Board) []Step`.
    *   Implement an infinite loop. Priority sequence: Naked Singles -> Hidden Singles -> Naked Pairs -> Hidden Pairs. If any returns a `Step`, record it, apply it, and `continue` to restart the loop. If none find a step, `break` (stalled or finished).
    *   Update `cmd/sudoku/solve.go` to print the initial board via `FormatSS`, execute `solver.Solve`, optionally print steps, and print the final board via `FormatSS`.

## 4. 🧪 Testing Strategy
*   **Unit Tests:**
    *   `parser_test.go`: Verify empty cell markers, error on invalid lengths.
    *   `formatter_test.go`: Assert exact string match against an expected `ss` ASCII board.
    *   `solver_test.go`: Synthetic boards testing `FindNakedSingles`, `FindHiddenSingles`, `FindNakedPairs`, `FindHiddenPairs`. Verify correct `Step` is returned.
*   **Integration Tests:**
    *   Use an 81-character known puzzle (solvable by pairs/singles) as standard input to the solver package, verify `Board.Resolved == 81`.
*   **Manual Verification:**
    *   Run `go run ./cmd/sudoku solve "310004069000000200008005040000000005006000017807030000590700006600003050000100002"`. Ensure it parses, runs the loop, and prints initial and final grids.

## 5. ✅ Success Criteria
*   Application compiles without errors.
*   `golangci-lint run` reports no warnings.
*   Zero heap allocations during the solver loop (`go build -gcflags="-m"` verified).
*   CLI successfully accepts an 81-character puzzle, applies Singles and Pairs logic, and outputs the final grid in `ss` ASCII format.
*   All unit and integration tests pass.