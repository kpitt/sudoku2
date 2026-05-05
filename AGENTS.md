# Agent Instructions

## Core Principles
- **Source of Truth**: The `docs/REQUIREMENTS.md` file is the foundational source of truth for all high-level product requirements and functional specifications. All implementation must strictly adhere to its definitions.
- Follow standard Go project layout (`/cmd`, `/pkg`, `/internal`).
- Write idiomatic Go: simple, clear, and efficient.
- Use `spf13/cobra` for command structure.
- Prefer high-performance solutions; choose clear readable code when performance difference is negligible.

## Development Workflow (MANDATORY)

### Test-Driven Development (TDD)
1. **Write failing tests FIRST** - No implementation without tests
2. Write the minimal code to make tests pass
3. Refactor while keeping tests green
4. Run `go test ./...` to verify all tests pass
5. Run `golangci-lint run` to ensure code quality

### Code Quality Checklist (Before Commit)
- [ ] All tests pass (`go test ./...`)
- [ ] All benchmarks pass (`go test -bench=. ./...`)
- [ ] golangci-lint passes with no errors (`golangci-lint run`)
- [ ] Code formatted with `gofmt` (or `goimports`)
- [ ] All exported symbols have doc comments
- [ ] `go mod tidy` has been run if dependencies changed

## Performance & Optimization

### Algorithm Selection
- **Prefer high-performance algorithms** suited to the task
- Document non-obvious algorithm choices with rationale
- Consider time/space complexity tradeoffs (Big-O analysis)
- Profile before optimizing (`go test -cpuprofile`, `-memprofile`)

### Memory Efficiency
- **Stack over Heap**: Prefer stack allocations when possible
  - Use value types for small structs (< 100 bytes)
  - Avoid unnecessary pointer indirection
  - Use escape analysis to verify: `go build -gcflags='-m'`
- **Pre-allocate slices** when size is known: `make([]T, 0, capacity)`
- **Reuse buffers** with `sync.Pool` for high-frequency allocations
- **Avoid allocations in hot paths**: minimize `[]byte` → `string` conversions

### CPU Cache Locality
- **Structure layout**: Order struct fields by size (largest first) to minimize padding
- **Data-oriented design**: Use arrays/slices of values, not pointers
- **Sequential access**: Process data in contiguous memory order
- **Batch operations**: Process multiple items in loops to maximize cache hits
- **Small hot structs**: Keep frequently accessed data structures compact

### Branch Predictability
- **Avoid conditionals in tight loops** when possible
- **Most common case first** in if/else chains
- Use table-driven approaches for dispatch logic
- Consider branchless techniques for critical paths

## Modern Go (1.24+) Constructs

### Range-Over-Int (Go 1.22+)
```go
// PREFER: Modern range-over-int
for i := range 10 {
    process(i)
}

// AVOID: Legacy C-style loop
for i := 0; i < 10; i++ {
    process(i)
}
```

### Benchmark Loop (Go 1.24+)
```go
// PREFER: Modern b.Loop()
func BenchmarkFoo(b *testing.B) {
    for b.Loop() {
        foo()
    }
}

// AVOID: Legacy b.N
func BenchmarkFoo(b *testing.B) {
    for i := 0; i < b.N; i++ {
        foo()
    }
}
```

### Other Modern Patterns
- Use `min()`, `max()` built-ins (Go 1.21+)
- Use `clear()` for maps/slices (Go 1.21+)
- Prefer `cmp.Compare` and `cmp.Or` (Go 1.21+)

## Technical Rules

### Dependency Management
- Use `go.mod` for all dependencies
- Run `go mod tidy` after modifying dependencies
- Pin versions for reproducible builds

### Error Handling
- **Explicitly check all errors** - never ignore with `_`
- Wrap errors with context: `fmt.Errorf("operation failed: %w", err)`
- Return errors, don't panic (except for programmer errors in init/main)
- Use sentinel errors sparingly; prefer error types with context

### Naming Conventions
- **Exported symbols**: CamelCase (e.g., `ProcessData`)
- **Unexported symbols**: mixedCaps (e.g., `processData`)
- **Packages**: short, lowercase, no underscores (e.g., `sudoku`, not `sudoku_solver`)
- **Interfaces**: Single-method interfaces end in `-er` (e.g., `Reader`, `Solver`)
- Use meaningful names; avoid cryptic abbreviations

### Documentation Standards
- **Every exported symbol** must have a doc comment
- Start with the symbol name: `// ProcessData validates and processes input data.`
- Follow [godoc conventions](https://go.dev/doc/comment)
- Document panics, special behaviors, and edge cases
- Keep comments concise but complete

### Testing Requirements
- Write tests in `*_test.go` files
- Use table-driven tests with `t.Run()` for multiple cases
- Test edge cases, error paths, and boundary conditions
- Aim for >80% code coverage on critical paths
- Run tests: `go test ./...`

### Benchmark Requirements
- **Write benchmarks for all performance-critical functions**
- Use `*_test.go` files with `Benchmark*` functions
- Use modern `b.Loop()` construct (Go 1.24+)
- Reset timers when needed: `b.ResetTimer()`
- Prevent compiler optimizations: assign to package-level var
- Run benchmarks: `go test -bench=. -benchmem ./...`

Example benchmark:
```go
var result int // prevent compiler optimization

func BenchmarkSolve(b *testing.B) {
    grid := setupTestGrid()
    b.ResetTimer()
    
    for b.Loop() {
        result = solve(grid)
    }
}
```

### Linting & Formatting
- **golangci-lint must pass** before commit: `golangci-lint run`
- Use `gofmt` or `goimports` for formatting
- Fix all linter errors; warnings should be addressed or explicitly ignored with `//nolint` and justification
- Configure `.golangci.yml` in project root for consistency

## CLI Behavior
- Always provide flags for configuration
- Implement proper `--help` documentation
- Use `os.Stderr` for errors and `os.Stdout` for data output
- Print human-readable errors, but keep output clean
- Exit codes: 0 for success, non-zero for errors

## Constraints (MUST NOT)
- **NO** panic (except in init/main for configuration errors)
- **NO** global variables for mutable state
- **NO** ignoring errors with `_` 
- **NO** premature optimization without benchmarks
- **NO** legacy loop constructs when modern alternatives exist (Go 1.22+)

## Performance Profiling Workflow
1. Write benchmark tests
2. Establish baseline: `go test -bench=. -benchmem > old.txt`
3. Make changes
4. Compare: `go test -bench=. -benchmem > new.txt && benchstat old.txt new.txt`
5. Profile if needed: `go test -bench=. -cpuprofile=cpu.prof`
6. Analyze: `go tool pprof cpu.prof`

## References
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Go Performance Tips](https://github.com/golang/go/wiki/Performance)
- [golangci-lint](https://golangci-lint.run/)

