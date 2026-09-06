# Sudoku Solver Architecture Design
Date: 2026-09-05

## 1. Core Grid & Data Model

To maximize CPU cache locality and avoid heap allocations, we are utilizing a Data-Oriented Design (DOD). 

### Data Structures
*   **`Grid` Struct**:
    *   `Values [81]byte`: Stores the known digit (1-9), or `0` if empty.
    *   `Candidates [81]uint16`: A bitmask where bits 1-9 represent valid candidates. `uint16` allows for ultra-fast bitwise operations (e.g., `candidates & (1 << digit)`) and zero-allocation candidate elimination.
    *   `Clues int`: A simple counter of filled cells (81 means solved).

### Lookups & Constants
We pre-compute lookup tables at startup to avoid re-calculating indices inside hot loops:
*   **`Houses[27][9]int`**: A unified array allowing iteration over all Rows (0-8), Columns (9-17), and Boxes (18-26) using generic logic loops. This prevents duplicating logic for techniques like "Naked Pair" across three different house types.
*   **`Peers[81][20]int`**: For any given flat cell index (0-80), this returns the 20 flat indices of its peers (cells in the same row, column, or box). This is the fastest way to handle candidate elimination when a digit is placed.
*   **Coordinate Math**: For specific one-off cell targeting, we use math: `row = i / 9`, `col = i % 9`, `box = (row / 3) * 3 + (col / 3)`.

### Utility Functions
*   We will use semantic, domain-specific naming (e.g., `CandidatesCount`, `HasCandidate`, `RemoveCandidate`) rather than generic bitwise names (e.g., `CountBits`), to ensure the logic algorithms read closely to plain English.

---

## 2. Solver Architecture

The solver manages both human-like deductive logic and a brute-force fallback mechanism.

### The Deductive Pipeline
*   **Technique Interface**: We will define an interface or function signature for a `Technique` (e.g., `func (g *Grid) (changed bool, step *Step)`).
*   **Execution**: 
    *   The solver maintains an ordered list of the 15 Simple Sudoku Technique Set (SSTS) techniques:
        1. Naked Single
        2. Hidden Single
        3. Locked Candidates (Pointing)
        4. Locked Candidates (Claiming)
        5. Naked Pair
        6. Naked Triple
        7. Naked Quad
        8. Hidden Pair
        9. Hidden Triple
        10. Hidden Quad
        11. X-Wing
        12. Swordfish
        13. Simple Colors
        14. Multi-Colors
        15. XY-Wing
    *   It loops through the techniques in order. If a technique makes *any* change (like eliminating a candidate or placing a value), the pipeline immediately restarts from Technique 1 (Naked Single), as simpler techniques may now be unblocked.
    *   If the pipeline runs through all 15 techniques without making a single change, the deductive engine stalls.

### Brute-Force Fallback
*   If the deductive engine stalls and the board is not fully solved, and the `--brute-force` flag is enabled, it hands the current grid state to a highly optimized, recursive backtracking algorithm.
*   This algorithm tries valid digits until it finds a solution or proves the puzzle unsolvable.

### Step Logging (HoDoKu format)
*   Techniques will return a `Step` struct (or `nil`) detailing their actions (e.g., "Naked Pair on 3,5 in Row 2, eliminated 3 from r2c8"). 
*   These steps will be appended to a log to support the HoDoKu output format.

---

## 3. I/O and CLI Architecture

### CLI Architecture (`cobra`)
*   `sudoku solve [file|string]`: Solves the puzzle (with optional `--step-log` and `--brute-force` flags).
*   `sudoku check [file|string]`: Tests if a given puzzle state is solvable (identifying if there are 0, 1, or multiple valid solutions).
*   `sudoku print [file|string]`: Renders a grid in a human-readable format.
*   `sudoku generate`: Scaffolded command to generate new puzzles (future scope).

### I/O Parsing
*   **`Parse(input string) (*Grid, error)`**: Normalizes the input string (stripping whitespace, borders, non-digit characters) to extract exactly 81 characters. `.` or `0` denote empty cells, and `1-9` denote clues. **Note**: In accordance with the Single Responsibility Principle, `Parse` strictly constructs the board state (`Values`) without performing candidate eliminations.
*   **Initialization**: A separate `Initialize(grid)` or `solver.New(grid)` function will use the `Peers` lookup to populate the initial `Candidates` bitmask, avoiding duplication of solver logic.

### Output Formatting
*   A `Display(style string)` method prints the 81-character string, a standard 9x9 grid, or a full "pencil mark" (candidates) grid if the puzzle is unsolved.
*   The `fatih/color` library will be used to highlight solved digits versus given clues in terminal output.
