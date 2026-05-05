# Sudoku CLI Solver Requirements Document

## 1. Product Overview
The product is a high-performance Command Line Interface (CLI) application designed to solve, validate, and print Sudoku puzzles. The solver utilizes human-style deductive solving techniques, providing a step-by-step log of the solution process, with an optional fallback to a brute-force algorithm. 

## 2. Command-Line Interface (CLI)
The application provides three primary subcommands: `solve`, `check`, and `print`. Output to a terminal console can utilize ANSI color codes to visually distinguish specific elements, such as given values versus user-placed values.

### 2.1. Subcommand: `solve`
Solves a given Sudoku puzzle using deductive techniques.
*   **Behavior:** 
    *   Prints the initial puzzle using the `pretty` format.
    *   Logs each deductive step used to solve the puzzle in a specific format (see Section 6).
    *   Prints the fully solved puzzle using the `pretty` format upon success.
    *   If deductive techniques are insufficient, finishes solving the puzzle using a brute-force algorithm (logged as a single "Brute Force" step).
    *   If the puzzle cannot be solved at all, prints a "failed to solve" error and displays the partially completed puzzle in the `prettypm` format.
*   **Options:**
    *   `--hint`: Modifies behavior to only process and display the first deductive step in the solution as a hint, rather than solving the entire puzzle.

### 2.2. Subcommand: `check`
Validates whether a given Sudoku puzzle has a single, unique solution.
*   **Behavior:** Evaluates the puzzle and reports if it is valid. If invalid, explicitly reports whether the puzzle has *multiple solutions* or *no solutions*.

### 2.3. Subcommand: `print`
Outputs a given Sudoku puzzle in a specified visual or text-based format.
*   **Options:**
    *   `--format <fmt>`: Specifies the desired output format (see Section 4 for supported formats).

## 3. Input Specifications
The application accepts an initial Sudoku puzzle state through mutually exclusive input methods and parses it flexibly.

### 3.1. Input Methods
Input can be provided in one of three ways:
1.  **Positional Argument:** An 81-character string provided directly after the subcommand.
2.  **File Input:** Read from a file on disk, specified via the `--file <filename>` CLI option (mutually exclusive with the positional argument).
3.  **Standard Input (stdin):** Read from stdin if neither a positional argument nor the `--file` option is provided.

### 3.2. Parsing Rules
The input parser is highly resilient and supports a wide variety of puzzle representations:
*   **Empty Cells:** Characters such as `0`, `.`, `x`, `X`, `*`, or `-` represent empty cells. The first non-digit character (excluding brackets and whitespace) encountered is treated as the empty cell marker for the remainder of the puzzle.
*   **Whitespace & Comments:** Empty lines, leading/trailing whitespace, and comments are ignored. Comments begin with a `#` at the start of a line or preceded by at least one space or tab character.
*   **Borders & Separators:** ASCII/Unicode borders and separators (e.g., `|`) are ignored.
*   **Process:** The parser simplifies the grid down to a basic 81-character string using the following exact steps:
    1. Remove any end-of-line comments, then trim all leading and trailing whitespace from each line.
    2. Remove all empty lines and lines that contain only a horizontal border or separator. The result should be either a single line or exactly 9 lines. If the number of remaining lines is not 1 or 9, then it is an invalid format.
    3. If there are 9 lines, remove the newlines and join them together into a single string.
    4. Remove all whitespace characters and vertical pipe separators (`|`). The result should be a single string with exactly 81 characters. If the number of characters is not exactly 81, then it is an invalid format.
    5. Process the 81-character string as a sequence of cells in row-major order. The first character in the string that is not a digit from 1 to 9 will be used as the empty cell marker for the remainder of the string.

### 3.3. Supported Input Formats
The parser accepts any format described in Section 4 (excluding `prettypm`), including variations with or without borders/separators. It natively supports parsing standard puzzle grids (givens/placed values) as well as candidate/pencil-mark grids (e.g., "Simple Sudoku" candidate grids and the "SadMan Sudoku" format).
*   **Candidate Grid Detection:** If the parsed input has exactly 9 lines, 9 groups of digits per line separated by whitespace or `|`, no non-digit empty markers, and at least one group with multiple digits, it is processed as a candidate grid. All single-digit groups in a candidate grid are assumed to be given values.
*   **SadMan Fallback:** If parsing the "SadMan" format and only the `[Puzzle]` section is provided, the `[Puzzle]` header may optionally be omitted. In this case, the format is identical to the standard 9x9-digit grid format.

