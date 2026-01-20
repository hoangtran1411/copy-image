# Learning from Copy Image Tool 🎓

If you're new to Go or building desktop applications with Go, this project is designed to be a great practical example. Here's a breakdown of the Go concepts and patterns you can learn by exploring this codebase.

## 📋 Table of Contents

- [1. Project Architecture](#1-project-architecture)
- [2. Blazing Fast Concurrency (Worker Pool)](#2-blazing-fast-concurrency-worker-pool)
- [3. Go + Wails: Desktop Apps with Web Tech](#3-go--wails-desktop-apps-with-web-tech)
- [4. Configuration & Domain Logic](#4-configuration--domain-logic)
- [5. Clean Error Handling](#5-clean-error-handling)
- [6. Real-world Infrastructure (Auto-Updates)](#6-real-world-infrastructure-auto-updates)
- [7. Testing in Go](#7-testing-in-go)

---

## 1. Project Architecture

This project follows the **Standard Go Project Layout**.
- `cmd/`: Entry points for applications.
- `internal/`: Private code you don't want others to import. This is a Go convention to keep your API clean.
- `frontend/`: All the UI code (HTML/CSS/JS).

**Lesson**: Keeping your core logic in `internal/` ensures that your business rules are separated from the presentation layer (CLI vs GUI).

---

## 2. Blazing Fast Concurrency (Worker Pool)

The heart of this tool is its ability to copy thousands of files in parallel. 
Check out `internal/copier/copier.go`:

```go
// We use a semaphore pattern to limit concurrency
semaphore := make(chan struct{}, workers)

for _, file := range files {
    go func(f string) {
        semaphore <- struct{}{}        // Acquire a slot
        defer func() { <-semaphore }() // Release the slot
        copyFile(f, dest)
    }(file)
}
```

**What you'll learn**:
-   How to use **Goroutines** for non-blocking execution.
-   How to use **Channels** as semaphores to prevent overwhelming the system (CPU/Disk/Network).
-   How to aggregate results using `sync.WaitGroup` or atomic counters.

---

## 3. Go + Wails: Desktop Apps with Web Tech

This project uses **Wails**, which allows you to write your backend in Go and your frontend in standard web technologies.

- **Binding**: See `main_wails.go`. We "bind" a Go struct to the frontend.
- **Calling Go from JS**: Look at `frontend/dist/app.js`. You literally call Go methods as if they were async JavaScript functions.
- **Events**: Go can push data to the UI using `runtime.EventsEmit`. This is how the real-time progress bar works!

---

## 4. Configuration & Domain Logic

In `internal/config/config.go`, we manage complex configurations.

- **YAML/JSON Tags**: See how the same struct can be parsed from a file and sent to the UI.
- **Defaults & Validation**: Always provide sane defaults and validate user input immediately.

---

## 5. Clean Error Handling

Go treats errors as values. 
**Pattern**: "Wrap" errors to provide context as they move up the stack.

```go
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}
```

**Lesson**: Using `%w` allows the caller to inspect the root cause while seeing the full history of what went wrong.

---

## 6. Real-world Infrastructure (Auto-Updates)

The `updater.go` file demonstrates how to build professional-grade features:
- **GitHub API Integration**: Fetching the latest releases dynamically.
- **Semantic Versioning**: Comparing strings like `v1.2.3` correctly.
- **OS Integration**: Using a temporary `.bat` script to update an executable that is currently running (a common Windows challenge!).

---

## 7. Testing in Go

Go has a built-in testing framework. Check out `updater_test.go` and `internal/config/config_test.go`.

- **Table-Driven Tests**: A clean way to test multiple scenarios for the same function.
- **Mocking**: How to test logic that normally interacts with the filesystem or network by isolating it.

---

## 💡 Pro-Tip for Beginners

Don't just read the code—**break it!**
1.  Try changing the `workers` count and see how it affects copy speed.
2.  Add a new button in the UI and bind it to a new Go function.
3.  Modify the progress bar colors in `style.css`.

Happy Learning! If you find this project helpful for your Go journey, consider giving it a ⭐ on GitHub!
