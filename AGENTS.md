# Agent Instructions

**CRITICAL:** This file (`AGENTS.md`) is the primary context file for this project. DO NOT use agent-specific files such as `GEMINI.md` or `CLAUDE.md`.

## Functional Mappings
- **Application Entry Points:** `cmd/sudoku/main.go`
- **CLI Command Logic:** `internal/cmd/`
- **Sudoku Core Logic:** `internal/solver/`
- **Project Documentation:** `conductor/`

## Conventions
- **Documentation Sync:** The project structure in `conductor/index.md` and the mappings in this file MUST be kept in sync with any changes made to the application code.
- **Loops:** Always use the modern range-over-int form (`for i := range N`) for simple counting iterations instead of traditional C-style loops.
- **Testing:** All tests must be in files suffixed with `_test.go` in the same package as the code being tested.
- **Dependency Management:** Use `go mod tidy` or `make install` for dependencies.
- **Build System:** Use the `Makefile` as the definitive source for development commands (`make build`, `make test`, etc.).
- **Workflow:** Always follow the lifecycle defined in `conductor/workflow.md`.

## Project Index
@conductor/index.md
