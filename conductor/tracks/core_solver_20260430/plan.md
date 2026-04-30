# Implementation Plan - Implement core Sudoku solving algorithm and CLI structure using Cobra

This plan outlines the steps to build the core Sudoku solver and CLI framework.

## Phase 1: Project Scaffolding and CLI Structure [checkpoint: e22b9d7]
- [x] Task: Initialize Cobra CLI and project structure [8589b29]
    - [x] Install Cobra dependency
    - [x] Initialize Cobra root command in `cmd/sudoku/`
    - [x] Setup basic project layout (internal/pkg/cmd)
- [x] Task: Implement 'solve' command skeleton [94bfdbf]
    - [x] Add `solve` subcommand
    - [x] Define flags for input (e.g., `--input` or positional arg)
- [x] Task: Implement 'check' command skeleton [1d16010]
    - [x] Add `check` subcommand
- [x] Task: Conductor - User Manual Verification 'Project Scaffolding and CLI Structure' (Protocol in workflow.md)

## Phase 2: Core Solver Implementation (TDD) [checkpoint: 0631304]
- [x] Task: Write failing tests for Sudoku board representation and validation [26fadec]
    - [x] Define `Board` struct
    - [x] Write tests for parsing 81-char strings
    - [x] Write tests for basic Sudoku rule validation (rows, cols, boxes)
- [x] Task: Implement Sudoku board representation and basic validation [3769ee7]
    - [x] Implement parsing logic
    - [x] Implement rule validation logic
- [x] Task: Implement backtracking algorithm for 'check' command [1550190]
    - [x] Implement Backtracking algorithm
    - [x] Ensure it can find a solution (if one exists)
- [x] Task: Implement deductive solving techniques for 'solve' command [608ed2c]
    - [x] Implement Naked Singles strategy
    - [x] Implement Hidden Singles strategy
    - [x] Implement iterative deduction loop
    - [x] Ensure all tests pass
- [x] Task: Conductor - User Manual Verification 'Core Solver Implementation' (Protocol in workflow.md)

## Phase 3: CLI Integration and Hints
- [x] Task: Integrate solver into 'solve' command with flag handling [1c2b3a4]
    - [x] Wire the `solve` command to the solver logic
    - [x] Implement pretty-printing for the output grid
- [ ] Task: Implement basic hint system and educational output
    - [ ] Add logic to provide a single move hint
- [ ] Task: Implement format conversion
    - [ ] Add support for different input/output formats
- [ ] Task: Conductor - User Manual Verification 'CLI Integration and Hints' (Protocol in workflow.md)