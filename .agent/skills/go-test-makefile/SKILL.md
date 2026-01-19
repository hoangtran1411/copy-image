---
name: Go Test & Makefile
description: Standard setup for Go projects with Makefile automation and Unit Testing patterns.
---

# Go Test & Makefile Skill

This skill provides a professional setup for Go projects, enabling automation of common tasks (build, test, lint) via a `Makefile` and standardized Unit Testing patterns.

## When to Use
- Starting a new Go project.
- Adding CI/CD readiness to an existing project.
- Wanting to standardize build and test commands across team members.

## Components

### 1. Makefile
A comprehensive `Makefile` that handles:
- **Build**: Compiles the application for multiple platforms (Windows/Linux/Mac).
- **Test**: Runs unit tests with verbose output.
- **Coverage**: Generates HTML coverage reports.
- **Lint**: Runs `golangci-lint` for code quality.
- **Clean**: Removes build artifacts.

### 2. Test Template
A `main_test.go` template showing how to:
- Use `testing` package.
- Setup/Teardown with `TestMain`.
- Write table-driven tests (idiomatic Go).
- Mock dependencies (basic interfaces).

### 3. GitHub Actions
Automated CI/CD workflows including:
- **`ci.yml`**: Main CI workflow with test, coverage, lint, and build jobs.
- **`release.yml`**: Automatic release creation when tags are pushed.
- **`dependabot.yml`**: Automatic dependency updates.

## Usage

### Step 1: Add Makefile
Copy the `templates/Makefile` to your project root.
Update the `BINARY_NAME` variable at the top of the file to match your project name.

### Step 2: Add Test File
Copy `templates/main_test.go` to your package directory (e.g., `internal/utils` or root) and rename it to `yourfile_test.go`.

### Step 3: Add GitHub Actions
Copy the `.github` folder to your project root. The workflows will:
- Run on every push/PR to main branches
- Check test coverage >= 70%
- Build binaries for Windows, Linux, macOS
- Create releases when you push version tags (e.g., `v1.0.0`)

### Step 4: Run Commands
Open your terminal and run:

```bash
# Run tests
make test

# View coverage
make coverage

# Build app
make build
```

## GitHub Actions Workflows

### CI Workflow (`ci.yml`)
Triggers on push/PR to `main`, `master`, `develop` branches.

| Job | Description |
|-----|-------------|
| `test` | Run all unit tests |
| `coverage` | Check coverage >= 70% |
| `lint` | Run golangci-lint |
| `build` | Build for Windows, Linux, macOS |

### Release Workflow (`release.yml`)
Triggers when you push a version tag:

```bash
git tag v1.0.0
git push origin v1.0.0
```

This will:
1. Run tests
2. Build binaries for all platforms
3. Create GitHub Release with binaries and checksums

### Dependabot (`dependabot.yml`)
Automatically creates PRs to update:
- Go module dependencies (weekly)
- GitHub Actions versions (weekly)

## Prerequisites
- **Go**: 1.20+
- **Make**:
    - **Linux/Mac**: Pre-installed.
    - **Windows**: Install via `choco install make` or use Git Bash.
- **GolangCI-Lint** (Optional): `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

