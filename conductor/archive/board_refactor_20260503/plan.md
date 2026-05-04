# Implementation Plan: Refactor Board Separation of Concerns

## Phase 1: Core Board State Package [checkpoint: 90acff9]
- [x] Task: Create `internal/board` package b084830
    - [ ] Create `internal/board/board.go` for the core grid structure, getters, and setters.
    - [ ] Create `internal/board/board_test.go` with unit tests for state access and validation.
- [x] Task: Conductor - User Manual Verification 'Phase 1: Core Board State Package' (Protocol in workflow.md) 90acff9

## Phase 2: IO Package (Parsing and Formatting) [checkpoint: e417053]
- [x] Task: Create `internal/io` package 12c6adc
    - [ ] Create `internal/io/parse.go` extracting parsing logic from the old board.
    - [ ] Create `internal/io/parse_test.go` with unit tests for parsing strings.
    - [ ] Create `internal/io/format.go` with formatting logic (to string) from the old board.
    - [ ] Create `internal/io/format_test.go` with unit tests for formatting.
- [x] Task: Conductor - User Manual Verification 'Phase 2: IO Package (Parsing and Formatting)' (Protocol in workflow.md) e417053

## Phase 3: Update Solver and CLI [checkpoint: 6781929]
- [x] Task: Update `internal/solver` to use new packages 986359c
    - [ ] Refactor `internal/solver/solver.go` to depend on `internal/board`.
    - [ ] Update `internal/solver/solver_test.go`.
    - [ ] Remove `internal/solver/board.go` and `internal/solver/board_test.go`.
- [x] Task: Update `internal/cmd` CLI commands 986359c
    - [ ] Update `check.go`, `convert.go`, `hint.go`, `solve.go` and their tests to use `internal/board` and `internal/io`.
- [x] Task: Conductor - User Manual Verification 'Phase 3: Update Solver and CLI' (Protocol in workflow.md) 6781929

## Phase: Review Fixes
- [x] Task: Apply review suggestions 5a05971