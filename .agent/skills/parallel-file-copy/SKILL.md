---
name: Parallel File Copy
description: Implement high-performance parallel file copying with progress tracking and cancellation.
---

# Parallel File Copy Skill 📁

This skill provides patterns for implementing high-performance file copying operations with parallel workers, progress tracking, retry logic, and graceful cancellation.

## When to Use

- Copying large numbers of files
- Implementing backup/sync utilities
- Adding progress indicators to file operations
- Need cancellation support for long-running copies

## Core Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Scanner   │────▶│  Worker Pool │────▶│  Reporter   │
│ (GetFiles)  │     │  (Parallel)  │     │ (Progress)  │
└─────────────┘     └──────────────┘     └─────────────┘
                           │
                    ┌──────┴──────┐
                    │  Semaphore  │
                    │ (Limit N)   │
                    └─────────────┘
```

## Key Components

### 1. Copier Struct

```go
type Copier struct {
    config  *config.Config
    results []CopyResult
}

type CopyResult struct {
    FileName string
    Success  bool
    Skipped  bool
    Error    error
}

type CopySummary struct {
    TotalFiles  int
    Successful  int
    Failed      int
    Skipped     int
    Duration    time.Duration
    FailedFiles []string
}
```

### 2. Semaphore Pattern for Worker Limiting

```go
func (c *Copier) CopyFilesParallel(files []string) CopySummary {
    var (
        successful int32
        failed     int32
        wg         sync.WaitGroup
    )
    
    // Create semaphore channel to limit concurrent workers
    semaphore := make(chan struct{}, c.config.Workers)
    
    for _, file := range files {
        wg.Add(1)
        go func(f string) {
            defer wg.Done()
            
            // Acquire worker slot
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            // Do the copy work
            result := c.CopyFileWithRetry(f)
            
            if result.Success {
                atomic.AddInt32(&successful, 1)
            } else {
                atomic.AddInt32(&failed, 1)
            }
        }(file)
    }
    
    wg.Wait()
    return CopySummary{...}
}
```

### 3. Context-Based Cancellation

```go
func (c *Copier) CopyFilesParallelWithEvents(
    ctx context.Context,
    files []string,
    onProgress ProgressCallback,
) CopySummary {
    
    for _, file := range files {
        // Check for cancellation before starting new work
        select {
        case <-ctx.Done():
            break  // Stop processing new files
        default:
            // Continue
        }
        
        wg.Add(1)
        go func(f string) {
            defer wg.Done()
            
            // Also check when acquiring worker slot
            select {
            case semaphore <- struct{}{}:
                defer func() { <-semaphore }()
            case <-ctx.Done():
                return  // Cancelled while waiting
            }
            
            // Do work...
        }(file)
    }
    
    wg.Wait()
    return summary
}
```

### 4. Retry Logic with Exponential Backoff

```go
func (c *Copier) CopyFileWithRetry(sourcePath string) CopyResult {
    var lastErr error
    
    for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
        err := c.CopyFile(sourcePath, c.config.Overwrite)
        if err == nil {
            return CopyResult{Success: true}
        }
        
        lastErr = err
        
        // Exponential backoff: 100ms, 200ms, 300ms...
        if attempt < c.config.MaxRetries {
            time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
        }
    }
    
    return CopyResult{Success: false, Error: lastErr}
}
```

### 5. Progress Callback Pattern

```go
type ProgressCallback func(current int, total int, fileName string, status string)

// Usage in parallel copy
current := int(atomic.AddInt32(&processed, 1))
if onProgress != nil {
    onProgress(current, total, fileName, status)
}

// CLI mode: use progressbar
bar := progressbar.NewOptions(len(files),
    progressbar.OptionEnableColorCodes(true),
    progressbar.OptionShowCount(),
    progressbar.OptionSetWidth(40),
)

// GUI mode: emit events
runtime.EventsEmit(ctx, "copy:progress", ProgressEvent{
    Current: current,
    Total:   total,
    Percent: float64(current) / float64(total) * 100,
})
```

## Atomic Counters

Use `sync/atomic` for thread-safe counters:

```go
var (
    successful int32
    failed     int32
    skipped    int32
)

// In goroutine
atomic.AddInt32(&successful, 1)

// At the end
summary := CopySummary{
    Successful: int(successful),
    Failed:     int(failed),
}
```

## Thread-Safe Slice Appending

```go
var (
    failedFiles []string
    failedMu    sync.Mutex
)

// In goroutine
failedMu.Lock()
failedFiles = append(failedFiles, fileName)
failedMu.Unlock()
```

## File Copy Best Practices

```go
func (c *Copier) CopyFile(sourcePath string, overwrite bool) error {
    // 1. Check if destination exists (skip if not overwriting)
    if utils.FileExists(destPath) && !overwrite {
        return nil
    }
    
    // 2. Check if source is locked
    if utils.IsFileLocked(sourcePath) {
        return fmt.Errorf("file is locked by another process")
    }
    
    // 3. Ensure destination directory exists
    if err := utils.EnsureDir(c.config.Destination); err != nil {
        return fmt.Errorf("failed to create destination: %w", err)
    }
    
    // 4. Open source
    srcFile, err := os.Open(sourcePath)
    if err != nil {
        return fmt.Errorf("failed to open source: %w", err)
    }
    defer func() { _ = srcFile.Close() }()
    
    // 5. Create destination
    dstFile, err := os.Create(destPath)
    if err != nil {
        return fmt.Errorf("failed to create destination: %w", err)
    }
    defer func() { _ = dstFile.Close() }()
    
    // 6. Copy with buffered I/O
    if _, err := io.Copy(dstFile, srcFile); err != nil {
        return fmt.Errorf("failed to copy content: %w", err)
    }
    
    // 7. Sync to disk (important for network drives)
    if err := dstFile.Sync(); err != nil {
        return fmt.Errorf("failed to sync: %w", err)
    }
    
    return nil
}
```

## Configuration Options

```go
type Config struct {
    Source      string   `yaml:"source"`
    Destination string   `yaml:"destination"`
    Workers     int      `yaml:"workers"`      // Concurrent workers (1-50)
    Overwrite   bool     `yaml:"overwrite"`    // Overwrite existing files
    Extensions  []string `yaml:"extensions"`   // Filter by extension
    MaxRetries  int      `yaml:"max_retries"`  // Retry count on failure
    DryRun      bool     `yaml:"dry_run"`      // Simulate without copying
}
```

## AI Prompt Templates

- **Add feature:** "Add file filtering by date to the copier"
- **Progress:** "Implement ETA calculation for copy progress"
- **Optimize:** "Optimize copy performance for large files"
