# 🚀 Implementation Plan: Nâng cấp Copy Image lên Wails Desktop App

> **Mục tiêu**: Chuyển đổi ứng dụng CLI hiện tại thành ứng dụng desktop hiện đại với giao diện đồ họa (GUI) sử dụng Wails framework.

---

## 📋 Tổng quan Project Hiện tại

### Cấu trúc hiện có
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

### Các thành phần chính có thể tái sử dụng
| Component | Mô tả | Tái sử dụng |
|-----------|-------|-------------|
| `config.Config` | Struct cấu hình với YAML parsing | ✅ 100% |
| `copier.Copier` | Logic copy file với worker pool | ✅ 90% (cần thêm events) |
| `copier.CopySummary` | Kết quả thống kê copy | ✅ 100% |
| `utils/*` | File utilities | ✅ 100% |

---

## 🎯 Tính năng mới với Wails

### So sánh CLI vs Desktop App

| Tính năng | CLI hiện tại | Wails Desktop |
|-----------|-------|---------------|
| Chọn thư mục | Nhập path thủ công | 📁 Native folder picker dialog |
| Progress | Text progress bar | 🎨 Real-time animated progress bar |
| Cấu hình | File YAML | ⚙️ Settings UI với form inputs |
| Thao tác | Terminal commands | 🖱️ Buttons, dropdowns, checkboxes |
| Kết quả | Print to console | 📊 Visual summary với charts |
| Notifications | Không có | 🔔 Desktop notifications |
| Dark mode | Không có | 🌙 Native dark/light mode |
| Drag & Drop | Không có | 📥 Kéo thả thư mục vào app |

---

## 📐 Kiến trúc mới

```
copy-image/
├── app.go                     # Wails app struct & bindings (NEW)
├── main.go                    # Wails entry point (REPLACE)
├── frontend/                  # React/Svelte UI (NEW)
│   ├── src/
│   │   ├── App.jsx
│   │   ├── components/
│   │   │   ├── FolderSelector.jsx
│   │   │   ├── ProgressBar.jsx
│   │   │   ├── SettingsPanel.jsx
│   │   │   └── SummaryCard.jsx
│   │   └── wailsjs/          # Auto-generated bindings
│   ├── index.html
│   └── package.json
├── internal/
│   ├── config/config.go       # (giữ nguyên)
│   ├── copier/
│   │   ├── copier.go          # (cập nhật thêm events)
│   │   └── events.go          # Event emitter cho progress (NEW)
│   └── utils/filelock.go      # (giữ nguyên)
├── wails.json                 # Wails config (NEW)
└── go.mod                     # Cập nhật deps
```

---

## 🎯 Tính năng đặc biệt: Copy Groups

### Mô tả
Cho phép tạo các **Copy Group** - mỗi group có 1 source và nhiều destinations. Giúp người dùng copy hình ảnh từ 1 thư mục nguồn đến nhiều thư mục đích cùng lúc.

### Use Cases
- Copy hình ảnh sản phẩm từ folder chung đến nhiều server khác nhau
- Backup đồng thời đến nhiều ổ đĩa/network shares
- Phân phối assets đến nhiều môi trường (dev, staging, production)

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
# config.yaml - Cấu hình mới với Groups

# Global settings
workers: 10
extensions:
  - .jpg
  - .jpeg
  - .png
  - .gif
max_retries: 3
dry_run: false

# Copy Groups - 1 source → nhiều destinations
groups:
  - id: "group-1"
    name: "📷 Hình mẫu sản phẩm"
    source: "\\\\192.1.1.1\\DM_DON_GIA_LUONG\\HINHMAUSP\\HÌNH CHƯA TẢI"
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
        enabled: false  # Tạm tắt

  - id: "group-2"
    name: "📁 Tài liệu kỹ thuật"
    source: "\\\\192.1.1.1\\TaiLieu\\KyThuat"
    enabled: true
    destinations:
      - id: "dest-4"
        path: "\\\\192.1.1.20\\dmdgl$\\TaiLieu"
        overwrite: true
        enabled: true