## 4. Output Formats
The `--format <fmt>` option for the `print` command supports the following outputs:

### 4.1. `raw`
A single 81-character string in row-major order. Uses `0` for empty cells in output.
```
310004069000000200008005040000000005006000017807030000590700006600003050000100002
```

### 4.2. `9x9`
A simple borderless 9x9 grid, one row per line. Uses `.` for empty cells in output.
```
31...4.69
......2..
..8..5.4.
........5
..6....17
8.7.3....
59.7....6
6....3.5.
...1....2
```

### 4.3. `ss` ("Simple Sudoku")
The standard grid output format for the "Simple Sudoku" program, with ASCII borders.
```
*-----------*
|.3.|4..|...|
|9.2|8.6|3.1|
|...|...|.2.|
|---+---+---|
|8..|.6.|7..|
|.6.|2.5|.9.|
|..3|.4.|..8|
|---+---+---|
|.7.|...|...|
|4.8|9.2|5.6|
|...|..8|.3.|
*-----------*
```

### 4.4. `ssnb`
An older variant of the "Simple Sudoku" format omitting the outer borders.
```
1..|...|7..
.2.|...|5..
6..|38.|...
-----------
.78|...|...
...|6.9|...
...|...|14.
-----------
...|.25|..9
..3|...|.6.
..4|...|..2
```

### 4.5. `sswide`
A variation of the "Simple Sudoku" format adding spaces between digits for a squarer appearance.
```
*-----------------------*
| . 3 . | 4 . . | . . . |
| 9 . 2 | 8 . 6 | 3 . 1 |
| . . . | . . . | . 2 . |
|-------+-------+-------|
| 8 . . | . 6 . | 7 . . |
| . 6 . | 2 . 5 | . 9 . |
| . . 3 | . 4 . | . . 8 |
|-------+-------+-------|
| . 7 . | . . . | . . . |
| 4 . 8 | 9 . 2 | 5 . 6 |
| . . . | . . 8 | . 3 . |
*-----------------------*
```

### 4.6. `pm` ("Simple Sudoku" Candidate Grid)
Shows remaining candidate digits (pencil marks). Every column dynamically adjusts its width to match the cell with the most candidates in that column.
```
*--------------------------------------------------------------------*
| 1567   3      1567   | 4      12579  179    | 689    5678   579    |
| 9      45     2      | 8      57     6      | 3      457    1      |
| 1567   1458   14567  | 1357   13579  1379   | 4689   2      4579   |
|----------------------+----------------------+----------------------|
| 8      12459  1459   | 13     6      139    | 7      145    2345   |
| 17     6      147    | 2      1378   5      | 14     9      34     |
| 1257   1259   3      | 17     4      179    | 126    156    8      |
|----------------------+----------------------+----------------------|
| 12356  7      1569   | 1356   135    134    | 12489  148    249    |
| 4      1      8      | 9      137    2      | 5      17     6      |
| 1256   1259   1569   | 1567   157    8      | 1249   3      2479   |
*--------------------------------------------------------------------*
```

### 4.7. `sadman`
"SadMan Sudoku" format containing up to 3 headers: `[Puzzle]` (initial state), `[State]` (current state, omitted if identical to Puzzle), and `[PencilMarks]`.
```
[Puzzle]
6.......7
....9..2.
3.1..259.
8....7.13
....8....
76.3....8
.782..1.6
.5..3....
2.......9
[State]
6.....387
547893621
381..2594
8....7.13
....8.765
76.3...48
.782..136
.5..3..72
2......59
[PencilMarks]
,29,29,145,145,145,,,
,,,,,,,,
,,,67,67,,,,
,29,2459,4569,2456,,29,,
149,1239,2349,149,,149,29,,
,,259,,125,159,29,,
49,,,,45,459,,,
149,,469,1469,,14689,48,,
,13,346,1467,1467,1468,48,,
```

