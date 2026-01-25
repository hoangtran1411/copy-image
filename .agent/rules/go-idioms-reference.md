# Go Idioms Reference - Copy Image Project

> **Full Reference Document** - Contains detailed idioms, code examples, and best practices.  
> For compact core rules, see `go-style-guide.md`

---

## Naming Conventions (Idiomatic Go)

- Use short but meaningful names, scoped by context:
  - `r`, `w`, `ctx`, `db`, `tx`, `cfg` are acceptable in small scopes.
  - Avoid `data`, `info`, `obj`, `temp`, `value` unless unavoidable.

- Prefer noun-based names for structs, verb-based names for functions:
  - `Parser.Parse()`, `Generator.Generate()`

- Boolean names should read naturally:
  - `isValid`, `hasHeader`, `enableStreaming`

- Avoid stuttering:
  - ❌ `pkg.PkgThing`
  - ✅ `pkg.Thing`

---

## Function Design

- Prefer small functions (≤ 40 lines).
- One function = one responsibility.
- Avoid flags that change behavior dramatically:

```go
// bad
func Process(input string, strict bool)

// good
func ProcessStrict(input string)
func ProcessLenient(input string)
```

- Return early (guard clauses):

```go
if err != nil {
    return nil, err
}
```

---

## Error Handling Idioms

- Never ignore errors explicitly:

```go
_ = f.Close() // ❌ unless justified in comment
```

- Wrap errors only at package boundaries:

```go
return nil, fmt.Errorf("operation failed: %w", err)
```

- Do not wrap errors multiple times in the same layer.

- Prefer `errors.Is` / `errors.As` for comparisons:

```go
if errors.Is(err, ErrNotFound) { ... }
```

- Avoid sentinel errors unless necessary; prefer typed errors:

```go
type ErrInvalidInput struct {
    Field string
}
```

---

## Package Design & Boundaries

- `internal` packages must be:
  - UI-agnostic
  - Framework-agnostic (no Wails imports)

- Each package should expose minimal API surface:

```go
// good
func Process(ctx context.Context, input string) (*Result, error)

// avoid exposing helpers
```

- Avoid circular dependencies at all cost.
- If two packages depend on each other → redesign.

---

## Context Usage Idioms

- `context.Context` must:
  - Be the first parameter
  - Never be stored in struct fields

- Do not pass nil context:

```go
ctx := context.Background()
```

- Respect cancellation in loops:

```go
select {
case <-ctx.Done():
    return ctx.Err()
default:
}
```

---

## Concurrency Patterns

- Prefer worker pool over unbounded goroutines.
- Always define ownership of goroutines:
  - Who starts?
  - Who stops?

- Use `errgroup.Group` for concurrent tasks with error propagation.

- Channels should have clear direction:

```go
func worker(in <-chan Job, out chan<- Result)
```

- Avoid closing channels you did not create.

---

## Struct & Interface Idioms

- Accept interfaces, return concrete types:

```go
func NewProcessor(r io.Reader) *Processor
```

- Interfaces should be small (1–3 methods):

```go
type Reader interface {
    Read() ([]byte, error)
}
```

- Do not define interfaces prematurely.

---

## Zero Value Philosophy

- Design structs so zero value is usable:

```go
var p Processor // should work
```

- Avoid constructors unless needed for invariants.

- Prefer empty slices over nil slices for JSON output:

```go
items := make([]Item, 0)
```

---

## Slice & Map Best Practices

- Pre-allocate when size is known:

```go
items := make([]Item, 0, estimatedSize)
```

- Check map existence properly:

```go
v, ok := m[key]
```

- Do not modify slices while ranging over them.

---

## Testing Idioms

- Test behavior, not implementation.
- Table-driven tests with descriptive names:

```go
name: "valid input returns expected result"
```

- Avoid `t.Fatal` inside loops.
- Use `cmp.Diff` or `reflect.DeepEqual` consistently.
- Tests must not depend on execution order.

---

## Logging (If Used)

- Do not log inside core business logic.
- Log at boundaries (UI / CLI / App layer).
- Logs must be structured and actionable.

---

## Comments & Documentation

- Comments explain **why**, not **what**.
- Avoid redundant comments:

```go
i++ // increment i ❌
```

- Exported comments must start with identifier name.
- Use TODO with owner & reason:

```go
// TODO(username): implement feature X
```

---

## Defensive Programming

- Validate inputs at package boundary.
- Never trust external data types.
- Fail fast on schema mismatch.
- Prefer explicit errors over silent correction.

---

## Build & Tooling Practices

- `go.mod` must be tidy:

```bash
go mod tidy
```

- CI must fail on:
  - lint
  - test
  - formatting

- **Linting**: `.golangci.yml` MUST use version 2 schema (`version: "2"`).
  - Use kebab-case for linter settings (e.g., `ignore-sigs`, `ignore-package-globs`).
  - Exclusions must be configured under `linters: exclusions: rules` instead of `issues: exclude-rules`.
  - Prefer global exclusions in config over redundant `//nolint` comments in test files.

- Avoid build tags unless justified.

---

## Reference Links

### Official Go Documentation
- Reference: https://go.dev/doc
- Primary source for Go syntax, tooling, modules.

### Effective Go
- Reference: https://go.dev/doc/effective_go
- Idiomatic Go practices. All `internal/` packages must comply.

### Go Modules
- Reference: https://go.dev/ref/mod
- Dependency management. Avoid unnecessary `replace` directives.

### Go Testing
- Reference: https://go.dev/doc/testing
- Standard patterns for unit tests, benchmarks, coverage.

### Go Context
- Reference: https://pkg.go.dev/context
- Mandatory for cancellation and timeouts in I/O operations.

### Go Error Handling
- Reference: https://go.dev/blog/error-handling-and-go
- Errors are values. Wrap errors; avoid panic in business logic.

### Go Concurrency
- Reference: https://go.dev/doc/effective_go#concurrency
- Use goroutines and channels deliberately.

### Go Standard Library
- Reference: https://pkg.go.dev/std
- Prefer stdlib before third-party dependencies.

### Wails Desktop App
- Reference: https://github.com/wailsapp/wails
- Follow Wails v2 patterns for Go-to-frontend binding.

### Linting
- Reference: https://github.com/golangci/golangci-lint
- Run before committing. Fix all issues to pass CI.