```

### UI Design cho Copy Groups

```
┌─────────────────────────────────────────────────────────────────┐
│  📷 Copy Image Tool v2.0                              [−][□][×] │
├─────────────────────────────────────────────────────────────────┤
│  [Copy] [Groups] [Settings] [About]                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌── Copy Groups ─────────────────────────────────────────────┐ │
│  │                                                             │ │
│  │  ☑ 📷 Hình mẫu sản phẩm                         [Edit][🗑️] │ │
│  │    └─ Source: \\192.1.1.1\...\HÌNH CHƯA TẢI                │ │
│  │    └─ Destinations:                                         │ │
│  │       ☑ \\192.1.1.20\dmdgl$\HinhAnh (overwrite: ✓)        │ │
│  │       ☑ \\192.1.1.30\backup\HinhAnh (overwrite: ✗)        │ │
│  │       ☐ D:\LocalBackup\HinhAnh (disabled)                  │ │
│  │                                                             │ │
│  │  ☑ 📁 Tài liệu kỹ thuật                         [Edit][🗑️] │ │
│  │    └─ Source: \\192.1.1.1\TaiLieu\KyThuat                  │ │
│  │    └─ Destinations:                                         │ │
│  │       ☑ \\192.1.1.20\dmdgl$\TaiLieu (overwrite: ✓)        │ │
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
│  │  Group: Hình mẫu sản phẩm                                  │ │
│  │  Dest: \\192.1.1.20\dmdgl$\HinhAnh                        │ │
│  │  [████████████░░░░░░░░] 65% (130/200 files)                │ │
│  │  Current: product_12345.jpg                                 │ │
│  └─────────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### Backend API cho Groups

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

### Copy Logic với Groups

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
// Events gửi đến frontend

// Khi bắt đầu copy một group
type GroupStartEvent struct {
    GroupID   string   `json:"groupId"`
    GroupName string   `json:"groupName"`
    DestCount int      `json:"destCount"`
    FileCount int      `json:"fileCount"`
}

// Khi bắt đầu copy đến một destination
type DestStartEvent struct {
    GroupID   string `json:"groupId"`
    DestID    string `json:"destId"`
    DestPath  string `json:"destPath"`
    FileCount int    `json:"fileCount"`
}

// Progress cho mỗi file
type FileProgressEvent struct {
    GroupID   string  `json:"groupId"`
    DestID    string  `json:"destId"`
    FileName  string  `json:"fileName"`
    Current   int     `json:"current"`
    Total     int     `json:"total"`
    Percent   float64 `json:"percent"`
    Status    string  `json:"status"` // "copying", "success", "failed", "skipped"
}

// Khi hoàn thành một destination
type DestCompleteEvent struct {
    GroupID    string `json:"groupId"`
    DestID     string `json:"destId"`
    Successful int    `json:"successful"`
    Failed     int    `json:"failed"`
    Skipped    int    `json:"skipped"`
}

