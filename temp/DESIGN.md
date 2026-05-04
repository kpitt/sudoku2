# High-Performance Architectural Design for Simple Sudoku Technique Set (SSTS) Algorithms

## Introduction to Deductive Sudoku Resolution and Systems Architecture

The engineering of an optimal logic-puzzle solver demands a meticulous equilibrium between combinatorial algorithm design and low-level systems architecture. Within the computational domain of Sudoku, solving engines must rapidly traverse a vast state space while executing complex deductive constraints. While stochastic search methods, simulated annealing, and brute-force backtracking algorithms—most notably Donald Knuth’s Algorithm X implemented via Dancing Links (DLX)—are highly effective for general exact cover problems, they fundamentally rely on speculative trial and error rather than pure deductive reasoning. In contrast, the "Simple Sudoku Technique Set" (SSTS) establishes a standardized framework of deterministic, deductive strategies utilized by human experts and advanced heuristic solvers to systematically reduce the puzzle state. The SSTS defines a progression of techniques ranging from foundational single-candidate deductions to highly complex multi-point graph-coloring algorithms, such as XY-Wings and Multi-Colors, serving as an industry benchmark for solver capability.

This architectural design specification details a high-performance, zero-allocation software engine implemented in Go (Golang) engineered specifically to resolve the complete SSTS. The primary objective of this architecture is to maximize computational throughput by aggressively exploiting modern processor topologies. This is achieved by enforcing strict CPU cache locality, eliminating branch mispredictions, leveraging hardware-accelerated bitwise intrinsic instructions, and completely eschewing heap allocations in favor of static stack memory.

Traditional implementations of advanced Sudoku solvers, including DLX, depend on pointer-linked data structures instantiated across the memory heap. While mathematically elegant, this object-oriented approach severely degrades CPU cache performance and introduces false sharing in multithreaded environments because threads continually invalidate each other's cache lines when traversing scattered memory addresses. Conversely, the architecture defined in this specification utilizes compact, contiguous, data-oriented bitmask arrays. This paradigm ensures that the entire board state and constraint matrices fit entirely within a minimal number of L1 CPU cache lines, effectively neutralizing main memory access latency. For the implementation engineer, this document serves as a comprehensive blueprint for translating abstract SSTS logic into optimally performing Go routines, prioritizing shared functionality across isomorphic mathematical properties to maintain a lean, robust codebase.

## Architectural Core: Memory Layout and Data Structures

To achieve peak computational efficiency, the representation of the puzzle state must prioritize absolute data density. The fundamental quantum of information in a Sudoku solver is the candidate list for any given cell. Because the game strictly utilizes nine digits, the candidate state of a cell can be perfectly encoded within the lowest nine bits of a standard 16-bit unsigned integer (`uint16`). In this encoding scheme, the least significant bit (bit 0) represents the digit 1, and bit 8 represents the digit 9. A value of zero denotes an unresolvable contradiction, while a fully populated mask of `0x01FF` (binary `111111111`) indicates that all nine digits remain mathematically viable candidates.

### The Stack-Allocated Board State

The central data structure of the solving engine is the `Board` object. This structure is engineered to exist entirely on the execution stack. By passing the `Board` structure by value during recursive evaluations or functional transformations, the application guarantees that no garbage collection (GC) overhead is incurred. Eliminating garbage collection pressure is a critical paradigm in high-performance Go engineering, as GC pauses can significantly degrade throughput in heavily parallelized workloads.

| Field Identifier | Go Data Type | Size (Bytes) | Architectural Purpose and Operational Description |
| :---- | :---- | :---- | :---- |
| `Cells` | `uint16` | 162 | The canonical representation of the grid. Indices 0 through 80 map directly to the 9x9 matrix in row-major order. Each element holds the active candidate mask for that specific spatial coordinate. |
| `Resolved` | `uint8` | 1 | A highly optimized counter tracking the absolute number of definitively solved cells (ranging from 0 to 81). This permits $O(1)$ verification of puzzle completion without iterating over the array. |
| `RowState` | `uint16` | 18 | Aggregated bitmasks representing the digits that have been permanently placed in each of the 9 horizontal rows. Facilitates $O(1)$ collision detection during candidate propagation. |
| `ColState` | `uint16` | 18 | Aggregated bitmasks representing the permanently placed digits across the 9 vertical columns. |
| `BoxState` | `uint16` | 18 | Aggregated bitmasks for the 9 distinct 9x9 subgrids. |
| **Total Memory Footprint** | **Struct** | **217** | The structure fits comfortably within four standard 64-byte x86-64 L1 cache lines, ensuring instantaneous retrieval by the arithmetic logic unit (ALU). |

The explicit decoupling of the isolated `Cells` array from the overarching house constraints (`RowState`, `ColState`, `BoxState`) ensures that intersection and subset techniques can perform localized constant-time $O(1)$ lookups to ascertain whether a target digit is already saturated within a particular region. This decoupling directly prevents redundant iteration over already solved domains.

### Immutable Global Lookup Tables (LUTs)

In conventional software development, mapping a one-dimensional array index back to two-dimensional or three-dimensional spatial coordinates requires division and modulo arithmetic (`x / 9`, `x % 9`). However, division instructions consume excessive CPU cycles and introduce pipeline latency. To guarantee branchless, zero-latency coordinate resolution, the architecture mandates the use of globally immutable Lookup Tables (LUTs) mapped statically into the application's read-only memory segment during initialization.

