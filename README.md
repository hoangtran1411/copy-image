# Copy Image Tool 📷

![Go Version](https://img.shields.io/github/go-mod/go-version/hoangtran1411/copy-image)
![License](https://img.shields.io/github/license/hoangtran1411/copy-image)
![Build Status](https://img.shields.io/github/actions/workflow/status/hoangtran1411/copy-image/ci.yml?branch=main)

> Bulk image copy tool with parallel processing support. Available as both CLI and Desktop application.

## ✨ Features

### Core Features
- 🚀 **Parallel Processing** - Worker pool for concurrent file copying
- 📊 **Real-time Progress** - Visual progress bar with file counts
- 🔄 **Retry Mechanism** - Auto-retry with exponential backoff
- 📝 **Detailed Reports** - Statistics for success/failed/skipped files
- 🎯 **Extension Filter** - Copy only specified file types
- 🔧 **Flexible Config** - YAML config file and CLI flags support
- 🌐 **UNC Path Support** - Works with network share paths

### Desktop App Features (Wails)
- 🖥️ **Native Desktop App** - Modern GUI with dark mode theme
- 📁 **Native Dialogs** - OS folder picker dialogs
- 🔔 **Toast Notifications** - Real-time status updates
- 🔄 **Auto-Update** - Check and install updates from GitHub Releases
- 📊 **Visual Progress** - Animated progress bar with file details
- ⚙️ **Settings UI** - Configure workers, extensions, and options

## 📁 Project Structure

```
copy-image/
├── cmd/copyimage/           # CLI entry point
│   └── main.go
├── main_wails.go            # Wails desktop entry point (Windows only)
├── app.go                   # Wails app bindings (Windows only)
├── updater.go               # Auto-update functionality (Windows only)
├── frontend/                # Desktop UI (HTML/CSS/JS)
│   └── dist/
│       ├── index.html
│       ├── style.css
│       └── app.js
├── internal/
│   ├── config/config.go     # Configuration with Copy Groups
│   ├── copier/copier.go     # Core copy logic
│   └── utils/filelock.go    # File utilities
├── wails.json               # Wails configuration
├── config.yaml              # Default configuration
└── Makefile                 # Build commands
```

## 🚀 Installation

### Requirements
- Go 1.21 or later
- For Desktop App: Windows 10/11, Wails CLI v2

### Build from Source

```bash
# Clone repo
git clone https://github.com/hoangtran1411/copy-image.git
cd copy-image

# Download dependencies
go mod tidy

# Build CLI
go build -o copyimage.exe ./cmd/copyimage

# Build Desktop App (Windows only)
wails build -clean
```

## 📖 Usage

### Interactive Mode (Default)

```bash
# Run with default config file
./copyimage.exe

# Or specify a config file
./copyimage.exe --config my-config.yaml
```

The program will display a menu:
```
┌─────────────────────────────────────┐
│         SELECT OPERATION            │
├─────────────────────────────────────┤
│  0: Don't copy (exit)               │
│  1: Copy and overwrite existing     │
│  2: Copy and skip existing files    │
└─────────────────────────────────────┘
```

### Command Line Mode

```bash
# Copy with options
./copyimage.exe \
  --source "\\192.1.1.1\share\images" \
  --dest "\\192.1.1.20\backup\images" \
  --overwrite \
  --workers 15 \
  --ext ".jpg,.png,.gif" \
  --interactive=false

# Dry-run mode (preview without copying)
./copyimage.exe --dry-run --interactive=false

# Show version
./copyimage.exe --version
```

### CLI Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--source` | Source directory path | (from config) |
| `--dest` | Destination directory path | (from config) |
| `--overwrite` | Overwrite existing files | false |
| `--workers` | Number of parallel workers | 10 |
| `--config` | Path to config file | config.yaml |
| `--dry-run` | Preview mode | false |
| `--ext` | Comma-separated list of extensions | (all files) |
| `--interactive` | Interactive mode | true |
| `--version` | Show version | - |

## ⚙️ Configuration

### config.yaml

```yaml
# Source directory - Network path to copy files from
source: "\\\\192.1.1.1\\DM_DON_GIA_LUONG\\ROUTING 2023 + HÌNH MẪU\\HINHMAUSP\\HÌNH CHƯA TẢI"

# Destination directory - Network path to copy files to
destination: "\\\\192.1.1.20\\dmdgl$\\HinhAnh"

# Number of concurrent workers
workers: 10

# Whether to overwrite existing files
overwrite: true

# File extensions to include (empty = all files)
extensions:
  - .jpg
  - .jpeg
  - .png
  - .gif

# Maximum retry attempts
max_retries: 3

# Dry run mode
dry_run: false
```

## 📊 Sample Output

```
╔═══════════════════════════════════════════════════════════╗
║          📷 Bulk Image Copy Tool - v2.0.0                 ║
╚═══════════════════════════════════════════════════════════╝

🔍 Scanning source directory...
📁 Found 100 file(s)

🚀 Starting file copy...
Copying files... [==================>          ] 75/100 7.5 it/s

========== RESULTS ===========
Total files:    100
Successful:     95 ✓
Failed:         3 ✗
Skipped:        2 ⊘
Duration:       5.20s
===============================
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests verbose
go test -v ./...

# Generate coverage report
make coverage
```

## 🔧 Makefile Commands

```bash
# Build CLI
make build

# Build Wails Desktop App
make wails-build

# Run Wails in development mode
make wails-dev

# Run all tests
make test

# Run linter
make lint

# Clean build artifacts
make clean
```

## 📝 License

MIT License - see [LICENSE](LICENSE) for details.

## 🤝 Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request
