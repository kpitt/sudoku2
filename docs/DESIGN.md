# Architectural Design for SSTS Algorithms

## Deductive Sudoku Resolution and Systems Architecture

A Sudoku solver requires balancing algorithm design with systems architecture to execute deductive constraints. Backtracking algorithms like Donald Knuth’s Algorithm X with Dancing Links (DLX) work well for exact cover problems but rely on trial and error. The Simple Sudoku Technique Set (SSTS) provides a framework of deterministic strategies used by human experts to reduce the puzzle state. The SSTS includes techniques from single-candidate deductions to multi-point graph-coloring algorithms like XY-Wings and Multi-Colors.

This document outlines a zero-allocation Go (Golang) solver designed to implement the SSTS. The architecture maximizes throughput by enforcing CPU cache locality, avoiding branch mispredictions, using hardware-accelerated bitwise instructions, and relying on stack memory instead of heap allocations.

Many Sudoku solvers use pointer-linked structures on the heap. This degrades CPU cache performance and can cause false sharing in multithreaded environments. This architecture uses contiguous bitmask arrays instead. The board state and constraint matrices fit into L1 CPU cache lines to minimize memory access latency. This specification explains how to translate SSTS logic into Go routines, sharing functionality across mathematically similar properties to reduce code duplication.

## Memory Layout and Data Structures

The state representation relies on data density. The basic unit of data is a cell's candidate list. Because Sudoku uses nine digits, a cell's state fits into the lowest nine bits of a `uint16`. The least significant bit (bit 0) represents the digit 1, and bit 8 represents the digit 9. A value of 0 means an unresolvable contradiction. A mask of `0x01FF` indicates all nine digits remain viable candidates.

### The Stack-Allocated Board State

The `Board` struct is the main data structure. It is designed to stay on the execution stack. Passing the `Board` by value avoids garbage collection (GC) overhead, which is important for maintaining throughput in concurrent workloads.

```go
type Board struct {
    Cells    [81]uint16
    RowState [9]uint16
    ColState [9]uint16
    BoxState [9]uint16
    Resolved uint8
}
```

| Field | Bytes | Description |
| :---- | :---- | :---- |
| `Cells` | 162 | The grid representation. Indices 0-80 map to the 9x9 matrix in row-major order. Each element is the candidate mask for that cell. |
| `RowState` | 18 | Bitmasks representing placed digits in each of the 9 rows. Used for O(1) collision detection. |
| `ColState` | 18 | Bitmasks representing placed digits in the 9 columns. |
| `BoxState` | 18 | Bitmasks representing placed digits in the 9 subgrids. |
| `Resolved` | 1 | A counter tracking definitively solved cells (0-81) to allow O(1) completion checks. |
| **Total** | **217** | The struct fits in four 64-byte x86-64 L1 cache lines. |

Separating the `Cells` array from the house constraints (`RowState`, `ColState`, `BoxState`) allows intersection and subset algorithms to use O(1) lookups to check if a digit is placed in a region, avoiding iteration over solved domains.

### Immutable Global Lookup Tables (LUTs)

Mapping a one-dimensional array index to two-dimensional coordinates typically requires division and modulo operations (`x / 9`, `x % 9`). These instructions take multiple CPU cycles. To avoid them, the architecture uses globally immutable Lookup Tables (LUTs) created during initialization.

```go
var (
    RowLUT   [81]uint8
    ColLUT   [81]uint8
    BoxLUT   [81]uint8
    PeersLUT [81][20]uint8
    HouseLUT [27][9]uint8
)
```

| LUT | Description |
| :---- | :---- |
| `RowLUT` | Maps cell index (0-80) to row index (0-8). |
| `ColLUT` | Maps cell index (0-80) to column index (0-8). |
| `BoxLUT` | Maps cell index (0-80) to box index (0-8). |
| `PeersLUT` | Precomputes the 20 intersecting peers for each cell (8 in the row, 8 in the column, 4 in the box). |
| `HouseLUT` | Maps the 27 houses (9 rows, 9 columns, 9 boxes) to their 9 cell indices. |

With `PeersLUT`, the solver propagates constraints by iterating through a 20-element array without boundary checks.

## Hardware-Accelerated Bitwise Primitives