| LUT Identifier | Data Type | Dimensionality Mapping |
| :---- | :---- | :---- |
| `RowLUT` | `uint8` | Instantly maps a continuous cell index (0-80) to its corresponding row index (0-8). |
| `ColLUT` | `uint8` | Instantly maps a continuous cell index (0-80) to its corresponding column index (0-8). |
| `BoxLUT` | `uint8` | Maps a continuous cell index to its corresponding 3x3 box index (0-8). |
| `PeersLUT` | `uint8` | Precomputes the 20 distinct intersecting peers for every individual cell (comprising 8 cells in the shared row, 8 in the shared column, and 4 in the shared box that have not already been counted). |
| `HouseLUT` | `uint8` | Maps the 27 structural houses (9 rows, 9 columns, 9 boxes) to their constituent 9 cell indices. |

By utilizing `PeersLUT`, the engine can execute candidate propagation across an entire intersecting domain by simply iterating through a flat 20-element integer array, avoiding any conditional boundary checks.

## Hardware-Accelerated Bitwise Primitives

The mathematical foundation of all SSTS algorithms relies entirely on core bitwise primitives. The Go language standard library provides the `math/bits` package, which contains intrinsic functions that compile directly down to dedicated hardware instructions on modern architectures (such as `POPCNT` for population counting and `TZCNT` for trailing zero detection on x86-64 processors). This hardware alignment completely eliminates the requirement for software-level `for` loops when analyzing candidate sets.

The implementation engineer must strictly adhere to the following primitive operations for state manipulation:

1. **Candidate Quantification:** `bits.OnesCount16(mask)` returns the exact number of active candidates within a cell. If the population count evaluates to 1, the cell is unambiguously solved.
2. **Lowest Candidate Isolation:** The bitwise operation `mask & -mask` mathematically isolates the least significant set bit. This is a property of two's complement arithmetic and is paramount for iterating over active candidates sequentially without engaging branch instructions.
3. **Candidate Eradication:** `mask &^ target_bit` clears a specific candidate from a mask without affecting any other state. This operator (AND NOT) is the primary mechanism for constraint propagation.
4. **Digit Extraction:** `bits.TrailingZeros16(mask)` returns the zero-based index of the lowest active bit, instantly translating a binary flag (e.g., `000010000`) into a workable integer index (4, mathematically mapping to the Sudoku digit 5).
5. **Subset Validation:** To check if a candidate mask `A` is a strict subset of mask `B` the engine evaluates `(A &^ B) == 0`.

## Group 1: Singularities and Foundational Eliminations

The most primitive techniques codified within the SSTS involve finalizing cells that have either been stripped down to a solitary candidate through peer constraints, or are the sole remaining valid geographic location for a specific digit within a predefined structural house. These techniques must be evaluated continuously.

### Naked Singles (Full House and Obvious Singles)

A Naked Single—often referred to as an Obvious Single or a Full House when completing a unit—materializes when a cell's bitmask contains exactly one active bit. This indicates that all other eight digits have been successfully eliminated by intersecting peers.

The detection routine iterates linearly across the `Cells` array. For any index where `bits.OnesCount16(Cells[i]) == 1`, the cell transitions from an unresolved state to a resolved state. The crucial subsequent action is propagation. The engine utilizes `PeersLUT[i]` to extract the 20 intersecting peers. The isolated candidate bit is subsequently stripped from all 20 peers using the `&^` bitwise operator. If any peer's candidate count is reduced to 1 during this stripping process, that peer is immediately pushed onto a stack-based resolution queue, cascading the eliminations without requiring a full board rescan.

### Hidden Singles (Cross-Hatching)

A Hidden Single manifests when a specific digit appears as an active candidate in exactly one cell across an entire house (row, column, or box), despite that specific cell possessing multiple other active candidates. Because the digit has no other valid placement within that spatial jurisdiction, the cell must be assigned that digit, implicitly obliterating all other candidates within the cell.

To optimize cache utilization and execution speed, the architecture prohibits scanning cell-by-cell for each of the nine digits independently. Instead, the engine computes the exact frequencies of all nine candidates within a house simultaneously by employing parallel bitwise accumulators.

The algorithm maintains two state tracking variables: `seenOnce` and `seenMultiple`. The engine iterates through the nine cell masks of the target house.

```go
seenOnce, seenMultiple := uint16(0), uint16(0)
for _, cellIdx := range HouseLUT[houseId] {
    mask := board.Cells[cellIdx]
    seenMultiple |= seenOnce & mask // Digits seen previously AND seen now
    seenOnce |= mask                // Register digits seen at least once
}
uniqueMask := seenOnce &^ seenMultiple
```

Upon loop completion, the `uniqueMask` variable holds the exact bit(s) corresponding to the Hidden Single(s) within that house. If `uniqueMask` is non-zero, a secondary fast-path loop isolates the exact cell containing that unique bit, and aggressively overwrites the cell's entire candidate mask to strictly equal the hidden digit mask. This operation instantaneously converts the Hidden Single into a Naked Single, allowing the standard propagation queue to broadcast the newly resolved digit.

## Group 2: Intersections (Locked Candidates)

When a specific candidate is spatially confined to a sub-segment of a house, its presence logically negates that candidate from intersecting houses. The SSTS defines two complementary types of Locked Candidates: Type 1 (Pointing Pairs/Triples) and Type 2 (Box-Line Reduction). These techniques analyze the structural overlap between 3x3 boxes and linear rows/columns.

### Type 1: Pointing Pairs and Triples

If a candidate within a 3x3 box only appears within cells that share a single global row or global column, the candidate is said to "point" along that line. Because the digit must be placed *somewhere* within the box, and all available options dictate it lands on that specific line, the candidate can safely be eliminated from all other cells situated on that line outside the boundaries of the originating box.

The implementation algorithm iterates through each of the 9 boxes. For each box, it isolates its three constituent rows and three constituent columns.

1. The engine constructs three bitmasks representing the union (bitwise OR) of candidates within each of the three rows of the box. Let these be $R_0, R_1, R_2$.
2. To mathematically identify digits exclusively confined to $R_0$, the engine computes: `exclusive = R_0 &^ (R_1 | R_2)`.
3. If the `exclusive` mask evaluates to a value greater than zero, those specific digits are locked to that row segment. The algorithm extracts the global row index, iterates over the 6 cells of that global row located entirely outside the box, and clears the `exclusive` bits from those cells.
4. This entire procedure is subsequently mirrored for the columns ($C_0, C_1, C_2$).

