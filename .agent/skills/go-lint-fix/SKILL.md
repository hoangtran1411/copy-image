---
name: Go Lint Fix
description: Quickly diagnose and fix golangci-lint errors to ensure CI passes.
---

# Go Lint Fix Skill 🔧

This skill provides a systematic approach to diagnose, understand, and fix linting errors reported by `golangci-lint`. Essential for ensuring CI pipelines pass.

## When to Use

- CI is failing due to lint errors
- Running `golangci-lint run ./...` shows issues
- Before pushing commits to main/develop branches
- Code review feedback mentions lint issues

## Quick Diagnosis Commands

### Step 1: Run Lint Check

```bash
# Run linter and capture output
golangci-lint run ./...

# Run with more verbose output
golangci-lint run ./... --verbose

# Run on specific package
golangci-lint run ./internal/...
```

### Step 2: Identify Error Types

Common linter errors in this project:

| Linter | Issue | Solution |
|--------|-------|----------|
| `gofmt` | File not properly formatted | Run `gofmt -w <file>` |
| `goimports` | Imports not sorted | Run `goimports -w <file>` |
| `errcheck` | Unchecked error | Handle error or use `_ = func()` pattern |
| `gosec` | Security issue | Review and fix or add exclusion in `.golangci.yml` |
| `govet` | Suspicious construct | Fix code logic issue |
| `staticcheck` | Static analysis issue | Follow suggestion in error message |
| `misspell` | Typo in comments/strings | Fix the spelling |

## Auto-Fix Commands

### Format All Code

```bash
# Format all Go files
gofmt -w .

# Sort imports for all files
goimports -w .

# Or use make target if available
make fmt
```

### Fix Specific File

```bash
# Format single file
gofmt -w internal/copier/copier.go

# Check what would change (dry run)
gofmt -d internal/copier/copier.go
```

## Project-Specific Exclusions

This project excludes certain files from CI linting (see `.golangci.yml`):

- `app.go`, `updater.go`, `main_wails.go` (Windows-only, require `//go:build windows`)
- `frontend/` directory (not Go code)
- `build/` directory (generated files)

### Adding New Exclusions

If you need to exclude a new pattern, edit `.golangci.yml`:

```yaml
issues:
  exclude-files:
    - "new_windows_only.go"
  
  exclude-rules:
    - path: _test\.go
      linters:
        - gosec
```

## Workflow

1. **Run lint** → `golangci-lint run ./...`
2. **Identify issue type** → Check error message for linter name
3. **Apply fix** → Use auto-fix or manual correction
4. **Verify fix** → Re-run lint to confirm
5. **Commit** → `git commit -m "fix(lint): <description>"`

## Common Patterns

### Ignoring Errors in Defer

```go
// BAD: errcheck will complain
defer file.Close()

// GOOD: Explicitly ignore error
defer func() { _ = file.Close() }()
```

### Error Wrapping

```go
// BAD: No context for debugging
return err

// GOOD: Wrap with context
return fmt.Errorf("failed to copy file: %w", err)
```

### Handling Unused Variables

```go
// BAD: unused variable
result, err := someFunc()

// GOOD: Explicitly ignore if not needed
_, err := someFunc()
```

## AI Prompt Templates

- **Quick fix:** "Run golangci-lint and fix all errors"
- **Specific file:** "Fix lint errors in internal/copier/copier.go"
- **Explain error:** "Explain this lint error: [paste error message]"