The algorithms rely on bitwise primitives. Go's `math/bits` package provides functions that compile to hardware instructions like `POPCNT` and `TZCNT` on x86-64 architectures. These instructions replace `for` loops for analyzing candidate sets.

Use the following bitwise operations for state manipulation:

1. **Candidate Count:** `bits.OnesCount16(mask)` returns the number of active candidates. If the count is 1, the cell is solved.
2. **Lowest Candidate Isolation:** `mask & -mask` isolates the least significant set bit using two's complement arithmetic, which helps iterate over active candidates without branching.
3. **Candidate Clearing:** `mask &^ target_bit` removes a specific candidate from a mask. This AND NOT operator is the primary way to propagate constraints.
4. **Digit Extraction:** `bits.TrailingZeros16(mask)` returns the zero-based index of the lowest active bit, translating a bit flag (e.g., `000010000`) into an integer (4, mapping to the digit 5).
5. **Subset Check:** Evaluate `(A &^ B) == 0` to check if mask `A` is a subset of mask `B`.

## Group 1: Foundational Eliminations

The simplest SSTS techniques solve cells that either have only one candidate remaining or are the only valid location for a digit in a house. The solver evaluates these continuously.

### Naked Singles

A Naked Single occurs when a cell's bitmask has exactly one active bit, meaning peers have eliminated the other eight digits.

The detection routine loops over the `Cells` array. When `bits.OnesCount16(Cells[i]) == 1`, the cell is resolved. The solver then uses `PeersLUT[i]` to get the 20 peers and clears the candidate bit from them using the `&^` operator. If clearing the bit reduces a peer's candidate count to 1, the solver adds that peer to a stack-based queue to cascade the eliminations without rescanning the board.

### Hidden Singles

A Hidden Single occurs when a digit is an active candidate in only one cell within a house, even if that cell has other candidates. Since the digit has no other valid placement, the cell must take that digit, clearing its other candidates.

To improve execution speed, the solver does not scan for each digit independently. It calculates candidate frequencies within a house simultaneously using bitwise accumulators.

The algorithm uses `seenOnce` and `seenMultiple` variables and iterates over the house's cell masks.

```go
seenOnce, seenMultiple := uint16(0), uint16(0)
for _, cellIdx := range HouseLUT[houseId] {
    mask := board.Cells[cellIdx]
    seenMultiple |= seenOnce & mask // Digits seen previously AND seen now
    seenOnce |= mask                // Register digits seen at least once
}
uniqueMask := seenOnce &^ seenMultiple
```

After the loop, `uniqueMask` contains the bits for any Hidden Singles in the house. If `uniqueMask` is non-zero, a secondary loop finds the specific cell and overwrites its mask to equal the hidden digit. This turns the Hidden Single into a Naked Single, which the propagation queue then broadcasts.

## Group 2: Intersections

When a candidate is confined to a section of a house, it cannot appear in intersecting houses. The SSTS defines two types of Locked Candidates: Type 1 (Pointing) and Type 2 (Box-Line Reduction). These analyze overlap between 3x3 boxes and rows or columns.

### Type 1: Pointing Pairs and Triples

If a candidate in a box only appears in cells sharing a row or column, it points along that line. The digit must go in that line within the box, so it can be cleared from the rest of the line outside the box.

The algorithm loops through the 9 boxes. For each box, it evaluates its three rows and three columns.

1. Create bitmasks for the union of candidates in each row segment of the box ($R_0, R_1, R_2$).
2. Identify digits confined to $R_0$ by computing `exclusive = R_0 &^ (R_1 | R_2)`.
3. If `exclusive > 0`, those digits are locked to that row segment. Find the global row index, iterate over the 6 cells in that row outside the box, and clear the `exclusive` bits.
4. Repeat for the columns ($C_0, C_1, C_2$).

### Type 2: Claiming Pairs and Triples (Box-Line Reduction)

Claiming is the inverse of Pointing. If a candidate in a row or column is confined to a single box, it can be cleared from the other 6 cells in that box.

The algorithm iterates through the 9 rows and 9 columns.

1. Divide the 9 cells of the line into three triplets corresponding to the boxes they intersect. Compute candidate unions $S_0, S_1, S_2$ for these segments.
2. Isolate exclusive digits for the first segment: `exclusive = S_0 &^ (S_1 | S_2)`.
3. If `exclusive > 0`, identify the target box for $S_0$. Iterate over the 6 cells in that box that do not overlap with the line, and apply `&^` to clear the locked candidates.

