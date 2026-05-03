# Sudoku2

**Solve puzzles faster and learn advanced techniques directly from your terminal.**

Sudoku2 is a powerful command-line tool designed for Sudoku enthusiasts. Whether you're stuck on a difficult puzzle and need a hint, or you're a developer looking for a fast solver to integrate into your workflow, Sudoku2 provides the logic and speed you need.

## Why use Sudoku2?

- **Human-Like Solving:** Uses deductive logic to solve puzzles just like a person would, rather than just brute-forcing the answer.
- **Learn as You Go:** Get step-by-step hints that explain the logic behind the next move, helping you master advanced techniques.
- **Instant Results:** Optimized for performance, solving even the toughest puzzles in milliseconds.
- **Verify Your Puzzles:** Check if your custom-made Sudoku has a unique solution.
- **Flexible Formats:** Easily convert puzzles between different text formats for use in other apps.

## Installation

### Building from Source

If you have Go installed (v1.26.2 or later), you can build Sudoku2 directly:

```bash
# Clone the repository
git clone https://github.com/kpitt/sudoku2.git
cd sudoku2

# Build the tool
make build
```

The `sudoku` executable will be located in the `bin/` directory.

## Getting Started

Sudoku2 uses simple commands to interact with puzzles. Puzzles are represented as 81-character strings, where `.` or `0` represents an empty cell.

### 1. Solve a Puzzle
Get the full solution for a puzzle instantly:

```bash
./bin/sudoku solve "53..7....6..195....98....6.8...6...34..8.3..17...2...6.6....28....419..5....8..79"
```

### 2. Get a Logic Hint
Stuck? Ask Sudoku2 for a hint. It will explain the next logical step without spoiling the entire puzzle:

```bash
./bin/sudoku hint "53..7....6..195....98....6.8...6...34..8.3..17...2...6.6....28....419..5....8..79"
```

### 3. Check for Uniqueness
Ensure a puzzle has exactly one valid solution:

```bash
./bin/sudoku check "53..7....6..195....98....6.8...6...34..8.3..17...2...6.6....28....419..5....8..79"
```

## Command Reference

- `solve`: Outputs the completed board.
- `hint`: Provides a logical explanation for the next move.
- `check`: Validates the puzzle state and uniqueness.
- `convert`: Changes the puzzle representation (e.g., from a string to a grid).

For a full list of options and flags, run:
```bash
./bin/sudoku --help
```