### Type 2: Box-Line Reduction

Box-Line Reduction operates as the exact mathematical inverse of the Pointing technique. If a candidate within a global row or column is entirely confined to the segment that intersects a single 3x3 box, it can be unilaterally eliminated from the remaining 6 cells of that box.

The algorithm iterates through each of the 9 rows and 9 columns.

1. The 9 cells of the linear house are divided into three triplets, each corresponding to the structural box they intersect. The engine computes the candidate unions $S_0, S_1, S_2$ for these segments.
2. The exclusive digits for the first segment are isolated: `exclusive = S_0 &^ (S_1 | S_2)`.
3. If `exclusive > 0`, the engine identifies the target box associated with $S_0$. It then iterates over the 6 cells within that box that do *not* overlap with the linear house, applying the `&^` operator to strip the locked candidates.

By leveraging bitwise block evaluations, the system executes these spatial intersection checks without maintaining any complex spatial memory, relying solely on sequential integer arrays.

## Group 3: The Unified Subset Engine (Naked and Hidden)

Subsets represent a critical threshold in constraint programming. A subset occurs when a group of $N$ candidates is strictly constrained to exactly $N$ cells within a shared house. The SSTS strictures mandate the comprehensive resolution of Pairs ($N = 2$), Triples ($N = 3$), and Quadruples ($N = 4$). To satisfy the architectural directive of maximizing shared functional logic and minimizing executable bloat, Naked and Hidden Subsets are evaluated using a singular, mathematically unified engine based on boolean matrix transposition.

### The Mathematical Duality of Subsets

A Naked Subset is defined as $N$ cells in a house that collectively contain exactly $N$ distinct candidate digits. If this condition holds, no other cell in that house may contain those $N$ digits. Conversely, a Hidden Subset consists of $N$ distinct digits whose valid placements in a house are restricted to exactly the same $N$ cells. In this scenario, those $N$ cells cannot contain any *other* digits.

Graph theoretically and algebraically, these two constructs are exact duals of one another. If a solver constructs an incidence matrix where the rows represent the 9 spatial cells of a house and the columns represent the 9 possible candidate digits, a Naked Subset manifests as a combination of matrix rows spanning a limited number of columns. A Hidden Subset manifests as a combination of matrix columns spanning a limited number of rows.

Furthermore, a fundamental theorem of exact cover spaces dictates that if a house contains $K$ unsolved cells, the presence of a Naked Subset of size $N$ absolutely guarantees the simultaneous presence of a complementary Hidden Subset of size $K - N$ existing within the remaining empty cells. Therefore, writing disparate algorithmic routines for Naked and Hidden Subsets is computationally redundant. A singular, highly optimized combinatorial subset engine operating over generic integer bitmasks can definitively discover both classes of techniques.

### Architecture of the Transposition Engine

The architecture defines a generalized, high-speed function: `FindSubsets(masks uint16, N int, mode Mode)`.

1. **Executing Naked Subsets:** The input `masks` array is populated directly with the standard candidate bitmasks of the cells comprising the house.
2. **Executing Hidden Subsets:** Before invoking the engine, the system performs a matrix transpose on the house constraints. It instantiates an empty `uint16` array. For each cell $C$ (from 0 to 8) and each digit $D$ (from 1 to 9), if cell $C$ contains digit $D$, the $C$-th bit is set in the $D$-th index of the transposed array. This transposed representation is then seamlessly fed into the exact same `FindSubsets` engine.

### Combinatorial Bitwise Evaluation

For target sizes $N \in \{2, 3, 4\}$, the combinatorial engine utilizes heavily unrolled, statically defined loops to evaluate all possible subset combinations, intentionally bypassing dynamic loop generation to ensure perfect branch prediction on the CPU.

* **Pairs ($N = 2$):** The algorithm requires testing all combinations of 2 digits: $C(9, 2) = 36$ combinations.

  ```go
  for i := 0; i < 8; i++ {
      for j := i + 1; j < 9; j++ {
          union := masks[i] | masks[j]
          if bits.OnesCount16(union) == 2 { /* Subset Verified */ }
      }
  }
  ```

* **Triples ($N = 3$):** Requires evaluating all combinations of 3 digits: $C(9, 3) = 84$ combinations.

  ```go
  union := masks[i] | masks[j] | masks[k]
  if bits.OnesCount16(union) == 3 { /* Triple Verified */ }
  ```

* **Quadruples ($N = 4$):** Evaluates all combinations of 4 digits: $C(9, 4) = 126$ combinations. To mitigate the processing time of the quadruple loop, an aggressive heuristic pruning strategy is implemented. Before entering the deepest nested loops, the system checks if the initial candidate cell possesses more than 4 active bits (`bits.OnesCount16(masks[i]) > 4`). If true, the cell mathematically cannot participate in a Naked Quad, and the loop advances immediately, preventing hundreds of redundant ALU operations.

Upon successful detection, the application context dictates the cleanup phase. If operating in Naked mode, the identified candidate bits are explicitly cleared from the masks of all *other* cells in the original house. If operating in Hidden mode, all *other* candidate bits are explicitly cleared from the original spatial cells representing the subset, effectively compressing them into Naked entities.

## Group 4: Fish Architectures (Orthogonal Projections)