## Group 3: The Unified Subset Engine

A subset occurs when $N$ candidates are constrained to exactly $N$ cells in a house. The solver finds Pairs ($N=2$), Triples ($N=3$), and Quadruples ($N=4$). To reduce code duplication, Naked and Hidden Subsets use a single engine based on matrix transposition.

### The Mathematical Duality of Subsets

A Naked Subset involves $N$ cells in a house containing exactly $N$ distinct candidates. Other cells in the house cannot contain those $N$ digits. A Hidden Subset involves $N$ digits restricted to exactly $N$ cells in a house. Those $N$ cells cannot contain other digits.

These constructs are mathematical duals. In an incidence matrix where rows are the 9 cells and columns are the 9 digits, a Naked Subset is a combination of rows spanning $N$ columns. A Hidden Subset is a combination of columns spanning $N$ rows.

If a house has $K$ unsolved cells, a Naked Subset of size $N$ guarantees a complementary Hidden Subset of size $K - N$ in the remaining cells. Therefore, a single subset engine operating on bitmasks can find both.

### Architecture of the Transposition Engine

The architecture uses a function: `FindSubsets(masks uint16, N int, mode Mode)`.

1. **Naked Subsets:** Pass the cell candidate bitmasks directly into the `masks` array.
2. **Hidden Subsets:** Transpose the house constraints before calling the function. Create an empty `uint16` array. For each cell $C$ and digit $D$, if $C$ contains $D$, set the $C$-th bit in the $D$-th index of the transposed array. Pass this array to `FindSubsets`.

### Combinatorial Bitwise Evaluation

For sizes $N \in \{2, 3, 4\}$, the engine uses unrolled loops to evaluate subset combinations, avoiding dynamic loops to help with CPU branch prediction.

* **Pairs ($N = 2$):** Tests 36 combinations.

  ```go
  for i := 0; i < 8; i++ {
      for j := i + 1; j < 9; j++ {
          union := masks[i] | masks[j]
          if bits.OnesCount16(union) == 2 { /* Subset Verified */ }
      }
  }
  ```

* **Triples ($N = 3$):** Tests 84 combinations.

  ```go
  union := masks[i] | masks[j] | masks[k]
  if bits.OnesCount16(union) == 3 { /* Triple Verified */ }
  ```

* **Quadruples ($N = 4$):** Tests 126 combinations. The algorithm skips cells with more than 4 active bits (`bits.OnesCount16(masks[i]) > 4`) before entering the nested loops because such a cell cannot be part of a Naked Quad.

After detecting a subset:

* **Naked mode:** Clear the subset's candidates from all other cells in the house.
* **Hidden mode:** Clear all other candidates from the cells containing the subset, turning them into Naked Subsets.

## Group 4: Fish Algorithms

Fish patterns occur when a candidate is constrained to $N$ columns across $N$ rows. The candidate must appear once in each row, filling the $N$ columns. Therefore, the candidate can be cleared from other rows within those $N$ columns. The symmetric logic applies to columns projected across rows.

The solver checks for X-Wing ($N = 2$), Swordfish ($N = 3$), and Jellyfish ($N = 4$).

### Orthogonal Engine

Fish patterns share an algorithmic framework. The architecture uses a single generic detection algorithm for all sizes to reduce code duplication.

**Execution Path:**

1. Iterate over each digit $D$.
2. Build a `uint16` array where index $r$ is the row. Set bit $c$ in the integer at index $r$ if the cell at $(r, c)$ contains candidate $D$.
3. Exclude rows with fewer than 2 active candidates (handled by Hidden Singles) and rows with more than 4 candidates (unlikely to form a Fish).
4. Evaluate combinations of the remaining rows for $N \in \{2, 3, 4\}$.
   * Compute the union of the selected rows: `union = row[i] | row[j]...`.
   * If `bits.OnesCount16(union) == N`, a Fish is present.
   * The active bits in `union` are the target columns. Clear candidate $D$ from those columns in all rows not part of the base $N$ rows.
5. Repeat the process using columns as the base and eliminating from rows.

## Group 5: Wings

These techniques involve functional dependencies across intersecting units using pivot nodes. The XY-Wing is the primary example.

### XY-Wing (Y-Wing)

