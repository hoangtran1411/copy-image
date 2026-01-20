---
trigger: always_on
---

# Go Style Guide - Copy Image Project

This project is a **file copy utility** built with:
- **Wails v2** for the Desktop GUI (Windows)
- **CLI mode** with progressbar for terminal usage
- **YAML configuration** for persistent settings

---

## Code Style

- "Ensure all Go code is formatted using `gofmt` or `goimports`. Run `golangci-lint run ./...` before committing."
- "Adhere to [Effective Go](https://golang.org/doc/effective_go.html) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)."
- "Organize code into domain-specific packages within `internal/` (e.g., `internal/config`, `internal/copier`, `internal/utils`)."
- "Keep Wails-specific code (app.go, updater.go, main_wails.go) in the root package with `//go:build windows` constraint."

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

- "Always wrap errors using `%w`: `fmt.Errorf(\"context: %w\", err)`. This is critical for tracing file I/O errors."
- "Implement 'fail fast' logic using guard clauses to minimize indentation."
- "Handle close errors in `defer` statements where appropriate, or use `defer func() { _ = file.Close() }()` pattern."

## Context & Concurrency

- "Functions performing I/O or long-running operations MUST accept `context.Context` as the first argument for cancellation support."
- "Use `context` to manage timeouts and cancellations for copy operations (see `CopyFilesParallelWithEvents`)."
- "Use `sync/atomic` for counters shared across goroutines."
- "Use semaphore pattern (`chan struct{}`) to limit concurrent workers."

## Wails Integration

- "All Wails-bound methods must be on the `*App` struct and be exported (PascalCase)."
- "Use `runtime.EventsEmit()` for progress updates to frontend."
- "Return structs with `json` tags for frontend consumption (e.g., `ProgressEvent`, `CopyResult`)."
- "Files with Wails bindings require `//go:build windows` constraint."

## CLI Mode (progressbar)

- "Use `github.com/schollz/progressbar/v3` for terminal progress display."
- "Keep CLI logic in `cmd/copyimage/` separate from core business logic."
- "Support both interactive and non-interactive modes."

## Configuration

- "Use `gopkg.in/yaml.v3` for configuration serialization."
- "Always provide `DefaultConfig()` function with sensible defaults."
- "Add both `yaml` and `json` struct tags for dual compatibility (file storage + Wails binding)."
- "Validate configuration before use with `Config.Validate()` method."

## Documentation

- "Every exported function, variable, and type must have clear documentation comments explaining 'Why' rather than just 'What'."
- "Document edge cases and design decisions in comments (e.g., why certain linter rules are disabled)."

## Testing

- "Prioritize Table-driven tests combined with `t.Run` for comprehensive test coverage."
- "Use build tags (`//go:build windows`) for Windows-specific tests."
- "Mock file system operations for unit tests when possible."
- "Target minimum 40% code coverage (CI enforced)."

## Linting (golangci-lint)

The project uses these linters (see `.golangci.yml`):
- `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`
- `gofmt`, `goimports`, `misspell`, `bodyclose`, `gocritic`, `gosec`

Excluded patterns:
- `frontend/` and `build/` directories are excluded
- Wails-specific files (`app.go`, `updater.go`, etc.) are excluded from CI lint due to Windows-only constraints

## Efficiency & Tone

- "Avoid greetings, apologies, or meta-commentary; focus strictly on code and execution logs."
- "Provide code as minimal diffs/blocks whenever possible."

---

## Reference & Resource Mapping

### Wails Desktop App
- **Reference:** [wailsapp/wails](https://github.com/wailsapp/wails)
- **Guideline:** Follow Wails v2 patterns for Go-to-frontend binding. Use event system for real-time updates.

### CLI Progress Display
- **Reference:** [schollz/progressbar](https://github.com/schollz/progressbar)
- **Guideline:** Use progressbar for CLI mode with proper terminal detection.

### Configuration
- **Reference:** [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)
- **Guideline:** Use YAML for human-readable configuration files.

### Linting
- **Reference:** [golangci/golangci-lint](https://github.com/golangci/golangci-lint)
- **Guideline:** Run `golangci-lint run ./...` before committing. Fix all issues to pass CI.
