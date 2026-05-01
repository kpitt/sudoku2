# Project Context

**NOTE:** The project structure below must be maintained and kept in sync with any changes made to the application code.

## Project Structure
- `bin/`            # Compiled binaries (ignored by Git)
- `cmd/`            # Entry points for the application (e.g., `cmd/sudoku/main.go`)
- `conductor/`      # Project orchestration and documentation (Source of Truth)
- `internal/`       # Private application code (Go-specific pattern)
    - `cmd/`        # Implementation logic for CLI commands
    - `solver/`     # Core Sudoku solving algorithms and board logic
- `Makefile`        # Standardized development commands

## Definition
- [Product Definition](./product.md)
- [Product Guidelines](./product-guidelines.md)
- [Tech Stack](./tech-stack.md)

## Workflow
- [Workflow](./workflow.md)
- [Code Style Guides](./code_styleguides/)

## Management
- [Tracks Registry](./tracks.md)
- [Tracks Directory](./tracks/)