Fish algorithms elevate the principle of locked candidates from isolated spatial houses to multiple parallel structural vectors. When a single specific candidate digit is constrained to exactly $N$ columns distributed across $N$ distinct rows, it establishes a geometric Fish pattern. Because the rules of Sudoku dictate that the candidate must appear exactly once in each of those $N$ rows, and there are exactly $N$ available columns among them, those $N$ columns are completely saturated by the candidate. Consequently, the candidate can be securely eliminated from all other rows existing within those $N$ target columns. The inverse geometric logic (Rows projected across Columns) is mathematically symmetric and equally valid.

The SSTS requires exhaustive detection of the X-Wing ($N = 2$), the Swordfish ($N = 3$), and the elusive Jellyfish ($N = 4$).

### The Unified Orthogonal Detection Engine

Similar to the subset architecture, Fish patterns of all scales share an identical underlying algorithmic and geometric framework. Writing heavily customized functions to individually handle X-Wings, Swordfish, and Jellyfish introduces severe instruction cache bloat and duplicates logic. The architecture specifies a singular, highly generic Fish detection algorithm that processes a specific target digit across orthogonal grid dimensions.

**Algorithmic Execution Path:**

1. The overarching engine iterates over each individual digit $D$.
2. The engine dynamically constructs a dense `uint16` array where the index $r$ represents the structural row $r$. The integer value assigned to index $r$ is a specialized bitmask where bit $c$ is flagged as 1 if the spatial cell at coordinates $(r, c)$ contains the active candidate $D$.
3. **Heuristic Filtration:** Before combinatorial evaluation begins, the engine excludes any row where `bits.OnesCount16(mask) < 2` (since a single candidate is already handled by the Hidden Singles routine) and any row where `bits.OnesCount16(mask) > 4` (as such rows are statistically improbable to formulate a pure Jellyfish without prior simplification from lesser techniques). This step drastically culls the search space.
4. **Combinatorial Search:** The algorithm evaluates combinations of the surviving rows for sizes $N \in \{2, 3, 4\}$.
   * For $N$ selected base rows, the engine computes the bitwise OR union of their specialized masks: `union = row[i] | row[j]...`.
   * If `bits.OnesCount16(union) == N`, a valid Fish structure has been formally identified.
   * The bits active within the `union` variable identify the specific target columns. The algorithm then iterates through all 9 rows of the grid; if an evaluated row is *not* one of the foundational $N$ base rows, it applies the `&^` operator to strip the candidate $D$ from the identified target columns within that row.
5. **Orthogonal Transpose Execution:** Following row evaluation, the entire sequence is repeated employing an orthogonal mapping schema (where columns act as the combinatorial base vectors, and the resulting target bitmasks dictate row eliminations).

By strictly bounding $N$ dynamically through an outer loop iterating from 2 to 4, branch prediction mechanisms remain optimized, and instruction cache lines remain tightly cohesive. The transition from backtracking tree search to combinatorial bitwise arithmetic transforms an operation with $O(R^N)$ algorithmic complexity into a microsecond-scale bitwise pipeline.

## Group 5: Wings and Interlocking Logic Gates

The highest echelon of techniques categorized within the SSTS involves multi-cell functional dependencies that cannot be resolved within simple orthogonal dimensions. These paradigms require the engine to track and evaluate state across intersecting units utilizing bipartite pivot nodes and extended logical chains. The foundational technique of this category is the XY-Wing.

### XY-Wing (Y-Wing) Bipartite Resolution

An XY-Wing functions analogously to a digital logic gate. It materializes when three specific cells form an interdependent web. The primary cell, termed the "pivot", must contain exactly two candidates (e.g., $X$ and $Y$). This pivot must structurally intersect with two peripheral cells, termed "pincers".

The first pincer must contain the candidates $X$ and $Z$. The second pincer must contain the candidates $Y$ and $Z$. The deductive logic operates as follows: If the pivot resolves to $X$, the first pincer is forcefully collapsed to $Z$. Conversely, if the pivot resolves to $Y$, the second pincer is forcefully collapsed to $Z$. Because the pivot is strictly bivalue (it *must* be either $X$ or $Y$), mathematical certainty dictates that at least one of the two pincers *must* eventually resolve to the digit $Z$. Therefore, any cell that exists within the intersection of *both* pincers' spatial jurisdictions cannot possibly contain the candidate $Z$, as it would invalidate the fundamental constraints of the board.

**High-Performance Bitwise Graph Search:**

1. The engine rapidly scans the board to identify all strictly bivalue cells (`bits.OnesCount16(mask) == 2`). The indices of these cells are stored sequentially in a fixed-size `uint8` stack array `bivalueCells`.
2. The algorithm iterates over `bivalueCells` to tentatively select a Pivot. It logs its candidate mask $M_p$.
3. The engine scans the 20 predefined peers of the selected Pivot (accessed via `PeersLUT`). To qualify as a valid Pincer, a peer must also exist within the `bivalueCells` array, must share exactly one candidate bit with $M_p$, and must introduce exactly one novel candidate bit (the $Z$ candidate).
4. Valid pincers sharing the first candidate of the Pivot are stored in an array $P_1$, and those sharing the second candidate are stored in $P_2$.
5. The engine iterates over all possible pairs generated from $P_1$ and $P_2$. It verifies that the chosen pincers share the exact same $Z$ candidate mask. It also verifies that the two pincers do *not* share a common house (if they did, the configuration would devolve into a simpler Naked Triple).
6. Upon validation, the engine computes the intersection of peers between the chosen $P_1$ and $P_2$ nodes. Because peer sets are statically precomputed in the global LUTs, this peer intersection is calculated as a fast bitwise AND operation across an array of three `uint32` integers (acting mathematically as an 81-bit boolean set).
7. For any spatial cell flagged as active within the resulting intersection set, the engine applies `&^ Z_mask` to permanently clear the $Z$ candidate bit.

## Group 6: Graph Coloring Networks

