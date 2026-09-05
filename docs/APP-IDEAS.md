# Sudoku Solver Ideas

Brainstorming ideas for Sudoku solver application, in no particular order.

- Sudoku solver using human-like deductive logic
- Cross-platform terminal CLI application
- Simple Sudoku Technique Set (SSTS) <http://www.angusj.com/sudoku/hints.php>, in order:
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
- Extensible: iteratively add more techniques later
- Multiple input formats: common ASCII formats used online
  - 81-character string
  - 9x9 ASCII grid (with or without borders and dividers)
  - Full "pencil mark" or "candidates" grid
  - HoDoKu library format for testing? <https://hodoku.sourceforge.net/en/libs.php>
- Same formats for output also
- Print completed puzzle if solved, pencil mark state if not fully solved
- Option to show log of solution steps: HoDoKu format <https://hodoku.sourceforge.net/en/docs_sol.php>
- Option to use a brute force fallback if can't solve with deductive logic (off by default)
- Written in Go language
- Modern Go 1.24+ idioms: range-over-int loops, use `slices` and `maps` packages, `b.Loop()` for benchmark tests
- Fast!
- Go optimization techniques: CPU cache locality, maximize branch prediction, function inlining, zero allocations inside hot search loops (limited allocations per found technique)
- Multiple subcommands:
  - `sudoku solve` (required) solve a puzzle grid
  - `sudoku check` (optional) check if a puzzle grid is valid (one unique solution)
  - `sudoku print` (optional) output a puzzle grid in a different format
  - `sudoku rate` (future) rate the difficulty of a puzzle grid
  - `sudoku generate` (future) generate a valid puzzle grid for a given difficulty
- Use `spf13/cobra` for CLI: de-facto standard, modular, provides automatic CLI help
- Use color in terminal output: `fatih/color` for cross-platform support?
- Directory structure conventions: commands in `cmd`, support packages in `internal`
- Can test each technique independently
- Can benchmark test each technique independently for performance testing and algorithm optimization
