# Implementation Plan - Implement core Sudoku solving algorithm and CLI structure using Cobra

This plan outlines the steps to build the core Sudoku solver and CLI framework.

## Phase 1: Project Scaffolding and CLI Structure
- [x] Task: Initialize Cobra CLI and project structure [8589b29]
    - [x] Install Cobra dependency
    - [x] Initialize Cobra root command in `cmd/sudoku/`
    - [x] Setup basic project layout (internal/pkg/cmd)
- [x] Task: Implement 'solve' command skeleton [94bfdbf]
    - [x] Add `solve` subcommand
    - [x] Define flags for input (e.g., `--input` or positional arg)
- [ ] Task: Implement 'validate' command skeleton
    - [ ] Add `validate` subcommand
- [ ] Task: Conductor - User Manual Verification 'Project Scaffolding and CLI Structure' (Protocol in workflow.md)

## Phase 2: Core Solver Implementation (TDD)
- [ ] Task: Write failing tests for Sudoku board representation and validation
    - [ ] Define `Board` struct
    - [ ] Write tests for parsing 81-char strings
    - [ ] Write tests for basic Sudoku rule validation (rows, cols, boxes)
- [ ] Task: Implement Sudoku board representation and basic validation
    - [ ] Implement parsing logic
    - [ ] Implement rule validation logic
- [ ] Task: Write failing tests for core solving algorithm
    - [ ] Create test cases for easy, medium, and hard puzzles
    - [ ] Define expected output for each
- [ ] Task: Implement core solving algorithm
    - [ ] Implement Backtracking or Dancing Links algorithm
    - [ ] Ensure all tests pass
- [ ] Task: Conductor - User Manual Verification 'Core Solver Implementation' (Protocol in workflow.md)

## Phase 3: CLI Integration and Hints
- [ ] Task: Integrate solver into 'solve' command with flag handling
    - [ ] Wire the `solve` command to the solver logic
    - [ ] Implement pretty-printing for the output grid
- [ ] Task: Implement basic hint system and educational output
    - [ ] Add logic to provide a single move hint
- [ ] Task: Implement format conversion
    - [ ] Add support for different input/output formats
- [ ] Task: Conductor - User Manual Verification 'CLI Integration and Hints' (Protocol in workflow.md)