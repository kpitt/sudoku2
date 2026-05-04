# System Architecture Overview

This document outlines the high-level architecture of the Sudoku CLI application. It maps the data flow from the command-line entry point through the parsing, solving, and formatting packages. It provides a structural blueprint based on the product requirements and the underlying technical design of the deductive solver.

## Application Scope

The application is a command-line tool built in Go. It accepts Sudoku puzzles in various text formats, solves them using a progression of deterministic logic techniques, and prints the solution steps and the final grid. If the deterministic logic fails, it uses a brute-force algorithm to finish the puzzle. It also provides utilities to check if a puzzle has a unique solution and to convert puzzles between different text formats.

## Package Structure

The codebase follows a standard Go project layout, dividing responsibilities into distinct packages to enforce separation of concerns.

### 1. `cmd/sudoku` (CLI Layer)

This package handles user interaction. It uses the `github.com/spf13/cobra` library to define the application's commands and flags.

*   **`root`:** The base command that sets up global configuration.
*   **`solve`:** The primary command. It reads the input puzzle, initializes the solver, executes the deductive loop, and coordinates the formatting of the initial board, the step log, and the final board. It also handles the `--hint` flag logic to short-circuit the solver.
*   **`check`:** Evaluates puzzle validity. It skips the deductive engine and invokes the exact cover (DLX) engine directly to count the total number of solutions (0, 1, or multiple).
*   **`print`:** Parses the input puzzle and immediately passes the board state to the formatter, ignoring the solving engine.

### 2. `internal/parser` (Input Ingestion)

The parser package isolates the logic required to read arbitrary text and convert it into a structured format.

*   **Standardization:** It reads data from positional arguments, files, or `stdin`. It strips whitespace, removes ASCII borders, and ignores comments starting with `#`.
*   **Heuristic Detection:** It examines the structure of the input to determine if it is a standard 81-character grid or a complex format like a candidate pencil-mark grid or "SadMan" format.
*   **Tokenization & Initialization:** It extracts the givens and explicit pencil marks, determines the character used to represent empty cells, and outputs an initialized `solver.Board` struct.

### 3. `internal/solver` (Core Computational Engine)

This package contains the high-performance logic required to manipulate and resolve the puzzle state. It operates entirely on bitwise representations and stack-allocated arrays to avoid garbage collection.

*   **`Board`:** The central struct. It contains arrays of `uint16` bitmasks representing the candidates for all 81 cells, along with the aggregated constraint masks for rows, columns, and boxes.
*   **Deductive Techniques:** A suite of independent, stateless functions (e.g., `FindNakedSingles`, `FindXWing`). Each function accepts a `Board` and returns a `Step` object containing any discovered placements or eliminations.
*   **The Solver Loop:** A priority queue that executes the deductive techniques from lowest algorithmic complexity to highest. If a technique succeeds, the loop applies the resulting `Step` to the `Board` and resets to the lowest priority to cascade the logic.
*   **DLX Fallback:** A brute-force exact cover algorithm using array-based Dancing Links. It runs if the deductive loop stalls before the board is fully resolved.

### 4. `internal/formatter` (Output Generation)

The formatter package translates internal memory structures back into human-readable text strings for terminal output.

*   **Grid Rendering:** It translates the `Board` bitmasks into the specific grid layouts requested by the user (`raw`, `9x9`, `ss`, `pm`, `pretty`, `sadman`).
*   **Dynamic Alignment:** For formats like `pm` (pencil marks), it calculates the maximum width of candidates in each column to ensure the printed grid remains aligned.
*   **Step Notation Translation:** It converts the raw integer array indices returned by the solver's `Step` objects into the compact `rncn` coordinate syntax required for the solution log.
*   **Terminal Colorization:** It uses the `github.com/fatih/color` library to apply ANSI color codes to the `pretty` and `prettypm` formats, visually distinguishing original givens from solver-placed digits.

## Data Flow Lifecycle

The typical execution path for the `solve` command illustrates how data moves through the architecture.

1.  **Ingestion:** The user runs `sudoku solve <input>`. The `cmd` package captures the input string and passes it to the `parser`.
2.  **Parsing and Initialization:** The `parser` strips formatting, detects the layout, and directly returns an initialized `solver.Board` struct. This allows the parser to natively apply explicit pencil marks defined in complex input formats.
3.  **Initial Output:** The `cmd` package passes a copy of the initial `Board` to the `formatter` to print the starting grid in the `pretty` format.
4.  **Solving Loop:** The `cmd` package passes the `Board` to the main deductive loop.
    *   The loop tests techniques. When one returns a `Step`, the loop applies it, updates the `Board`, and saves the `Step` to a log array.
    *   If the deductive loop stalls, it triggers the DLX fallback to complete the `Board`, appending a final "Brute Force" step to the log.
5.  **Log Formatting:** The `cmd` package passes the array of `Step` objects to the `formatter`, which translates them into `rncn` syntax and prints them to the terminal.
6.  **Final Output:** The `cmd` package passes the finalized `Board` to the `formatter` to print the completed puzzle.