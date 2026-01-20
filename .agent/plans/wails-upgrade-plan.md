# 🚀 Implementation Plan: Upgrade Copy Image to Wails Desktop App

> **Goal**: Convert the current CLI application into a modern desktop application with a graphical user interface (GUI) using the Wails framework.

---

## 📋 Existing Project Overview

### Existing Structure
```
copy-image/
├── cmd/copyimage/main.go      # CLI Entry point (231 lines)
├── internal/
│   ├── config/config.go       # Config loading & validation (86 lines)
│   ├── copier/copier.go       # Core copy logic (267 lines)
│   └── utils/filelock.go      # File utilities (40 lines)
├── config.yaml                # Configuration file
└── go.mod                     # Go 1.23
```

### Reusable Key Components
| Component | Description | Reusable |
|-----------|-------------|----------|
| `config.Config` | Configuration structure with YAML parsing | ✅ 100% |
| `copier.Copier` | File copy logic with worker pool | ✅ 90% (needs events) |
| `copier.CopySummary` | Copy statistics results | ✅ 100% |
| `utils/*` | File utilities | ✅ 100% |

---

## 🎯 New Features with Wails

### CLI vs Desktop App Comparison

| Feature | Current CLI | Wails Desktop |
|---------|-------------|---------------|
| Folder Selection | Manual path entry | 📁 Native folder picker dialog |
| Progress | Text progress bar | 🎨 Real-time animated progress bar |
| Configuration | YAML file | ⚙️ Settings UI with form inputs |
| Interaction | Terminal commands | 🖱️ Buttons, dropdowns, checkboxes |
| Results | Print to console | 📊 Visual summary with results |
| Notifications | None | 🔔 Desktop notifications |
| Dark mode | None | 🌙 Native dark/light mode |
| Drag & Drop | None | 📥 Drag and drop folders into the app |

---

## 📐 New Architecture

```
copy-image/
├── app.go                     # Wails app struct & bindings (NEW)
├── main_wails.go              # Wails entry point (NEW)
├── frontend/                  # Web-based UI (NEW)
│   └── dist/                  # Built frontend assets
│       ├── index.html
│       ├── style.css
│       └── app.js
├── internal/
│   ├── config/config.go       # (kept)
│   ├── copier/
│   │   ├── copier.go          # (updated with events)
│   └── utils/filelock.go      # (kept)
├── wails.json                 # Wails configuration (NEW)
└── go.mod                     # Updated dependencies
```

---

## 🎯 Special Feature: Copy Groups

### Description
Allows creating **Copy Groups** - each group has 1 source and multiple destinations. Helps users copy images from one source folder to multiple destination folders simultaneously.

### Use Cases
- Copying product images from a common folder to multiple different servers
- Simultaneous backup to multiple drives/network shares
- Distributing assets to multiple environments (dev, staging, production)

### Data Structure

```go
// internal/config/config.go

// CopyGroup represents a copy configuration with one source and multiple destinations
type CopyGroup struct {
    ID           string        `yaml:"id" json:"id"`
    Name         string        `yaml:"name" json:"name"`
    Source       string        `yaml:"source" json:"source"`
    Destinations []Destination `yaml:"destinations" json:"destinations"`
    Enabled      bool          `yaml:"enabled" json:"enabled"`
}

// Destination represents a single destination with its own settings
type Destination struct {
    ID        string `yaml:"id" json:"id"`
    Path      string `yaml:"path" json:"path"`
    Overwrite bool   `yaml:"overwrite" json:"overwrite"`
    Enabled   bool   `yaml:"enabled" json:"enabled"`
}

// Config represents the application configuration
type Config struct {
    // Legacy single source/dest (for backward compatibility)
    Source      string   `yaml:"source"`
    Destination string   `yaml:"destination"`
    
    // New: Copy Groups
    Groups []CopyGroup `yaml:"groups" json:"groups"`
    
    // Global settings
    Workers    int      `yaml:"workers"`
    Extensions []string `yaml:"extensions"`
    MaxRetries int      `yaml:"max_retries"`
    DryRun     bool     `yaml:"dry_run"`
}
```

### Config YAML Example