Coloring methodologies represent the zenith of the SSTS, exploiting the strict binary nature of conjugate pairs to trace rippling boolean implications across the grid. A conjugate pair exists when a specific digit is isolated to exactly two viable cells within any house. Due to the fundamental rules of Sudoku, these two cells represent mutually exclusive boolean states: if one evaluates to true, the other must evaluate to false. By linking intersecting conjugate pairs, the system generates an extensive chain of alternating boolean realities.

The SSTS defines two coloring tiers: Simple Colors and Multi-Colors. Simple Colors evaluates single continuous contradiction chains, while Multi-Colors expands this logic across disconnected graph clusters that interact via common peers. Both techniques operate upon the exact same underlying topological graph structure.

### Zero-Allocation Adjacency List Graph Engine

Constructing and traversing dynamic graph structures in high-level managed languages typically invokes significant heap allocations (via `new` or `make` keywords), generating severe memory fragmentation and GC churn. To preserve rigorous cache locality, the graph engine is materialized utilizing fixed-size integer arrays instantiated entirely on the function stack.

| Data Structure Component | Go Type Definition | Structural Description |
| :---- | :---- | :---- |
| `Nodes` | `uint8` | Maps raw spatial cell indices to internally tracked Graph Node IDs (0 = unmapped/inactive). |
| `Edges` | `uint8` | A flat adjacency list. The absolute maximum degree of connectivity for a node in a single-digit conjugate graph is 4, bounding the array dimension. |
| `EdgeCount` | `uint8` | An accumulator tracking the number of active, verified edges terminating at a given node. |
| `Color` | `int8` | State tracking enum: `0` = Unvisited, `1` = Color A (Reality 1), `-1` = Color B (Reality 2). |

### Algorithmic Breadth-First Search (BFS) Execution

1. For a given target digit $D$, the engine constructs the logical graph by systematically scanning all 27 houses (rows, columns, boxes).
2. If a house contains exactly two unresolved cells actively holding candidate $D$, the system registers a bidirectional conjugate edge between them within the `Edges` array.
3. The algorithm initializes a stack-allocated Breadth-First Search (BFS) utilizing a fixed `uint8` ring buffer array `queue` to traverse the topology.
4. An unvisited origin node is seeded with Color `1`. For each sequentially connected neighbor discovered via BFS, the engine forcefully assigns the inverse logical color (`-CurrentColor`).
5. **Elimination Rule 1 (Color Contradiction):** If the BFS traversal attempts to assign a color to a node that has already been painted with the *opposite* color, the internal topology of the graph is demonstrably invalid, meaning the origin assumption is fatally flawed. More commonly in Simple Colors, if the engine detects that two distinct nodes painted with the *same* color share a geographic peer relation, that specific color represents an impossibility. The engine immediately flags Color 1 (or -1) as false and eliminates candidate $D$ from all cells painted with the contradicted color.
6. **Elimination Rule 2 (Common Peer Elimination):** The engine systematically iterates over all uncolored cells still containing the candidate $D$. It verifies their peer network. If an uncolored cell is identified as a direct peer of *at least one Color 1 node* AND *at least one Color -1 node*, it is logically positioned to "see" both potential realities of the conjugate chain. Regardless of which boolean state ultimately collapses into truth, that uncolored cell can never house digit $D$. Candidate $D$ is purged from the uncolored cell.

By strictly enforcing the adjacency list's stack memory footprint, populating and traversing the entire graph topology requires fractions of a microsecond per digit, rendering complex Multi-Color logic chaining computationally trivial.

## System Orchestration, Execution Heuristics, and Compiler Directives

The individual algorithms detailed throughout this specification exhibit vastly divergent computational complexities. Consequently, the operational sequence in which they are invoked fundamentally dictates the system's overall throughput. An optimal deductive solving engine must not operate as a linear script, but rather as a highly responsive, priority-driven state machine.

### Execution Hierarchy and Escalation

Lower-complexity, high-yield techniques must execute continuously. Only when an entire subset of techniques fails to mutate the overall board state does the engine escalate to a mathematically heavier tier. Crucially, if a higher-tier algorithmic technique successfully strips even a single candidate bit from the board, the execution state immediately interrupts and resets to the absolute lowest tier (Naked Singles) to rapidly cascade the newly simplified logic.

| Priority Tier | SSTS Technique Classification | Algorithmic Complexity | Execution Frequency |
| :---- | :---- | :---- | :---- |
| **0** | Naked Singles / Full Houses | $O(1)$ constant ALU ops per cell | Extremely High (Primary loop) |
| **1** | Hidden Singles (Cross-Hatching) | $O(1)$ via bitwise accumulators | High |
| **2** | Intersections (Type 1 & 2) | $O(1)$ sequential block evaluation | Medium |
| **3** | Subsets (Naked/Hidden Pairs to Quads) | $O(N_4)$ dynamically constrained subset | Low-Medium |
| **4** | Fish (X-Wing to Jellyfish) | $O(R_4)$ dynamically constrained search | Low |
| **5** | Wings & Intersections (XY-Wing) | $O(N_2)$ localized graph lookup | Very Low |
| **6** | Graph Coloring (Simple/Multi-Colors) | $O(V + E)$ BFS stack traversal | Last Resort |

### Branch Predictability and Compiler Escape Analysis

In the Go runtime environment, conditional branching (`if/else` statements) located inside tightly nested inner loops severely disrupts the CPU's instruction pipeline, leading to costly stall cycles. To effectively mitigate this, the architecture advocates for branchless mathematical equivalencies wherever logically feasible. For example, instead of relying on a conditional jump to increment a counter:

```go
// Suboptimal branching code
if bitmask & targetMask > 0 {
    count++
}
```

The architecture employs bit shifting and arithmetic accumulation:

```go
// Optimal branchless implementation
count += uint8((bitmask & targetMask) >> shift_offset)
```

