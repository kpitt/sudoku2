# Sudoku 2

## Project Overview

`sudoku2` is a high-performance deductive Sudoku puzzle solver written in Go.

- **Primary Language:** Go
- **Module Path:** `github.com/kpitt/sudoku2`
- **Go Version:** 1.26+

## Core Principles

- Follow standard Go project layout (`/cmd`, `/pkg`, `/internal`).
- Write idiomatic Go: simple, clear, and efficient.
- Write high-performance code: choose efficient algorithms and data structures, minimize heap allocations.
- Use `spf13/cobra` for command structure.

## Building and Running

- **Build:** `go build ./...`
- **Run:** `go run ./cmd/sudoku`
- **Test:** `go test ./...`
- **Lint:** `go vet ./...`

## Technical Rules

- **Dependency Management:** Use `go.mod`. Run `go mod tidy` after modifying dependencies.
- **Error Handling:** Explicitly check errors. Do not ignore them.
- **Naming:** Follow standard Go naming conventions: CamelCase for exported symbols, mixedCaps for packages and private symbols. Use interfaces sparingly.
- **Testing:** Write unit tests in `_test.go` files, using `t.Run` and the standard `testing` package. Use the table-driven pattern for unit tests, with named subtests. Always use a TDD test-first approach: write a failing test first, then update the code until it passes.
- **Git Checkpoints:** Create regular checkpoints using `git` after key milestones: updating design/plans, adding failing tests, implementing atomic functionality, or completing features.
- **Commit Messages:** Use consistent, descriptive formatting for commit messages. Do NOT use "semantic commit" style prefixes (e.g., `feat:`, `fix:`).
- **Performance:** All performance improvements MUST be verified with benchmarks.
- **Documentation:** Every exported symbol must have a doc comment starting with the symbol name.
- **Formatting:** All code must be formatted with `go fmt`.

## CLI Behavior

- Always provide flags for configuration.
- Implement proper `--help` documentation.
- Use `os.Stderr` for errors and `os.Stdout` for data output.
- Print human-readable errors, but keep output clean.

## Constraints

- **NO** panic. Use error returning.
- **NO** global variables for state.
- Strictly limited to standard 9x9 puzzles, no extended variants.
