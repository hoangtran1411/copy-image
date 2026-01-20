---
name: Wails Desktop App
description: Patterns and best practices for building Wails v2 desktop applications with Go backend.
---

# Wails Desktop App Skill 🖥️

This skill provides patterns and best practices for building desktop applications using Wails v2 framework with Go backend and vanilla JS/HTML/CSS frontend.

## When to Use

- Creating new Wails desktop application features
- Adding Go functions to expose to frontend
- Implementing real-time progress updates via events
- Handling native dialogs (file picker, message boxes)

## Project Structure

```
project/
├── app.go              # Main app struct with Wails bindings (Windows only)
├── updater.go          # Auto-update functionality
├── main_wails.go       # Wails entry point
├── frontend/
│   ├── dist/           # Production build output
│   ├── index.html      # Main HTML file
│   ├── app.js          # Frontend JavaScript
│   └── style.css       # Styles
└── wails.json          # Wails configuration
```

## Core Patterns

### 1. App Struct with Context

```go
//go:build windows

package main

import (
    "context"
    "github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
    ctx        context.Context
    config     *config.Config
    cancelFunc context.CancelFunc  // For cancellable operations
}

func NewApp() *App {
    return &App{}
}

// startup is called when app starts - save context here
func (a *App) startup(ctx context.Context) {
    a.ctx = ctx
    // Initialize other components
}
```

### 2. Exposing Functions to Frontend

All methods on `*App` struct that are **exported** (PascalCase) are automatically available in frontend.

```go
// GetConfig returns the current configuration.
// Frontend calls: window.go.main.App.GetConfig()
func (a *App) GetConfig() *config.Config {
    return a.config
}

// UpdateConfig updates settings.
// Frontend calls: window.go.main.App.UpdateConfig(cfg)
func (a *App) UpdateConfig(cfg *config.Config) error {
    if err := cfg.Validate(); err != nil {
        return fmt.Errorf("invalid configuration: %w", err)
    }
    a.config = cfg
    return nil
}
```

### 3. JSON Tags for Frontend Communication

Always add `json` tags to structs returned to frontend:

```go
type ProgressEvent struct {
    Current  int     `json:"current"`
    Total    int     `json:"total"`
    Percent  float64 `json:"percent"`
    FileName string  `json:"fileName"`
    Status   string  `json:"status"`
}

type CopyResult struct {
    Success     bool     `json:"success"`
    Message     string   `json:"message"`
    TotalFiles  int      `json:"totalFiles"`
    Duration    float64  `json:"duration"`
}
```

### 4. Emitting Events to Frontend

Use events for real-time updates (progress bars, notifications):

```go
// Emit progress event
runtime.EventsEmit(a.ctx, "copy:progress", ProgressEvent{
    Current:  current,
    Total:    total,
    Percent:  float64(current) / float64(total) * 100,
    FileName: fileName,
    Status:   "copying",
})

// Emit completion event
runtime.EventsEmit(a.ctx, "copy:complete", result)
```

**Frontend listener:**

```javascript
runtime.EventsOn("copy:progress", (event) => {
    updateProgressBar(event.percent);
    showCurrentFile(event.fileName);
});

runtime.EventsOn("copy:complete", (result) => {
    showCompletionMessage(result);
});
```

### 5. Native Dialogs

```go
// File/Folder picker
folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
    Title: "Select Source Folder",
})

// Message dialog
runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
    Type:    runtime.InfoDialog,
    Title:   "Success",
    Message: "Operation completed!",
})
```

### 6. Cancellable Operations

```go
func (a *App) StartLongOperation() Result {
    // Create cancellable context
    ctx, cancel := context.WithCancel(a.ctx)
    a.cancelFunc = cancel
    defer func() { a.cancelFunc = nil }()
    
    // Pass context to worker
    result := doWork(ctx)
    return result
}

func (a *App) CancelOperation() {
    if a.cancelFunc != nil {
        a.cancelFunc()
        runtime.EventsEmit(a.ctx, "operation:cancelled", nil)
    }
}
```

## Build Tags

Files with Wails bindings MUST have build constraint:

```go
//go:build windows

package main
```

This ensures:
- Code only compiles on Windows
- CI can run on Linux without Wails dependencies
- Clear separation of platform-specific code

## Build Commands

```bash
# Development mode (hot reload)
wails dev

# Production build
wails build -clean

# Build with version injection
wails build -clean -ldflags "-s -w -X main.CurrentVersion=v2.1.0"
```

## Frontend Communication Cheatsheet

| Go Side | JavaScript Side |
|---------|-----------------|
| `func (a *App) GetData() Data` | `await window.go.main.App.GetData()` |
| `runtime.EventsEmit(ctx, "event", data)` | `runtime.EventsOn("event", callback)` |
| `runtime.OpenDirectoryDialog(...)` | N/A (Go-initiated) |
| `runtime.Quit(ctx)` | `runtime.Quit()` |

## AI Prompt Templates

- **Add feature:** "Add a new Wails-bound method to [do something]"
- **Progress updates:** "Implement progress events for [operation]"
- **Dialog:** "Add a file picker dialog for [purpose]"