Furthermore, by strictly passing the 217-byte `Board` struct via explicit pointers (`*Board`) within the tightly bound local execution package, the system ensures negligible stack copying. Because the entire constraint resolution matrix is self-contained and explicitly avoids generating polymorphic interfaces or dynamic slices, the Go compiler's escape analysis (verifiable via `go build -gcflags="-m"`) mathematically guarantees that zero variables migrate to the heap. This zero-allocation design ensures that the garbage collector remains entirely dormant during massive inference bursts, permitting the processor to dedicate 100% of L1 and L2 cache capacity directly to the bitwise solving logic.

## Architectural Synthesis

The high-performance resolution of the complete Simple Sudoku Technique Set fundamentally shifts the engineering bottleneck from abstract algorithmic logic to concrete memory architecture. By explicitly rejecting object-oriented, pointer-centric matrix models in favor of tightly packed, data-oriented bitmask arrays, the system natively aligns with the physical realities of modern superscalar CPU cache architectures.

Furthermore, the mathematical unification of theoretically dual algorithms—such as deploying isolated matrix transpositions to solve both Naked and Hidden subsets simultaneously, and architecting a single contiguous computational routine to execute all geometric variations of Fish patterns—drastically reduces cyclomatic complexity and instruction cache misses. This architectural blueprint ultimately yields a deductive inference engine capable of executing millions of advanced technique cycles per second, ensuring exhaustive problem-space reduction with absolute maximal computational efficiency.

#### **Works cited**