// Khi hoàn thành toàn bộ group
type GroupCompleteEvent struct {
    GroupID   string        `json:"groupId"`
    GroupName string        `json:"groupName"`
    Duration  time.Duration `json:"duration"`
    Results   []DestinationCopyResult `json:"results"`
}
```

---

## � Tính năng Auto-Update (Tham khảo từ GoExcelImageImporter)

### Mô tả
Tự động kiểm tra và cập nhật phiên bản mới từ GitHub Releases. Đây là tính năng rất hay từ project [GoExcelImageImporter](https://github.com/hoangtran1411/GoExcelImageImporter).

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

    // Create update batch script
    batchPath := filepath.Join(tempDir, "update_copyimage.bat")
    batchContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
del "%s"
move /y "%s" "%s"
start "" "%s"
del "%%~f0"
`, exePath, tempFile, exePath, exePath)

    if err := os.WriteFile(batchPath, []byte(batchContent), 0644); err != nil {
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

### Frontend: Update Button & Check

```javascript
// Global variable to store update info
let updateInfo = null;

// Check for updates on startup
async function checkForUpdates() {
    try {
        updateInfo = await window.go.main.App.CheckForUpdate();
        
        if (updateInfo && updateInfo.available) {
            const updateBtn = document.getElementById('updateBtn');
            updateBtn.classList.add('visible');
            updateBtn.title = `Update to ${updateInfo.latestVersion} available!`;
        }
    } catch (err) {
        console.error('Failed to check for updates:', err);
    }
}

// Perform the update
async function performUpdate() {
    if (!updateInfo || !updateInfo.downloadUrl) {
        showToast('No update information available', 'error');
        return;
    }
    
    showToast(`Downloading ${updateInfo.latestVersion}...`, 'info');
    
    try {
        await window.go.main.App.PerformUpdate(updateInfo.downloadUrl);
        showToast('Update installed! Restarting...', 'success');
    } catch (err) {
        showToast('Update failed: ' + err, 'error');
    }
}

// Listen for update progress events
runtime.EventsOn('updateProgress', function(message) {
    showToast(message, 'info');
});
```

### CSS: Update Button Animation

```css
.update-btn {
    display: none;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    padding: 0;
    background: linear-gradient(135deg, var(--accent-success), #059669);
    border: none;
    border-radius: 50%;
    cursor: pointer;
    transition: var(--transition);
    animation: pulse-glow 2s ease-in-out infinite;
}

.update-btn.visible {
    display: flex;
}

@keyframes pulse-glow {
    0%, 100% { box-shadow: 0 0 8px rgba(16, 185, 129, 0.4); }
    50% { box-shadow: 0 0 16px rgba(16, 185, 129, 0.7); }
}
```

---

## 🎨 Design System (Tham khảo từ GoExcelImageImporter)

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

### Card Component với Hover Effect

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

## �📅 Phases triển khai

### Phase 1: Setup Wails Project (Day 1)
**Ước tính: 2-3 giờ**

- [ ] **1.1** Cài đặt Wails CLI
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  wails doctor  # Kiểm tra dependencies
  ```

- [ ] **1.2** Khởi tạo project Wails
  ```bash
  # Backup code hiện tại
  git checkout -b feature/wails-upgrade

  # Khởi tạo với template React (hoặc Svelte)
  wails init -n copy-image-gui -t react-ts
  ```

- [ ] **1.3** Migrate existing `internal/` packages
  - Copy toàn bộ `internal/` folder
  - Cập nhật `go.mod` để include Wails dependency

- [ ] **1.4** Tạo file `app.go` với basic bindings
  ```go
  type App struct {
      ctx    context.Context
      config *config.Config
      copier *copier.Copier
  }

  func (a *App) GetConfig() *config.Config
  func (a *App) SaveConfig(cfg *config.Config) error
  func (a *App) SelectFolder(dialogType string) (string, error)
  func (a *App) StartCopy() error
  ```

---

### Phase 2: Backend Bindings (Day 2)
**Ước tính: 3-4 giờ**

- [ ] **2.1** Tạo `app.go` - Main application struct
  ```go
  package main

  import (
      "context"
      "copy-image/internal/config"
      "copy-image/internal/copier"
      "github.com/wailsapp/wails/v2/pkg/runtime"
  )

  type App struct {
      ctx    context.Context
      config *config.Config
  }

  func NewApp() *App {
      return &App{}
  }

  func (a *App) startup(ctx context.Context) {
      a.ctx = ctx
      a.config = config.DefaultConfig()
  }
  ```

- [ ] **2.2** Implement folder selection dialog
  ```go
  func (a *App) SelectSourceFolder() (string, error) {
      return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
          Title: "Chọn thư mục nguồn",
      })
  }

  func (a *App) SelectDestFolder() (string, error) {
      return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
          Title: "Chọn thư mục đích",
      })
  }
  ```

- [ ] **2.3** Implement config management
  ```go
  func (a *App) GetConfig() *config.Config
  func (a *App) UpdateConfig(cfg *config.Config) error
  func (a *App) LoadConfigFromFile(path string) (*config.Config, error)
  func (a *App) SaveConfigToFile(path string) error
  ```

- [ ] **2.4** Update `copier.go` để emit events
  ```go
  // Thêm event emitter để gửi progress đến frontend
  type ProgressEvent struct {
      Current   int     `json:"current"`
      Total     int     `json:"total"`
      Percent   float64 `json:"percent"`
      FileName  string  `json:"fileName"`
      Status    string  `json:"status"` // "copying", "success", "failed", "skipped"
  }

  func (c *Copier) CopyFilesParallelWithEvents(ctx context.Context, files []string) CopySummary {
      // Emit events thay vì print to console
      runtime.EventsEmit(ctx, "copy:progress", ProgressEvent{...})
  }
  ```

- [ ] **2.5** Implement copy operations với events
  ```go
  func (a *App) ScanFiles() ([]string, error)
  func (a *App) StartCopy(overwrite bool) error
  func (a *App) CancelCopy() error
  ```

---

### Phase 3: Frontend UI (Day 3-4)
**Ước tính: 6-8 giờ**

- [ ] **3.1** Setup modern styling
  - Sử dụng CSS Variables cho theming
  - Dark mode support
  - Glassmorphism effects

- [ ] **3.2** Tạo `FolderSelector` component
  ```jsx
  // Hiển thị source/dest paths với nút Browse
  <FolderSelector
    label="Thư mục nguồn"
    value={config.source}
    onChange={handleSourceChange}
    onBrowse={handleBrowseSource}
  />
  ```

- [ ] **3.3** Tạo `SettingsPanel` component
  ```jsx
  // Workers slider, extensions checkboxes, overwrite toggle
  <SettingsPanel
    workers={config.workers}
    extensions={config.extensions}
    overwrite={config.overwrite}
    dryRun={config.dryRun}
    onChange={handleConfigChange}
  />
  ```

- [ ] **3.4** Tạo `ProgressBar` component với animations
  ```jsx
  // Animated progress bar với file count
  <ProgressBar
    current={progress.current}
    total={progress.total}
    currentFile={progress.fileName}
    status={progress.status}
  />
  ```

- [ ] **3.5** Tạo `SummaryCard` component
  ```jsx
  // Hiển thị kết quả với icons
  <SummaryCard
    total={summary.totalFiles}
    success={summary.successful}
    failed={summary.failed}
    skipped={summary.skipped}
    duration={summary.duration}
  />
  ```

- [ ] **3.6** Tạo main `App.jsx` layout
  - Header với logo và version
  - Body với tabs: Copy | Settings | About
  - Footer với action buttons

---

### Phase 4: Event Integration (Day 5)
**Ước tính: 2-3 giờ**

- [ ] **4.1** Subscribe to backend events trong frontend
  ```jsx
  useEffect(() => {
    EventsOn("copy:progress", (data) => {
      setProgress(data);
    });

    EventsOn("copy:complete", (summary) => {
      setSummary(summary);
      setIsCopying(false);
    });

    return () => {
      EventsOff("copy:progress");
      EventsOff("copy:complete");
    };
  }, []);
  ```

- [ ] **4.2** Implement cancel functionality
  ```go
  // Backend: sử dụng context cancellation
  type App struct {
      cancelFunc context.CancelFunc
  }

  func (a *App) CancelCopy() {
      if a.cancelFunc != nil {
          a.cancelFunc()
      }
  }
  ```

- [ ] **4.3** Error handling và notifications
  ```jsx
  // Toast notifications cho errors
  runtime.EventsOn("copy:error", (error) => {
    showToast({ type: "error", message: error });
  });
  ```

---

### Phase 5: Polish & Testing (Day 6)
**Ước tính: 3-4 giờ**

- [ ] **5.1** Window configuration
  ```go
  wails.Run(&options.App{
      Title:            "Copy Image Tool",
      Width:            900,
      Height:           650,
      MinWidth:         600,
      MinHeight:        500,
      WindowStartState: options.Normal,
      AssetServer: &assetserver.Options{
          Assets: assets,
      },
      OnStartup: app.startup,
  })
  ```

- [ ] **5.2** App icon và branding
  - Tạo `appicon.png` (1024x1024)
  - Build icons cho các platforms

- [ ] **5.3** Testing
  - Test trên Windows 10/11
  - Test với UNC paths (network shares)
  - Test drag & drop folders
  - Test với large file sets (1000+ files)

- [ ] **5.4** Build và packaging
  ```bash
  wails build -platform windows/amd64
  ```

---

### Phase 6: Advanced Features (Optional - Future)
**Ước tính: 4-6 giờ**

- [ ] **6.1** Drag & Drop support
  ```go
  // Wails v2.9+ hỗ trợ drag & drop
  OnDragDrop: func(filenames []string) { ... }
  ```

- [ ] **6.2** System tray integration
  - Minimize to tray
  - Background copy notifications

- [ ] **6.3** File preview thumbnails
  - Hiển thị thumbnail của images đang copy

- [ ] **6.4** Copy history
  - Lưu lịch sử các lần copy
  - Quick repeat last copy

- [ ] **6.5** Multiple copy queues
  - Hỗ trợ queue nhiều tasks

---

## 📦 Dependencies mới

```go
// go.mod additions
require (
    github.com/wailsapp/wails/v2 v2.9.2
)
```

```json
// frontend/package.json
{
  "dependencies": {
    "@wailsio/runtime": "^2.0.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "lucide-react": "^0.300.0"  // Icons
  }
}
```

---

## ✅ Definition of Done

### MVP Requirements
- [x] Có thể chọn source/dest folders qua dialog
- [x] Hiển thị progress bar real-time
- [x] Hiển thị kết quả sau khi copy xong
- [x] Settings có thể edit trong UI
- [x] Build được file .exe standalone

### Nice to Have
- [ ] Dark mode support
- [ ] Drag & drop folders
- [ ] Desktop notifications
- [ ] System tray

---

## 🔗 Tài liệu tham khảo

- [Wails Documentation](https://wails.io/docs/introduction)
- [Wails Examples](https://github.com/wailsapp/wails/tree/master/examples)
- [React + Wails Template](https://github.com/wailsapp/wails/tree/master/v2/internal/frontend/templates/react-ts)

---

## 📝 Notes

1. **Giữ nguyên CLI mode**: Có thể giữ lại `cmd/copyimage/main.go` để hỗ trợ headless/automated scenarios.

2. **Config compatibility**: Đảm bảo `config.yaml` format không thay đổi để người dùng hiện tại có thể migrate dễ dàng.

3. **WebView2 requirement**: Wails trên Windows yêu cầu WebView2. Cần document hoặc bundle WebView2 installer.