```yaml
# config.yaml - New configuration with Groups

# Global settings
workers: 10
extensions:
  - .jpg
  - .jpeg
  - .png
  - .gif
max_retries: 3
dry_run: false

# Copy Groups - 1 source → multiple destinations
groups:
  - id: "group-1"
    name: "📷 Product Images"
    source: "\\\\192.1.1.1\\DM_DON_GIA_LUONG\\HINHMAUSP\\PENDING_UPLOAD"
    enabled: true
    destinations:
      - id: "dest-1"
        path: "\\\\192.1.1.20\\dmdgl$\\HinhAnh"
        overwrite: true
        enabled: true
      - id: "dest-2"
        path: "\\\\192.1.1.30\\backup\\HinhAnh"
        overwrite: false
        enabled: true
      - id: "dest-3"
        path: "D:\\LocalBackup\\HinhAnh"
        overwrite: true
        enabled: false  # Temporarily disabled

  - id: "group-2"
    name: "📁 Technical Docs"
    source: "\\\\192.1.1.1\\Docs\\Technical"
    enabled: true
    destinations:
      - id: "dest-4"
        path: "\\\\192.1.1.20\\dmdgl$\\Docs"
        overwrite: true
        enabled: true
```

### UI Design for Copy Groups

```
┌─────────────────────────────────────────────────────────────────┐
│  📷 Copy Image Tool v2.0                              [−][□][×] │
├─────────────────────────────────────────────────────────────────┤
│  [Copy] [Groups] [Settings] [About]                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌── Copy Groups ─────────────────────────────────────────────┐ │
│  │                                                             │ │
│  │  ☑ 📷 Product Images                            [Edit][🗑️] │ │
│  │    └─ Source: \\192.1.1.1\...\PENDING_UPLOAD               │ │
│  │    └─ Destinations:                                         │ │
│  │       ☑ \\192.1.1.20\dmdgl$\HinhAnh (overwrite: ✓)        │ │
│  │       ☑ \\192.1.1.30\backup\HinhAnh (overwrite: ✗)        │ │
│  │       ☐ D:\LocalBackup\HinhAnh (disabled)                  │ │
│  │                                                             │ │
│  │  ☑ 📁 Technical Docs                            [Edit][🗑️] │ │
│  │    └─ Source: \\192.1.1.1\Docs\Technical                   │ │
│  │    └─ Destinations:                                         │ │
│  │       ☑ \\192.1.1.20\dmdgl$\Docs (overwrite: ✓)            │ │
│  │                                                             │ │
│  │                                    [+ Add New Group]        │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌── Actions ─────────────────────────────────────────────────┐ │
│  │  ○ Copy selected groups    ○ Copy all enabled groups       │ │
│  │                                                             │ │
│  │  [🔍 Scan Files]  [🚀 Start Copy]  [⏹️ Cancel]              │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                 │
│  ┌── Progress ────────────────────────────────────────────────┐ │
│  │  Group: Product Images                                     │ │
│  │  Dest: \\192.1.1.20\dmdgl$\HinhAnh                         │ │
│  │  [████████████░░░░░░░░] 65% (130/200 files)                │ │
│  │  Current: product_12345.jpg                                 │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Backend API for Groups

```go
// app.go - Wails bindings

// Group Management
func (a *App) GetGroups() []config.CopyGroup
func (a *App) AddGroup(group config.CopyGroup) error
func (a *App) UpdateGroup(group config.CopyGroup) error
func (a *App) DeleteGroup(groupID string) error
func (a *App) ToggleGroup(groupID string, enabled bool) error

// Destination Management
func (a *App) AddDestination(groupID string, dest config.Destination) error
func (a *App) UpdateDestination(groupID, destID string, dest config.Destination) error
func (a *App) DeleteDestination(groupID, destID string) error
func (a *App) ToggleDestination(groupID, destID string, enabled bool) error

// Copy Operations
func (a *App) ScanGroupFiles(groupID string) ([]string, error)
func (a *App) StartGroupCopy(groupID string) error
func (a *App) StartAllGroupsCopy() error
func (a *App) CancelCopy() error
```

### Copy Logic with Groups

```go
// internal/copier/group_copier.go

