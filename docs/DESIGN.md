# Architectural Design: Sudoku Solver (SSTS)

## Introduction

This solver implements the Simple Sudoku Technique Set (SSTS). Our goal was to build a engine that prioritizes deductive logic over brute force. While something like Knuth's Algorithm X is faster for simply finding a solution, it doesn't help with "human-style" hints or step-by-step logging. We use brute force only as a fallback when deductive logic hits a wall.

Performance-wise, we focused on three things: cache locality, bitwise operations, and zero heap allocations. Most solvers use pointers and objects that live on the heap, which is fine for simple scripts but slow for high-throughput tasks. By keeping everything in contiguous bitmask arrays on the stack, we stay in the CPU's L1 cache and avoid garbage collection entirely.

## Memory Layout and Data Structures

The board state is represented by bitmasks. We use `uint16` for each cell, where bits 0-8 represent digits 1-9. A mask of `0x01FF` (binary `111111111`) means every digit is still a candidate.

### Board State

The `Board` struct holds the grid and some metadata to speed up lookups. It's 217 bytes, which fits comfortably in four 64-byte cache lines. We pass this around by pointer to avoid copying, but the Go compiler keeps it on the stack because it never escapes the local scope.

| Field | Type | Size | Notes |
| :--- | :--- | :--- | :--- |
| `Cells` | `[81]uint16` | 162B | The 9x9 grid in row-major order. |
| `Resolved` | `uint8` | 1B | Tracks how many cells are solved (0-81). |
| `RowState` | `[9]uint16` | 18B | Digits already placed in each row. |
| `ColState` | `[9]uint16` | 18B | Digits already placed in each column. |
| `BoxState` | `[9]uint16` | 18B | Digits already placed in each 3x3 box. |

These state arrays let us do $O(1)$ collision checks. If we need to know if '5' is already in row 3, we just check bit 4 of `RowState[3]`.

### Precomputed Tables (LUTs)

We use global lookup tables for everything that would otherwise require division or modulo. Division is slow, and for a solver doing millions of iterations, it adds up.

- **Coordinate Mapping:** `RowLUT`, `ColLUT`, and `BoxLUT` map a flat index (0-80) to its house index.
- **Peers:** `PeersLUT` gives us the 20 intersecting cells for any index in a single array access.
- **Houses:** `HouseLUT` defines the indices for all 27 houses (9 rows, 9 columns, 9 boxes).

## Bitwise Operations

We use the `math/bits` package for everything. These functions usually compile down to single hardware instructions like `POPCNT` or `TZCNT`.

- `bits.OnesCount16(mask)`: Returns the number of candidates.
- `mask & -mask`: Isolates the lowest set bit.
- `bits.TrailingZeros16(mask)`: Converts a bitmask back to a digit index.
- `(A &^ B) == 0`: Fast check to see if subset A is entirely contained within B.

## Group 1: The Basics

### Naked Singles
If a cell has only one candidate left, it's solved. We update the cell and immediately strip that digit from its 20 peers. This often triggers a chain reaction, so we use a simple stack-based queue to process these without rescanning the whole board.

### Hidden Singles
Sometimes a digit can only go in one spot in a row, even if that spot has other candidates. We find these by counting how many times each digit appears in a house. If a digit's count is 1, we found a hidden single.

```go
seenOnce, seenMultiple := uint16(0), uint16(0)
for _, cellIdx := range HouseLUT[houseId] {
    mask := board.Cells[cellIdx]
    seenMultiple |= (seenOnce & mask)
    seenOnce |= mask
}
uniqueMask := seenOnce &^ seenMultiple
```

## Group 2: Intersections

These techniques look for candidates that are "locked" into a specific intersection between a box and a line.

- **Pointing:** All instances of a digit in a box are in one row. That digit can't be anywhere else in that row.
- **Box-Line Reduction:** All instances of a digit in a row are inside one box. That digit can't be anywhere else in that box.

## Group 3: Subsets (Pairs, Triples, Quads)

A Naked Subset is when $N$ cells in a house share $N$ candidates. For example, if two cells in a row both only have $\{1, 2\}$, then 1 and 2 can't be anywhere else in that row.

We handle Naked and Hidden subsets with the same engine. To find Hidden subsets, we just transpose the house data (turning "which digits are in this cell" into "which cells contain this digit") and run the Naked subset logic. This kept the code much cleaner.

We use unrolled loops for pairs, triples, and quads to keep the branch predictor happy. For quads, we skip cells that have more than 4 candidates early on to save time.

## Group 4: Fish (X-Wing to Jellyfish)

Fish patterns (X-Wing, Swordfish, Jellyfish) are just orthogonal projections. If a digit appears in $N$ rows but is restricted to the same $N$ columns, we can clear that digit from those columns in all other rows.

The solver uses one generic function for all fish sizes. It builds bitmasks for each row and looks for combinations of $N$ rows where the union of their masks has exactly $N$ bits set.

## Group 5: Wings

An XY-Wing uses a "pivot" cell with two candidates $\{X, Y\}$ and two "pincer" cells with $\{X, Z\}$ and $\{Y, Z\}$. No matter what the pivot turns out to be, one of the pincers must be $Z$. So any cell that "sees" both pincers can't be $Z$.

We find these by scanning for cells with exactly 2 candidates and checking their peer networks for the required pincer pattern.

## Group 6: Coloring

Coloring tracks "conjugate pairs"—houses where a digit appears exactly twice. These form a binary link: if one is true, the other is false. We build a graph of these links and use a simple BFS to alternate colors.

- **Contradictions:** If two cells with the same color are peers, that color is impossible.
- **Multi-Colors:** If an uncolored cell sees both colors in a chain, it can't be that digit.

We use stack arrays for the graph adjacency list to avoid heap allocations.

## Execution Order

The solver is a loop that escalates. We always try the easiest thing first. If a complex technique (like an XY-Wing) manages to clear a candidate, we immediately jump back to the start of the loop to see if that triggered any new Naked Singles.

| Level | Technique | Cost |
| :--- | :--- | :--- |
| 0 | Naked/Hidden Singles | $O(1)$ |
| 1 | Intersections | $O(1)$ |
| 2 | Subsets (2-4) | $O(N^4)$ |
| 3 | Fish (2-4) | $O(R^4)$ |
| 4 | Wings / Coloring | $O(V+E)$ |

By keeping the data structures simple and the loops tight, we can run through millions of these checks per second.
