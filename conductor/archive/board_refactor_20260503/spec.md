# Specification: Refactor Board Separation of Concerns

## Overview
Refactor `internal/solver/board.go` to implement proper separation of concerns. The goal is to separate the storage and management of the board state from parsing, formatting, and the solving logic itself. The new components will reside in new dedicated sub-packages. All existing callers (CLI, tests) will be updated simultaneously to use the new architecture.

## Functional Requirements
- **Board State Management (`internal/board`):** Create a dedicated sub-package for the core Sudoku grid state, including getters, setters, and basic state validation.
- **I/O & Formatting (`internal/io`):** Extract logic for both parsing (instantiating a board from strings/text) and formatting (converting a board state into strings) into a single, unified package.
- **Solver Logic (`internal/solver`):** Ensure the solver only depends on the public interface of the new board state package and is stripped of parsing/formatting concerns.
- **Caller Updates:** Update `cmd/sudoku` and `internal/cmd` to use the new packages for reading input, solving, and printing output.

## Non-Functional Requirements
- **Performance:** The refactoring must not degrade the performance of the solver.
- **Testability:** The separated components (state, io, solving) must be independently testable with their own unit tests.

## Acceptance Criteria
- [ ] The `internal/solver/board.go` file is removed or significantly reduced, containing no parsing or formatting logic.
- [ ] New packages (`internal/board`, `internal/io`) exist with clear, single responsibilities.
- [ ] All unit tests pass.
- [ ] The CLI application builds and functions correctly for solving, hinting, and converting puzzles.

## Out of Scope
- Adding new solving algorithms.
- Adding new CLI features.