### 4.8. `pretty`
A clean grid utilizing Unicode line-drawing characters. *(Default for `solve` command's initial and final state).*
```
┏━━━━━━━┯━━━━━━━┯━━━━━━━┓
┃ . 3 . │ 4 . . │ . . . ┃
┃ 9 . 2 │ 8 . 6 │ 3 . 1 ┃
┃ . . . │ . . . │ . 2 . ┃
┠───────┼───────┼───────┨
┃ 8 . . │ . 6 . │ 7 . . ┃
┃ . 6 . │ 2 . 5 │ . 9 . ┃
┃ . . 3 │ . 4 . │ . . 8 ┃
┠───────┼───────┼───────┨
┃ . 7 . │ . . . │ . . . ┃
┃ 4 . 8 │ 9 . 2 │ 5 . 6 ┃
┃ . . . │ . . 8 │ . 3 . ┃
┗━━━━━━━┷━━━━━━━┷━━━━━━━┛
```

For direct terminal output, this format can be enhanced with ANSI color codes, e.g. dark gray for empty cell markers, and different colors to distinguish given values from placed values.

### 4.9. `prettypm`
A visually rich candidate grid utilizing Unicode characters, displaying placed values centered with `[]` brackets around givens. *(Default for `solve` command's failure state).*
```
┏━━━━━┯━━━━━┯━━━━━┳━━━━━┯━━━━━┯━━━━━┳━━━━━┯━━━━━┯━━━━━┓
┃     ┆1    ┆1    ┃     ┆     ┆     ┃     ┆     ┆     ┃
┃  6  ┆4    ┆     ┃ [7] ┆ [2] ┆ [3] ┃  5  ┆4    ┆4 5  ┃
┃     ┆    9┆    9┃     ┆     ┆     ┃  8 9┆  8 9┆     ┃
┠╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌┨
┃  2  ┆  2  ┆     ┃1    ┆     ┆1    ┃     ┆     ┆     ┃
┃4 5  ┆4    ┆ [3] ┃4    ┆4 5  ┆  5  ┃ [6] ┆4    ┆  7  ┃
┃  8  ┆    9┆     ┃    9┆    9┆  8  ┃     ┆    9┆     ┃
┠╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌┨
┃     ┆     ┆     ┃     ┆     ┆     ┃     ┆     ┆     ┃
┃4 5  ┆4    ┆  5  ┃4    ┆  6  ┆  5  ┃  3  ┆ [2] ┆ [1] ┃
┃  8  ┆7   9┆7   9┃    9┆     ┆  8  ┃     ┆     ┆     ┃
┣━━━━━┿━━━━━┿━━━━━╋━━━━━┿━━━━━┿━━━━━╋━━━━━┿━━━━━┿━━━━━┫
┃     ┆     ┆     ┃     ┆     ┆     ┃     ┆     ┆     ┃
┃ [9] ┆  5  ┆  2  ┃  6  ┆  7  ┆ [4] ┃  1  ┆ [3] ┆  8  ┃
┃     ┆     ┆     ┃     ┆     ┆     ┃     ┆     ┆     ┃
┠╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌┨
┃     ┆     ┆     ┃     ┆     ┆     ┃     ┆     ┆     ┃
┃ [1] ┆  3  ┆  6  ┃  5  ┆ [8] ┆  2  ┃  4  ┆  7  ┆ [9] ┃
┃     ┆     ┆     ┃     ┆     ┆     ┃     ┆     ┆     ┃
┠╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌┨
┃     ┆     ┆     ┃     ┆     ┆     ┃     ┆     ┆     ┃
┃  7  ┆ [8] ┆ [4] ┃  3  ┆ [1] ┆  9  ┃  2  ┆  5  ┆  6  ┃
┃     ┆     ┆     ┃     ┆     ┆     ┃     ┆     ┆     ┃
┣━━━━━┿━━━━━┿━━━━━╋━━━━━┿━━━━━┿━━━━━╋━━━━━┿━━━━━┿━━━━━┫
┃     ┆     ┆1    ┃     ┆     ┆     ┃     ┆1    ┆     ┃
┃ [3] ┆ [6] ┆  5  ┃  8  ┆4 5  ┆  7  ┃  5  ┆4    ┆  2  ┃
┃     ┆     ┆    9┃     ┆    9┆     ┃    9┆    9┆     ┃
┠╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌┨
┃  2  ┆  2  ┆     ┃1    ┆     ┆1    ┃     ┆     ┆     ┃
┃4 5  ┆4    ┆ [8] ┃4    ┆4 5  ┆  5  ┃ [7] ┆  6  ┆  3  ┃
┃     ┆    9┆     ┃    9┆    9┆     ┃     ┆     ┆     ┃
┠╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌╂╌╌╌╌╌┼╌╌╌╌╌┼╌╌╌╌╌┨
┃     ┆1    ┆1    ┃     ┆     ┆     ┃     ┆1    ┆     ┃
┃4 5  ┆     ┆     ┃ [2] ┆ [3] ┆ [6] ┃     ┆     ┆4 5  ┃
┃     ┆7   9┆7   9┃     ┆     ┆     ┃  8 9┆  8 9┆     ┃
┗━━━━━┷━━━━━┷━━━━━┻━━━━━┷━━━━━┷━━━━━┻━━━━━┷━━━━━┷━━━━━┛
```

For direct terminal output, this format can be enhanced with ANSI color codes, e.g. dark gray for candidate values, and different colors to distinguish given values from placed values.

## 5. Solving Capabilities
The solver implements human-style deductive techniques, progressing in complexity.

### 5.1. Core Techniques (SSTS)
The solver must implement the Simple Sudoku Technique Set (SSTS):
1. Naked Single
2. Hidden Single
3. Locked Candidates (Type 1 Pointing, Type 2 Claiming)
4. Naked Pair / Triple / Quad
5. Hidden Pair / Triple / Quad
6. X-Wing / Swordfish / Jellyfish
7. XY-Wing
8. Simple Colors / Multi-Colors

### 5.2. Advanced Techniques (Optional)
The solver may implement additional advanced techniques to improve deductive success rates (e.g., XYZ-Wing, W-Wing, Remote Pair, Skyscraper, 2-String Kite, Empty Rectangle, Unique Rectangle, Finned Fish).

### 5.3. Brute Force Fallback
A constraint-propagating backtracking algorithm (e.g., Dancing Links) can finish solving the puzzle if it cannot be fully solved using the implemented deductive techniques. This is logged as a single "Brute Force" step.

## 6. Solution Step Format
The `solve` command outputs a log of steps. Each step defines the technique used and the result applied to the grid, formatted generally as `Technique Name: Result`.

### 6.1. Step Results
*   **Placement:** Assigns a single value to a cell. Format: `rncn=N` (e.g., `Hidden Single: r1c3=5`).
*   **Elimination:** Removes candidate values from cells. For eliminations, the result is formatted as an implication: `Conditions => Eliminations`. 
    *   **Elimination format:** `rncn<>N`. Multiple digits eliminated from different cells are separated by commas followed by a single space (e.g., `r8c7<>3, r8c7<>8, r8c8<>9`), sorted in increasing order of the eliminated value.
    *   **Strict Whitespace Rules:** There are **no spaces** around commas separating disjoint cell references (e.g., `r1c36,r2c1`), but there **is a single space** after a comma separating different digits (as shown above).

**Examples:**
```
1. Naked Triple: 3,9,6 in r245c2 => r1c2<>6
2. X-Wing: 1 c15 r25 => r2c4789,r5c34789<>1
3. Locked Candidates Type 1 (Pointing): 5 in b1 => r3c7<>5
4. XY-Wing: 5/7/2 in r1c36,r2c1 => r2c6<>2
5. Simple Colors: 3 (r1c4,r3c3,r5c6,r6c1,r7c5,r9c7) / (r6c4,r7c3,r8c9,r9c6) => r13c9,r8c1<>3
```

### 6.2. Cell Reference Notation
Cell references utilize a compact `rncn` syntax:
*   `r`: Row index (1-9).
*   `c`: Column index (1-9).
*   `b`: Box index (1-9).
*   **Grouping:** Digits can be grouped to reference multiple cells (e.g., `r456c4` represents r4c4, r5c4, r6c4). Disjoint groups are separated by commas without spaces (e.g., `r1c36,r2c1`). Entire houses can be referenced by omitting the intersecting coordinate (e.g., `r3` for row 3, `b3` for box 3).
