# CopyImageDMDGL - Method Documentation

> This document summarizes the logic and main methods of the CopyImageDMDGL application (C#) to support the conversion to Go.

---

## 📋 Application Purpose

A console application for **bulk image copying** from a network share source directory to a destination directory.

**Real-world use case:**
- Copy product sample images from temporary storage server (`IMAGES PENDING UPLOAD`) to main server (`Images`)
- Support file overwrite when needed
- Parallel processing to speed up copying

---

## ⚙️ Configuration (Constants)

```
SOURCE_PATH      = "\\192.1.1.1\DM_DON_GIA_LUONG\ROUTING 2023 + HÌNH MẪU\HINHMAUSP\HÌNH CHƯA TẢI"
DESTINATION_PATH = "\\192.1.1.20\dmdgl$\HinhAnh"
```

**Suggestions for Go:**
- Use config file (JSON, YAML, TOML) or environment variables
- Support command-line flags: `--source`, `--dest`, `--overwrite`

---

## 🔄 Main Flow

```
START
  │
  ▼
┌─────────────────────────────────────┐
│ 1. Display selection menu           │
│    - 0: Don't copy (exit)           │
│    - 1: Copy and overwrite          │
└─────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────┐
│ 2. Validate input                   │
│    - Loop until valid input (0/1)   │
└─────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────┐
│ 3. Check source directory exists    │
│    - If not exists → notify         │
└─────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────┐
│ 4. Get list of files                │
│    - If empty → notify              │
└─────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────┐
│ 5. Copy in parallel                 │
│    - For each file:                 │
│      • Create destination path      │
│      • Copy file (overwrite if set) │
│      • Handle exceptions            │
│      • Log result                   │
└─────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────┐
│ 6. Complete                         │
│    - Display notification           │
│    - Wait for user keypress         │
└─────────────────────────────────────┘
  │
  ▼
END
```

---

## 📦 Main Methods

### 1. `Main()` - Entry Point

**Purpose:** Orchestrate the entire processing flow

**Logic:**
```
1. Print console menu
2. Read user input → validate (only accept 0 or 1)
3. If selected 0 → exit
4. If selected 1:
   a. Check source directory
   b. Get list of files
   c. Copy in parallel with overwrite = true
5. Print results and wait for keypress
```

**Input:** No parameters
**Output:** Console output

---

### 2. `IsFileLocked(filePath string) bool`

**Purpose:** Check if a file is locked (being opened by another process)

**Logic:**
```
1. Try to open file with ReadWrite mode and FileShare.None
2. If successful → file is not locked → return false
3. If IOException → file is locked → return true
```

**Input:** `filePath` - Absolute path to the file to check
**Output:** `bool` - `true` if file is locked, `false` otherwise

**Original C# code:**
```csharp
static bool IsFileLocked(string filePath)
{
    try
    {
        using (FileStream stream = File.Open(filePath, FileMode.Open, FileAccess.ReadWrite, FileShare.None))
        {
            return false;
        }
    }
    catch (IOException)
    {
        return true;
    }
}
```

**Suggestion for Go:**
```go
func isFileLocked(filePath string) bool {
    file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
    if err != nil {
        return true // File is locked or doesn't exist
    }
    defer file.Close()
    return false
}
```

---

### 3. `CopyFile(sourcePath, destPath string, overwrite bool) error`

**Purpose:** Copy a file from source to destination

**Logic:**
```
1. Get filename from source path
2. Create destination path = destPath + fileName
3. Copy file:
   - If overwrite = true → overwrite if exists
   - If overwrite = false → skip if exists
4. Handle exceptions:
   - File is locked → log and skip
   - Other errors → log error
```

**Input:**
- `sourcePath` - Source file path
- `destPath` - Destination directory
- `overwrite` - Whether to overwrite

**Output:** `error` or `nil`

---

### 4. `CopyFilesParallel(files []string, destPath string, overwrite bool)`

**Purpose:** Copy multiple files in parallel to increase performance

**Logic (C# uses Parallel.ForEach):**
```
1. For each file in the list (parallel):
   a. Call CopyFile()
   b. Log result: ✓ success or ✗ failure
```

**Suggestion for Go (using goroutines + WaitGroup):**
```go
func copyFilesParallel(files []string, destPath string, overwrite bool) {
    var wg sync.WaitGroup
    semaphore := make(chan struct{}, 10) // Limit concurrent goroutines
    
    for _, file := range files {
        wg.Add(1)
        go func(f string) {
            defer wg.Done()
            semaphore <- struct{}{}        // Acquire
            defer func() { <-semaphore }() // Release
            
            err := copyFile(f, destPath, overwrite)
            if err != nil {
                fmt.Printf("✗ %s: %v\n", filepath.Base(f), err)
            } else {
                fmt.Printf("✓ %s copied.\n", filepath.Base(f))
            }
        }(file)
    }
    wg.Wait()
}
```

---

## 🛡️ Error Handling

| Error Type | Handling |
|------------|----------|
| Source directory doesn't exist | Log message and exit |
| No files in directory | Log message and exit |
| File is locked | Skip, log with ✗ prefix |
| Other IOException | Log error details with ✗ prefix |
| General exception | Log and continue with next file |

---

## 🚀 Improvement Suggestions for Go

### 1. **CLI with Cobra/Flag**
```
copyimage --source "/path/to/source" --dest "/path/to/dest" --overwrite --workers 10
```

### 2. **Progress Bar**
Use library like `github.com/schollz/progressbar/v3`

### 3. **Structured Logging**
Use `log/slog` (Go 1.21+) or `zerolog`/`zap`

### 4. **Retry Mechanism**
Retry when copy fails (max 3 attempts)

### 5. **Dry-run Mode**
Option `--dry-run` to preview files that will be copied

### 6. **Filter Files**
Option `--ext .jpg,.png` to only copy certain file types

### 7. **Worker Pool**
Control the number of concurrent goroutines to avoid overload

### 8. **Report/Summary**
```
========== RESULTS ===========
Total files:    100
Successful:     95
Failed:         3
Skipped:        2
Duration:       5.2s
===============================
```

### 9. **Config File**
```yaml
# config.yaml
source: "\\\\192.1.1.1\\path\\to\\source"
destination: "\\\\192.1.1.20\\path\\to\\dest"
workers: 10
overwrite: true
extensions:
  - .jpg
  - .png
  - .gif
```

---

## 📁 Suggested Go Project Structure

```
copyimage/
├── cmd/
│   └── copyimage/
│       └── main.go           # Entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Load config
│   ├── copier/
│   │   ├── copier.go         # Copy logic
│   │   └── copier_test.go    # Unit tests
│   └── utils/
│       └── filelock.go       # IsFileLocked helper
├── config.yaml               # Default config
├── go.mod
├── go.sum
└── README.md
```

---

## ✅ Conversion Checklist

- [x] Create new Go project with `go mod init`
- [x] Implement `config` package (load from file/flags/env)
- [x] Implement `isFileLocked()` function
- [x] Implement `copyFile()` function
- [x] Implement `copyFilesParallel()` with worker pool
- [x] Add CLI flags (cobra or flag package)
- [x] Add progress bar
- [x] Add summary report
- [x] Write unit tests
- [x] Build and test on Windows with UNC paths
- [x] Add Wails desktop application
- [x] Add auto-update functionality
- [x] Add Copy Groups support

---

*Document created: 2026-01-19*
*Last updated: 2026-01-20 - Added Wails desktop app features*