type GroupCopyResult struct {
    GroupID      string                   `json:"groupId"`
    GroupName    string                   `json:"groupName"`
    Destinations []DestinationCopyResult  `json:"destinations"`
    TotalFiles   int                      `json:"totalFiles"`
    Duration     time.Duration            `json:"duration"`
}

type DestinationCopyResult struct {
    DestID     string `json:"destId"`
    DestPath   string `json:"destPath"`
    Successful int    `json:"successful"`
    Failed     int    `json:"failed"`
    Skipped    int    `json:"skipped"`
}

// CopyGroup copies files from source to all enabled destinations
func (c *Copier) CopyGroup(ctx context.Context, group config.CopyGroup, files []string) GroupCopyResult {
    result := GroupCopyResult{
        GroupID:    group.ID,
        GroupName:  group.Name,
        TotalFiles: len(files),
    }

    for _, dest := range group.Destinations {
        if !dest.Enabled {
            continue
        }

        // Emit event: starting destination copy
        runtime.EventsEmit(ctx, "copy:dest-start", map[string]any{
            "groupId": group.ID,
            "destId":  dest.ID,
            "destPath": dest.Path,
        })

        destResult := c.copyToDestination(ctx, files, dest)
        result.Destinations = append(result.Destinations, destResult)
    }

    return result
}
```

### Progress Events Structure

```go
// Events sent to the frontend

// When starting group copy
type GroupStartEvent struct {
    GroupID   string   `json:"groupId"`
    GroupName string   `json:"groupName"`
    DestCount int      `json:"destCount"`
    FileCount int      `json:"fileCount"`
}

// When starting copy to a destination
type DestStartEvent struct {
    GroupID   string `json:"groupId"`
    DestID    string `json:"destId"`
    DestPath  string `json:"destPath"`
    FileCount int    `json:"fileCount"`
}

// File progress update
type FileProgressEvent struct {
    GroupID   string  `json:"groupId"`
    DestID    string  `json:"destId"`
    FileName  string  `json:"fileName"`
    Current   int     `json:"current"`
    Total     int     `json:"total"`
    Percent   float64 `json:"percent"`
    Status    string  `json:"status"` // "copying", "success", "failed", "skipped"
}

// When a destination is complete
type DestCompleteEvent struct {
    GroupID    string `json:"groupId"`
    DestID     string `json:"destId"`
    Successful int    `json:"successful"`
    Failed     int    `json:"failed"`
    Skipped    int    `json:"skipped"`
}

