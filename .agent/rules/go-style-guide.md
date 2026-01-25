---
trigger: always_on
---

# Go Style Guide - Copy Image Project

> **Core Rules** - For full idioms reference, see `go-idioms-reference.md`

This project is a **file copy utility** built with:
- **Wails v2** for the Desktop GUI (Windows)
- **CLI mode** with progressbar for terminal usage
- **YAML configuration** for persistent settings

---

## Code Style

- Format with `gofmt`/`goimports`. Run `golangci-lint` (v2.8.0+) `run ./...` before commit.
- **Linting Configuration**: MUST use `golangci-lint` v2 configuration schema (v2.8.x+). 
  - Top-level `version: "2"` is mandatory.
  - Use kebab-case for all linter settings.
  - Exclusions move to `linters: exclusions: rules`.
- Adhere to [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
- Organize code into domain-specific packages within `internal/` (e.g., `internal/config`, `internal/copier`, `internal/utils`).
- Keep Wails-specific code (app.go, updater.go, main_wails.go) in the root package with `//go:build windows` constraint.

## Project Structure

```
copy-image/
├── cmd/copyimage/       # CLI entry point
├── internal/
│   ├── config/          # Configuration loading/saving (YAML)
│   ├── copier/          # Core file copy logic with parallel workers
│   └── utils/           # Utility functions (file checks, etc.)
├── frontend/            # Wails frontend (HTML/CSS/JS)
├── app.go               # Wails app bindings (Windows only)
├── updater.go           # Auto-update functionality (Windows only)
├── main_wails.go        # Wails entry point (Windows only)
├── config.yaml          # User configuration file
└── .golangci.yml        # Linter configuration
```

## Error Handling

- Always wrap errors using `%w`: `fmt.Errorf("context: %w", err)`. Critical for tracing file I/O errors.
- Implement 'fail fast' logic using guard clauses.
- Handle close errors in `defer` statements: `defer func() { _ = file.Close() }()`
- Do not log and return the same error.

## Context & Concurrency

- Functions performing I/O or long-running operations MUST accept `context.Context` as the first argument.
- Use `context` to manage timeouts/cancellations (e.g., `CopyFilesParallelWithEvents`).
- Use `sync/atomic` for counters shared across goroutines.
- Use semaphore pattern (`chan struct{}`) to limit concurrent workers.

## Wails Integration

- Bound methods on `*App` struct must be exported (PascalCase).
- Use `runtime.EventsEmit()` for progress updates to frontend.
- Return structs with `json` tags for frontend consumption.
- Files with Wails bindings require `//go:build windows` constraint.

## CLI Mode (progressbar)

- Use `github.com/schollz/progressbar/v3` for terminal progress display.
- Keep CLI logic in `cmd/copyimage/` separate from core business logic.
- Support both interactive and non-interactive modes.

## Configuration

- Use `gopkg.in/yaml.v3` for configuration serialization.
- Always provide `DefaultConfig()` function with sensible defaults.
- Add both `yaml` and `json` struct tags.
- Validate configuration via `Config.Validate()`.

## Testing & Linting

- Table-driven tests with `t.Run`.
- Target 40% coverage (CI enforced).
- Use build tags (`//go:build windows`) for Windows-specific tests.
- Mock file system operations for unit tests.
- Lint with `golangci-lint`: `errcheck`, `gosimple`, `govet`, `staticcheck`, `bodyclose`, `gosec`, etc.

---

## AI Agent Rules (Critical)

### Enforcement

- Prefer clarity over cleverness
- Prefer idiomatic Go over Java/C#/JS patterns
- If unsure, follow Effective Go first

### Context Accuracy

- Documentation links ≠ guarantees of correctness
- For external APIs: prefer explicit function signatures in context
- State assumptions when context is missing

### Library Version Awareness

- Check `go.mod` for actual versions before suggesting APIs
- LLMs hallucinate APIs for newer features not in training data
- Prefer stable APIs over experimental features

### Context Engineering

- Right context at right time, not all docs at once
- Reference existing codebase patterns first
- State missing context rather than guessing

---

## Quick Reference Links

- [Effective Go](https://go.dev/doc/effective_go)
- [Wails v2](https://github.com/wailsapp/wails)
- [schollz/progressbar](https://github.com/schollz/progressbar)
- [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)
- [golangci-lint](https://github.com/golangci/golangci-lint)

> **Full Reference:** See `.agent/rules/go-idioms-reference.md` for detailed idioms, code examples, and best practices.