1. Solving Sudoku in Go: Backtracking, Constraint Caching, and Benchmarks \- Manuel Fedele, accessed May 4, 2026, [https://manuelfedele.github.io/posts/build-a-sudoku-solver-in-golang/](https://manuelfedele.github.io/posts/build-a-sudoku-solver-in-golang/)
2. Sudoku Solver — From a Father-Son Coding Challenge to a CPU Benchmark Tool, accessed May 4, 2026, [https://allenkuo.medium.com/sudoku-solver-from-a-father-son-coding-challenge-to-a-cpu-benchmark-tool-38300c968bf5](https://allenkuo.medium.com/sudoku-solver-from-a-father-son-coding-challenge-to-a-cpu-benchmark-tool-38300c968bf5)
3. Sudoku solving algorithms \- Wikipedia, accessed May 4, 2026, [https://en.wikipedia.org/wiki/Sudoku\_solving\_algorithms](https://en.wikipedia.org/wiki/Sudoku_solving_algorithms)
4. User Manual (Chapter 3: Configuring the Solver) \- HoDoKu, accessed May 4, 2026, [https://hodoku.sourceforge.net/en/docs\_solv.php](https://hodoku.sourceforge.net/en/docs_solv.php)
5. SSTS \- Sudopedia, accessed May 4, 2026, [https://www.sudopedia.org/wiki/SSTS](https://www.sudopedia.org/wiki/SSTS)
6. User Manual (Chapter 4: Creating Sudokus) \- HoDoKu, accessed May 4, 2026, [https://hodoku.sourceforge.net/en/docs\_cre.php](https://hodoku.sourceforge.net/en/docs_cre.php)
7. Sudoku, Go and WebAssembly \- Eli Bendersky's website, accessed May 4, 2026, [https://eli.thegreenplace.net/2022/sudoku-go-and-webassembly/](https://eli.thegreenplace.net/2022/sudoku-go-and-webassembly/)
8. What would be the fastest (performance wise) way to implement a sudoku-like possibilities set? : r/rust \- Reddit, accessed May 4, 2026, [https://www.reddit.com/r/rust/comments/1ekwori/what\_would\_be\_the\_fastest\_performance\_wise\_way\_to/](https://www.reddit.com/r/rust/comments/1ekwori/what_would_be_the_fastest_performance_wise_way_to/)
9. Algorithm that solves sudoku with bit operators \[closed\] \- Stack Overflow, accessed May 4, 2026, [https://stackoverflow.com/questions/26787428/algorithm-that-solves-sudoku-with-bit-operators](https://stackoverflow.com/questions/26787428/algorithm-that-solves-sudoku-with-bit-operators)
10. Algorithm for Hidden Sets : r/sudoku \- Reddit, accessed May 4, 2026, [https://www.reddit.com/r/sudoku/comments/18z95v3/algorithm\_for\_hidden\_sets/](https://www.reddit.com/r/sudoku/comments/18z95v3/algorithm_for_hidden_sets/)
11. Solving Sudoku using Bitwise Algorithm \- GeeksforGeeks, accessed May 4, 2026, [https://www.geeksforgeeks.org/dsa/solving-sudoku-using-bitwise-algorithm/](https://www.geeksforgeeks.org/dsa/solving-sudoku-using-bitwise-algorithm/)
12. Writing a High Performance Sudoku Solver in Rust \- Reddit, accessed May 4, 2026, [https://www.reddit.com/r/rust/comments/1slni1k/writing\_a\_high\_performance\_sudoku\_solver\_in\_rust/](https://www.reddit.com/r/rust/comments/1slni1k/writing_a_high_performance_sudoku_solver_in_rust/)
13. 37\. Sudoku Solver \- In-Depth Explanation \- AlgoMonster, accessed May 4, 2026, [https://algo.monster/liteproblems/37](https://algo.monster/liteproblems/37)
14. Solving Sudoku Puzzles with Go \- 1729.org.uk, accessed May 4, 2026, [https://1729.org.uk/posts/sudoku/](https://1729.org.uk/posts/sudoku/)
15. nilostolte/Sudoku: Simple 9x9 Sudoku brute force solver with intrinsic parallel candidate set processing using bits to represent digits in the \[1, 9\] range, and bitwise operations to test a candidate against the candidate set, all at once. · GitHub, accessed May 4, 2026, [https://github.com/nilostolte/Sudoku](https://github.com/nilostolte/Sudoku)
16. Sudoku techniques \- Conceptis Puzzles, accessed May 4, 2026, [https://www.conceptispuzzles.com/index.aspx?uri=puzzle/sudoku/techniques](https://www.conceptispuzzles.com/index.aspx?uri=puzzle/sudoku/techniques)
17. Sudoku Rules \- Strategies, solving techniques and tricks, accessed May 4, 2026, [https://sudoku.com/sudoku-rules/](https://sudoku.com/sudoku-rules/)
18. Sudoku Assistant \-- Solving Techniques \- St. Olaf College, accessed May 4, 2026, [https://www.stolaf.edu/people/hansonr/sudoku/explain.htm](https://www.stolaf.edu/people/hansonr/sudoku/explain.htm)
19. Sudoku Solving Techniques | PDF \- Scribd, accessed May 4, 2026, [https://www.scribd.com/document/666602496/Sudoku-Solving-Techniques](https://www.scribd.com/document/666602496/Sudoku-Solving-Techniques)
20. How to Find and Solve Jellyfish \- Sudoku Handmade Classics \#5 \- YouTube, accessed May 4, 2026, [https://www.youtube.com/watch?v=AsTVGxpJZEg](https://www.youtube.com/watch?v=AsTVGxpJZEg)
21. Exact Method for Generating Strategy-Solvable Sudoku Clues \- MDPI, accessed May 4, 2026, [https://www.mdpi.com/1999-4893/13/7/171](https://www.mdpi.com/1999-4893/13/7/171)
22. Terminology \- Sudopedia, accessed May 4, 2026, [https://www.sudopedia.org/wiki/Terminology](https://www.sudopedia.org/wiki/Terminology)
23. Strategies and Algorithms of Sudoku \- Louisiana Tech Digital Commons, accessed May 4, 2026, [https://digitalcommons.latech.edu/cgi/viewcontent.cgi?article=1012\&context=mathematics-senior-capstone-papers](https://digitalcommons.latech.edu/cgi/viewcontent.cgi?article=1012&context=mathematics-senior-capstone-papers)
24. Using solving techniques for generating puzzles : Help with puzzles, accessed May 4, 2026, [http://forum.enjoysudoku.com/using-solving-techniques-for-generating-puzzles-t30670.html](http://forum.enjoysudoku.com/using-solving-techniques-for-generating-puzzles-t30670.html)
25. What if Naked pair/triplet found within chute \- The New Sudoku Players' Forum, accessed May 4, 2026, [http://forum.enjoysudoku.com/what-if-naked-pair-triplet-found-within-chute-t32696.html](http://forum.enjoysudoku.com/what-if-naked-pair-triplet-found-within-chute-t32696.html)
26. Solving the Sudoku Minimum Number of Clues Problem via Hitting Set Enumeration \- Mathematical Notes, accessed May 4, 2026, [https://www.math.ie/McGuire\_V2.pdf](https://www.math.ie/McGuire_V2.pdf)
27. LOGGING TIME MATH BITS SUDOKU ANSWERS PLOVERORE \- Carnaval de Rua, accessed May 4, 2026, [https://carnavalderua.prefeitura.sp.gov.br/book-search/BpMM8V/7FE126/LoggingTimeMathBitsSudokuAnswersPloverore.pdf](https://carnavalderua.prefeitura.sp.gov.br/book-search/BpMM8V/7FE126/LoggingTimeMathBitsSudokuAnswersPloverore.pdf)
28. logging time math bits sudoku answers ploverore \- Carnaval de Rua, accessed May 4, 2026, [https://carnavalderua.prefeitura.sp.gov.br/Resources/BpMM8V/7FE126/logging\_time\_math\_bits-sudoku-answers\_ploverore.pdf](https://carnavalderua.prefeitura.sp.gov.br/Resources/BpMM8V/7FE126/logging_time_math_bits-sudoku-answers_ploverore.pdf)
29. a new (?) view of fish (naked or hidden) : Advanced solving techniques, accessed May 4, 2026, [http://forum.enjoysudoku.com/a-new-view-of-fish-naked-or-hidden-t5017.html](http://forum.enjoysudoku.com/a-new-view-of-fish-naked-or-hidden-t5017.html)
30. A study of Sudoku solving algorithms, accessed May 4, 2026, [https://www.csc.kth.se/utbildning/kth/kurser/DD143X/dkand12/Group6Alexander/final/Patrik\_Berggren\_David\_Nilsson.report.pdf](https://www.csc.kth.se/utbildning/kth/kurser/DD143X/dkand12/Group6Alexander/final/Patrik_Berggren_David_Nilsson.report.pdf)
31. Naked Candidates \- SudokuWiki.org, accessed May 4, 2026, [https://www.sudokuwiki.org/naked\_candidates](https://www.sudokuwiki.org/naked_candidates)
32. Solve Sudoku using Naked Pairs Triples or Quads, accessed May 4, 2026, [https://sudokubliss.com/guides/naked-pairs-triples-quads](https://sudokubliss.com/guides/naked-pairs-triples-quads)
33. Can someone explain this to Naked triples and naked quads to me like I'm 5? : r/sudoku, accessed May 4, 2026, [https://www.reddit.com/r/sudoku/comments/1qqr2ax/can\_someone\_explain\_this\_to\_naked\_triples\_and/](https://www.reddit.com/r/sudoku/comments/1qqr2ax/can_someone_explain_this_to_naked_triples_and/)
34. How to Identify Hidden Pairs in Sudoku: Unlocking Advanced Strategies, accessed May 4, 2026, [https://www.247sudoku.com/news/how-to-identify-hidden-pairs-in-sudoku-unlocking-advanced-strategies/](https://www.247sudoku.com/news/how-to-identify-hidden-pairs-in-sudoku-unlocking-advanced-strategies/)
35. Let the education continue\! Possible Quad in Center Block? : r/sudoku \- Reddit, accessed May 4, 2026, [https://www.reddit.com/r/sudoku/comments/xtpwvy/let\_the\_education\_continue\_possible\_quad\_in/](https://www.reddit.com/r/sudoku/comments/xtpwvy/let_the_education_continue_possible_quad_in/)
36. The Hidden Logic of Sudoku \- ResearchGate, accessed May 4, 2026, [https://www.researchgate.net/profile/Denis-Berthier/publication/280301600\_The\_Hidden\_Logic\_of\_Sudoku/links/5ed33f62299bf1c67d2cbb74/The-Hidden-Logic-of-Sudoku.pdf](https://www.researchgate.net/profile/Denis-Berthier/publication/280301600_The_Hidden_Logic_of_Sudoku/links/5ed33f62299bf1c67d2cbb74/The-Hidden-Logic-of-Sudoku.pdf)
37. Category:Solving Techniques \- Sudopedia, accessed May 4, 2026, [https://www.sudopedia.org/wiki/Category:Solving\_Techniques](https://www.sudopedia.org/wiki/Category:Solving_Techniques)
38. Sunnie-Shine's C\# Sudoku solver/rater : Software, accessed May 4, 2026, [http://forum.enjoysudoku.com/sunnie-shine-s-c-sudoku-solver-rater-t37788.html](http://forum.enjoysudoku.com/sunnie-shine-s-c-sudoku-solver-rater-t37788.html)
39. Jellyfish: The 4×4 Fish Pattern \- Minimal Sudoku, accessed May 4, 2026, [https://www.minimal-sudoku.com/learn/jellyfish](https://www.minimal-sudoku.com/learn/jellyfish)
40. HoDoKu: User Manual (Chapter 7: Reference) \- SourceForge, accessed May 4, 2026, [https://hodoku.sourceforge.net/en/docs\_ref.php](https://hodoku.sourceforge.net/en/docs_ref.php)
41. Swordfish Strategy \- SudokuWiki.org, accessed May 4, 2026, [https://www.sudokuwiki.org/sword\_fish\_strategy](https://www.sudokuwiki.org/sword_fish_strategy)
42. The X-wing, and its sister solving methods (swordfish, jellyfish, and squirmbag), are found in some of the hardest sudoku puzzles., accessed May 4, 2026, [https://www.math.uci.edu/\~brusso/Scan%202(1).pdf](https://www.math.uci.edu/~brusso/Scan%202\(1\).pdf)
43. ecprice/wordlist \- MIT, accessed May 4, 2026, [https://www.mit.edu/\~ecprice/wordlist.100000](https://www.mit.edu/~ecprice/wordlist.100000)
44. 6 Advanced Sudoku Strategies explained \- SudokuOnline.io, accessed May 4, 2026, [https://www.sudokuonline.io/tips/advanced-sudoku-strategies](https://www.sudokuonline.io/tips/advanced-sudoku-strategies)
45. Skyscraper, How to solve sudoku puzzles, accessed May 4, 2026, [https://www.sudoku9981.com/sudoku-solving/skyscraper.php](https://www.sudoku9981.com/sudoku-solving/skyscraper.php)
46. 123 Parallel Computational Technologies \- ResearchGate, accessed May 4, 2026, [https://www.researchgate.net/profile/Mikhail\_Zymbler/publication/321501607\_Parallel\_Computational\_Technologies\_11th\_International\_Conference\_PCT\_2017\_Kazan\_Russia\_April\_3-7\_2017\_Revised\_Selected\_Papers/links/5e58ea4d299bf1bdb8411d6c/Parallel-Computational-Technologies-11th-International-Conference-PCT-2017-Kazan-Russia-April-3-7-2017-Revised-Selected-Papers.pdf](https://www.researchgate.net/profile/Mikhail_Zymbler/publication/321501607_Parallel_Computational_Technologies_11th_International_Conference_PCT_2017_Kazan_Russia_April_3-7_2017_Revised_Selected_Papers/links/5e58ea4d299bf1bdb8411d6c/Parallel-Computational-Technologies-11th-International-Conference-PCT-2017-Kazan-Russia-April-3-7-2017-Revised-Selected-Papers.pdf)
47. A Graph-Theoretic Solution to NP-Complete Sudoku Puzzles: Malatya Centrality for Large-Scale Optimization \- IEEE Xplore, accessed May 4, 2026, [https://ieeexplore.ieee.org/iel8/6287639/10820123/11205359.pdf](https://ieeexplore.ieee.org/iel8/6287639/10820123/11205359.pdf)
48. Best Sudoku Multi-Coloring Video Ever\! Sudoku Advanced Tutorial 34 \- YouTube, accessed May 4, 2026, [https://www.youtube.com/watch?v=beRuCYE7j74](https://www.youtube.com/watch?v=beRuCYE7j74)
49. GitHub \- alexdej/hodoku-py: A complete Python port of HoDoKu's, accessed May 4, 2026, [https://github.com/alexdej/hodoku-py](https://github.com/alexdej/hodoku-py)
50. Defining the Spectrum of Difficulty? : Help with puzzles and solving techniques \- The New Sudoku Players' Forum, accessed May 4, 2026, [http://forum.enjoysudoku.com/defining-the-spectrum-of-difficulty-t30293.html](http://forum.enjoysudoku.com/defining-the-spectrum-of-difficulty-t30293.html)
51. UMAP 29.3, accessed May 4, 2026, [https://people.math.wisc.edu/\~spagnolie/COMAP/solns/2008mcmsolutions.pdf](https://people.math.wisc.edu/~spagnolie/COMAP/solns/2008mcmsolutions.pdf)
52. Solving the Sudoku Minimum Number of Clues Problem via Hitting Set Enumeration \- arXiv, accessed May 4, 2026, [https://arxiv.org/pdf/1201.0749](https://arxiv.org/pdf/1201.0749)