// When the entire group copy is complete
type GroupCompleteEvent struct {
    GroupID   string        `json:"groupId"`
    GroupName string        `json:"groupName"`
    Duration  time.Duration `json:"duration"`
    Results   []DestinationCopyResult `json:"results"`
}
```

---

## 🔄 Auto-Update Feature (Referenced from GoExcelImageImporter)

### Description
Automatically checks and updates to the latest version from GitHub Releases. This is a great feature from the [GoExcelImageImporter](https://github.com/hoangtran1411/GoExcelImageImporter) project.

### Implementation: `updater.go`

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strings"

    "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Current version - update this when releasing new versions, or use ldflags
var CurrentVersion = "v1.0.0"

// GitHub repository info
const (
    GitHubOwner = "hoangtran1411"
    GitHubRepo  = "copy-image"
)

// UpdateInfo holds information about available updates
type UpdateInfo struct {
    Available   bool   `json:"available"`
    CurrentVer  string `json:"currentVersion"`
    LatestVer   string `json:"latestVersion"`
    DownloadURL string `json:"downloadUrl"`
    ReleaseURL  string `json:"releaseUrl"`
}

// GitHubRelease represents a GitHub release API response
type GitHubRelease struct {
    TagName string `json:"tag_name"`
    HTMLURL string `json:"html_url"`
    Assets  []struct {
        Name               string `json:"name"`
        BrowserDownloadURL string `json:"browser_download_url"`
    } `json:"assets"`
}

// GetCurrentVersion returns the current app version
func (a *App) GetCurrentVersion() string {
    return CurrentVersion
}

// CheckForUpdate checks GitHub for newer versions
func (a *App) CheckForUpdate() UpdateInfo {
    info := UpdateInfo{
        Available:  false,
        CurrentVer: CurrentVersion,
    }

    // Call GitHub API
    url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", GitHubOwner, GitHubRepo)
    resp, err := http.Get(url)
    if err != nil {
        return info
    }
    defer resp.Body.Close()

    if resp.StatusCode != 200 {
        return info
    }

    var release GitHubRelease
    if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
        return info
    }

    info.LatestVer = release.TagName
    info.ReleaseURL = release.HTMLURL

    // Find Windows exe asset
    for _, asset := range release.Assets {
        if strings.HasSuffix(strings.ToLower(asset.Name), ".exe") {
            info.DownloadURL = asset.BrowserDownloadURL
            break
        }
    }

    // Compare versions
    if info.LatestVer != "" && CompareVersions(info.LatestVer, CurrentVersion) {
        info.Available = true
    }

    return info
}

// CompareVersions returns true if v1 is newer than v2
func CompareVersions(v1, v2 string) bool {
    v1 = strings.TrimPrefix(v1, "v")
    v2 = strings.TrimPrefix(v2, "v")

    parts1 := parseVersion(v1)
    parts2 := parseVersion(v2)

    for i := 0; i < 3; i++ {
        if parts1[i] > parts2[i] {
            return true
        }
        if parts1[i] < parts2[i] {
            return false
        }
    }
    return false
}

func parseVersion(v string) [3]int {
    var result [3]int
    parts := strings.Split(v, ".")
    for i := 0; i < len(parts) && i < 3; i++ {
        fmt.Sscanf(parts[i], "%d", &result[i])
    }
    return result
}

// PerformUpdate downloads and installs the new version
func (a *App) PerformUpdate(downloadURL string) (bool, error) {
    if downloadURL == "" {
        return false, fmt.Errorf("no download URL provided")
    }

    exePath, err := os.Executable()
    if err != nil {
        return false, fmt.Errorf("failed to get executable path: %w", err)
    }
    exePath, _ = filepath.Abs(exePath)

    tempDir := os.TempDir()
    tempFile := filepath.Join(tempDir, "copyimage_update.exe")

    runtime.EventsEmit(a.ctx, "updateProgress", "Downloading update...")

    resp, err := http.Get(downloadURL)
    if err != nil {
        return false, fmt.Errorf("failed to download: %w", err)
    }
    defer resp.Body.Close()

    out, err := os.Create(tempFile)
    if err != nil {
        return false, fmt.Errorf("failed to create temp file: %w", err)
    }

    _, err = io.Copy(out, resp.Body)
    out.Close()
    if err != nil {
        return false, fmt.Errorf("failed to save update: %w", err)
    }

    runtime.EventsEmit(a.ctx, "updateProgress", "Installing update...")

    // Create update batch script for Windows
    batchPath := filepath.Join(tempDir, "update_copyimage.bat")
    batchContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
del "%s"
move /y "%s" "%s"
start "" "%s"
del "%%~f0"
`, exePath, tempFile, exePath, exePath)

    if err := os.WriteFile(batchPath, []byte(batchContent), 0600); err != nil {
        return false, fmt.Errorf("failed to create update script: %w", err)
    }

    cmd := exec.Command("cmd", "/c", "start", "/min", "", batchPath)
    if err := cmd.Start(); err != nil {
        return false, fmt.Errorf("failed to start update script: %w", err)
    }

    runtime.Quit(a.ctx)
    return true, nil
}
```

---

## 🎨 Design System (Referenced from GoExcelImageImporter)

### CSS Variables - Dark Mode Premium Theme

```css
:root {
    /* Background Colors */
    --bg-primary: #0f1419;
    --bg-secondary: #1a1f2e;
    --bg-card: #1e2533;
    --bg-input: #252d3d;
    --bg-hover: #2a3447;

    /* Text Colors */
    --text-primary: #e7eaf0;
    --text-secondary: #8b95a5;
    --text-muted: #5c6778;

    /* Accent Colors */
    --accent-primary: #3b82f6;
    --accent-primary-hover: #2563eb;
    --accent-success: #10b981;
    --accent-error: #ef4444;
    --accent-warning: #f59e0b;

    /* Border & Shadow */
    --border-color: #2d3748;
    --border-radius: 12px;
    --border-radius-sm: 8px;
    --shadow-sm: 0 2px 8px rgba(0, 0, 0, 0.2);
    --shadow-md: 0 4px 16px rgba(0, 0, 0, 0.3);

    /* Transitions */
    --transition: all 0.2s ease;
}
```

### Card Component with Hover Effect

```css
.card {
    background: var(--bg-card);
    border-radius: var(--border-radius);
    border: 1px solid var(--border-color);
    box-shadow: var(--shadow-sm);
    transition: var(--transition);
}