An XY-Wing involves three cells. The "pivot" cell must have exactly two candidates (e.g., $X$ and $Y$). It must intersect with two "pincer" cells. The first pincer contains $X$ and $Z$, and the second contains $Y$ and $Z$. If the pivot is $X$, the first pincer becomes $Z$. If the pivot is $Y$, the second pincer becomes $Z$. Since the pivot must be $X$ or $Y$, one of the pincers will be $Z$. Therefore, any cell that intersects with both pincers cannot contain $Z$.

**Algorithm:**

1. Find all cells with exactly two candidates (`bits.OnesCount16(mask) == 2`). Store their indices in a `uint8` stack array `bivalueCells`.
2. Iterate over `bivalueCells` to select a Pivot and note its candidates $M_p$.
3. Check the Pivot's 20 peers. A valid Pincer must be in `bivalueCells`, share exactly one candidate with $M_p$, and introduce one new candidate ($Z$).
4. Group pincers sharing the first candidate in array $P_1$, and those sharing the second in $P_2$.
5. Evaluate all pairs from $P_1$ and $P_2$. The two pincers must share the same $Z$ candidate and must not share a house (otherwise it's a Naked Triple).
6. Calculate the intersection of peers between the chosen $P_1$ and $P_2$ cells using a bitwise AND across the precomputed `PeersLUT` sets.
7. Clear the $Z$ candidate bit from any cell in the intersection set.

## Group 6: Graph Coloring

Graph coloring relies on conjugate pairs, which occur when a digit has exactly two valid cells in a house. These cells are mutually exclusive: if one is the digit, the other is not. Linking intersecting conjugate pairs creates chains of alternating states.

Simple Colors processes single continuous chains, while Multi-Colors handles disconnected clusters interacting through common peers. Both use the same graph structure.

### Adjacency List Graph Engine

Dynamic graphs often cause heap allocations and memory fragmentation. To maintain cache locality, this graph engine uses fixed-size integer arrays allocated on the stack.

```go
type Graph struct {
    Nodes     [81]uint8
    Edges     [81][4]uint8
    EdgeCount [81]uint8
    Color     [81]int8
}
```

| Component | Description |
| :---- | :---- |
| `Nodes` | Maps cell indices to Graph Node IDs (0 = inactive). |
| `Edges` | Flat adjacency list. Maximum degree is 4. |
| `EdgeCount` | Number of active edges for a node. |
| `Color` | State enum: 0 = Unvisited, 1 = Color A, -1 = Color B. |

### Breadth-First Search (BFS)

1. For a digit $D$, build the graph by scanning the 27 houses.
2. If a house has exactly two cells containing $D$, add a bidirectional edge between them in `Edges`.
3. Start a stack-allocated Breadth-First Search (BFS) using a `uint8` array as a ring buffer queue.
4. Assign Color `1` to an unvisited origin node. For each neighbor found via BFS, assign the inverse color (`-CurrentColor`).
5. **Color Contradiction:** If the BFS assigns a color to a node that already has the opposite color, the graph is invalid. More commonly, if two nodes with the same color share a peer relationship, that color must be false. Eliminate candidate $D$ from all cells matching the invalid color.
6. **Common Peer Elimination:** Check all uncolored cells containing $D$. If an uncolored cell is a peer of at least one Color 1 node and at least one Color -1 node, it sees both states of the chain. It can never be $D$. Clear candidate $D$ from this cell.

## Deductive Solver Loop

The main deductive engine operates as a priority-driven state machine rather than a linear script. This ensures the solver always applies the simplest, most computationally efficient logic before attempting complex geometric or graph-based techniques.

### Technique Independence and Statelessness

A critical architectural constraint is that every individual solving technique (e.g., `FindHiddenSingles`, `FindXWing`) must be strictly stateless and functionally independent. 

1. **No Assumed Preconditions:** A technique cannot assume that a simpler technique has already run. For example, the `FindNakedPairs` function must operate correctly even if there are unresolved Naked Singles on the board.
2. **Immutable Input:** Each technique receives a copy of the current `Board` state (or reads it safely) and evaluates it. It does not mutate the board directly during evaluation.
3. **Discrete Output:** Instead of modifying state, a successful technique returns a `Step` object containing the specific eliminations or placements discovered.

This isolation guarantees that developers can unit test and benchmark each technique independently by passing it synthetic board states, without needing to mock the entire solver loop or execution history.

### The Execution Priority Queue

The solver evaluates techniques in a strict hierarchy based on their algorithmic complexity:

1. **Tier 0:** Naked Singles ($O(1)$)
2. **Tier 1:** Hidden Singles ($O(1)$)
3. **Tier 2:** Intersections (Pointing/Claiming) ($O(1)$ block evaluation)
4. **Tier 3:** Subsets (Pairs, Triples, Quads) ($O(N)$ combinatorial search)
5. **Tier 4:** Fish (X-Wing, Swordfish, Jellyfish) ($O(R)$ orthogonal search)
6. **Tier 5:** Wings (XY-Wing) ($O(N)$ local graph lookup)
7. **Tier 6:** Coloring (Simple/Multi-Colors) ($O(V + E)$ BFS traversal)

### The Loop and Short-Circuiting

The core loop iterates continuously. It begins at Tier 0.

1. If the current tier finds a valid deduction, the loop records the `Step`, applies the changes to the `Board` (updating candidate masks), and immediately restarts the loop back at Tier 0. This "cascade" ensures that simple deductions uncovered by complex techniques are resolved instantly.
2. If the current tier finds nothing, the loop increments to the next tier.
3. The loop terminates when either the board is fully resolved (`Board.Resolved == 81`), or the engine evaluates all tiers without finding a single deduction (a "stalled" state).

**Hint Mode (`--hint`):**
When the user requests a single hint via the CLI, the solver modifies this loop. Instead of applying the `Step` and restarting, the solver immediately halts execution the moment *any* tier returns a valid `Step`. It returns that single `Step` to the CLI layer to be formatted and printed, bypassing all further evaluation and the brute-force fallback.

### Branching and Escape Analysis

Conditional branching in tight loops disrupts CPU pipelines. The architecture uses branchless logic where possible. For example, instead of a conditional jump:

```go
// Branching
if bitmask & targetMask > 0 {
    count++
}
```

Use bit shifting and arithmetic:

```go
// Branchless
count += uint8((bitmask & targetMask) >> shift_offset)
```

The system passes the 217-byte `Board` struct via pointers (`*Board`) to limit stack copying. By avoiding interfaces and dynamic slices, Go's escape analysis (`go build -gcflags="-m"`) confirms that variables stay off the heap. This keeps the garbage collector dormant, leaving CPU cache available for solving logic.

## Group 7: Brute-Force Fallback (Algorithm X / DLX)

While the solver prioritizes deductive SSTS logic, it requires a fallback mechanism for puzzles with multiple solutions or those that resist implemented deductive techniques. The solver uses Donald Knuth's Algorithm X implemented via Dancing Links (DLX) as a deterministic, exhaustive search of the remaining possibility space.

Like the deductive engines, this DLX implementation uses contiguous arrays rather than pointer-linked heap nodes to maintain cache locality and avoid garbage collection overhead.

### Exact Cover Matrix

A Sudoku grid translates into an exact cover problem. The matrix has constraints (columns) and candidate placements (rows).

1. **Constraints (Columns):** There are 324 exact cover constraints in a standard 9x9 grid:
   * 81 Cell constraints: Each cell must contain exactly one digit.
   * 81 Row constraints: Each row must contain digits 1-9 exactly once.
   * 81 Column constraints: Each column must contain digits 1-9 exactly once.
   * 81 Box constraints: Each box must contain digits 1-9 exactly once.
2. **Placements (Rows):** There are up to 729 possible placements (9 cells *9 rows* 9 digits). Each valid placement satisfies exactly 4 constraints.

### Array-Based Dancing Links

Instead of creating node structs linked by `Left`, `Right`, `Up`, and `Down` pointers, the solver uses parallel arrays indexed by an integer Node ID.

```go
type DLX struct {
    Left   []int
    Right  []int
    Up     []int
    Down   []int
    Col    []int
    Row    []int
    Count  []int // For column headers: number of nodes in the column
    Header int   // Root node ID
}
```

By preallocating these arrays to the maximum possible number of nodes (around 3000 integers total, fitting in a few kilobytes), the entire DLX matrix resides in a contiguous block of memory.

### Fallback Execution

When the deductive engine exhausts all available techniques without finding a solution, the solver triggers the fallback:

1. **Matrix Construction:** The solver builds the DLX matrix based on the *current* state of the `Board`'s candidate masks, not the initial puzzle state. This takes advantage of the deductive eliminations already performed, significantly reducing the size of the exact cover problem.
2. **Algorithm X Search:** The DLX engine recursively selects the column with the fewest nodes (to branch as little as possible) and covers it, along with all intersecting rows and columns.
3. **State Restoration:** If a branch fails, DLX efficiently uncovers the nodes, restoring the previous state in $O(1)$ time per node using the parallel array structure.
4. **Resolution:** Upon finding a valid exact cover, the solver translates the selected rows back into digit placements, updating the `Board` state. This entire process is recorded as a single "Brute Force" step in the output log.

## Command-Line Interface (CLI) Architecture

The application uses the `github.com/spf13/cobra` library to manage its command-line interface. The CLI structure consists of three subcommands that route to distinct functional packages.

1. **`solve`:** This command parses the input, initializes the `Board`, and runs the deductive engine in a loop. If the deductive techniques finish, it checks the `Board.Resolved` counter. If the count is less than 81, it invokes the DLX fallback. The command prints the initial board, the step log, and the final board. If the `--hint` flag is active, it halts after executing and printing the first successful deductive step.
2. **`check`:** This command evaluates a puzzle's validity. It bypasses the step logger and runs the DLX engine directly to count the total number of valid solutions. It reports "valid" for exactly one solution, "multiple solutions" for more than one, and "no solutions" for zero.
3. **`print`:** This command parses the input and immediately formats and prints the grid based on the `--format` flag, bypassing the solving engine entirely.

## Input Parsing Engine

The input parser converts text (from standard input, files, or positional arguments) into the internal `Board` structure.

1. **Standardization:** The parser reads the input line by line. It strips comments starting with `#`, removes leading and trailing whitespace, and drops empty lines and horizontal borders.
2. **Format Detection:** If the remaining text has exactly 9 lines with 9 space- or pipe-separated groups per line, the parser treats it as a candidate grid. It counts the number of digits in each group; groups with one digit are treated as givens, and groups with multiple digits are treated as explicit pencil marks. Otherwise, the parser removes all spaces and pipes, expecting exactly an 81-character string.
3. **Tokenization:** The parser scans the 81-character string. The first non-digit character it encounters becomes the empty cell marker. It maps digits 1-9 to their bitwise representations and sets empty cells to the full `0x01FF` candidate mask.

## Step Logging and Explanation System

The deductive engine needs to record its actions to provide a human-readable log. Since the `Board` struct is optimized for calculation and is passed by value, the logging mechanism operates externally.

The solver function maintains a slice of `Step` structs.

```go
type Step struct {
    Technique string
    Target    []int    // Cell indices affected
    Values    []int    // Digits involved
    Action    ActionType // Placement or Elimination
}
```

When a deductive function finds a valid pattern, it returns a `Step` object containing the raw cell indices and digits. The main solver loop appends this to the log and applies the changes to the `Board`.

Before output, a dedicated formatter translates the raw array indices into the `rncn` notation (e.g., mapping index 14 to `r2c6`). It groups continuous sequences into compact references (like `r1c36` instead of `r1c3, r1c6`).

## Output Formatting and Rendering

The formatter package translates the internal `Board.Cells` array into the required text formats.

For simple formats like `raw` or `9x9`, it iterates over the cells and outputs the solved digits or the empty cell marker. For complex formats like `ss` or `pretty`, it injects border characters based on the row and column indices using modulo logic.

### Dynamic Rendering for Pencil Marks (`pm`)

The `pm` format displays remaining candidates and requires dynamic column widths.

1. The formatter iterates through the `Board` and calculates the string length of the candidates for each cell (e.g., `bits.OnesCount16(mask)`), and tracks the maximum string length encountered.
2. When printing the grid, it pads each cell's output string with spaces to match the maximum width determined in step 1, ensuring the grid remains aligned and all boxes have a consistent width.

### ANSI Colorization

For terminal output, the `pretty` and `prettypm` formats support ANSI color codes using the `github.com/fatih/color` library. The parser records which cells were populated in the initial input. The formatter references this record during output, utilizing `fatih/color` to wrap given digits in a specific ANSI color code (e.g., white) and placed digits or candidate values in a contrasting color (e.g., green or gray).
.