.card:hover {
    border-color: var(--accent-primary);
    box-shadow: var(--shadow-md), 0 0 0 1px var(--accent-primary);
}
```

### Toast Notification

```javascript
function showToast(message, type) {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    
    let icon = type === 'success' ? '✓' : type === 'error' ? '✗' : 'ℹ';
    toast.innerHTML = `<span class="toast-icon">${icon}</span>${message}`;
    
    container.appendChild(toast);
    
    setTimeout(() => {
        toast.classList.add('hiding');
        toast.addEventListener('transitionend', () => toast.remove());
    }, 4000);
}
```

---

## 📅 Implementation Phases

### Phase 1: Setup Wails Project
**Estimated: 2-3 hours**

- [x] **1.1** Install Wails CLI
- [x] **1.2** Initialize Wails Project
- [x] **1.3** Migrate existing `internal/` packages
- [x] **1.4** Create `app.go` with basic bindings

---

### Phase 2: Backend Bindings
**Estimated: 3-4 hours**

- [x] **2.1** Create `app.go` - Main application struct
- [x] **2.2** Implement folder selection dialog
- [x] **2.3** Implement config management
- [x] **2.4** Update `copier.go` to emit events (CopyFilesParallelWithEvents)
- [x] **2.5** Implement copy operations with events

---

### Phase 3: Frontend UI
**Estimated: 6-8 hours**

- [x] **3.1** Setup modern styling with dark mode support
- [x] **3.2** Create folder selection UI components
- [x] **3.3** Create settings UI components
- [x] **3.4** Create progress bar component with animations
- [x] **3.5** Create summary results view
- [x] **3.6** Assemble Main App layout

---

### Phase 4: Event Integration
**Estimated: 2-3 hours**

- [x] **4.1** Subscribe to backend events in frontend (progress, completion)
- [x] **4.2** Implement cancel functionality using context cancellation
- [x] **4.3** Error handling and toast notifications

---

### Phase 5: Polish & Testing
**Estimated: 3-4 hours**

- [x] **5.1** Window configuration (title, size, background color)
- [x] **5.2** App icon and branding setup
- [x] **5.3** Comprehensive testing (UNC paths, large file sets)
- [x] **5.4** Build and packaging for Windows

---

### Phase 6: Advanced Features (Optional - Future)
**Estimated: 4-6 hours**

- [ ] **6.1** Drag & Drop support
- [ ] **6.2** System tray integration
- [ ] **6.3** File preview thumbnails during copy
- [ ] **6.4** Copy history logging
- [ ] **6.5** Multiple copy queues support

---

## 📦 New Dependencies

```go
// go.mod additions
require (
    github.com/wailsapp/wails/v2 v2.9.2
)
```

---

## ✅ Definition of Done

### MVP Requirements
- [x] Native source/dest folder selection dialogs
- [x] Real-time animated progress bar
- [x] Visual summary results after copy
- [x] Settings editable directly in UI
- [x] Standalone .exe build successful
- [x] Auto-update feature functional
- [x] Copy Groups support implemented in backend

### Nice to Have
- [x] Premium Dark Mode support
- [ ] Drag & drop folder support
- [x] Toast notifications
- [ ] System tray integration

---

## 🔗 References

- [Wails Documentation](https://wails.io/docs/introduction)
- [Wails Examples](https://github.com/wailsapp/wails/tree/master/examples)
- [GoExcelImageImporter](https://github.com/hoangtran1411/GoExcelImageImporter) - Design & Update inspiration

---

## 📝 Notes

1. **Keep CLI mode**: Retain `cmd/copyimage/main.go` for automated/headless scenarios.
2. **Config compatibility**: Ensure YAML format remains backward compatible.
3. **WebView2 requirement**: Document that WebView2 is required on Windows for Wails.
