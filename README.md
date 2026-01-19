# Copy Image Tool 📷

> Công cụ sao chép hình ảnh hàng loạt từ thư mục nguồn sang thư mục đích, hỗ trợ xử lý song song để tăng tốc độ.

## ✨ Tính năng

- 🚀 **Xử lý song song** - Sử dụng worker pool để copy nhiều file đồng thời
- 📊 **Progress bar** - Hiển thị tiến trình copy trực quan
- 🔄 **Retry mechanism** - Tự động thử lại khi copy thất bại
- 📝 **Báo cáo chi tiết** - Thống kê số file thành công/thất bại/bỏ qua
- 🎯 **Lọc file** - Chỉ copy các file theo extension được chỉ định
- 🔧 **Cấu hình linh hoạt** - Hỗ trợ config file YAML và CLI flags
- 🌐 **Hỗ trợ UNC path** - Hoạt động với đường dẫn network share

## 📁 Cấu trúc project

```
copy-image/
├── cmd/
│   └── copyimage/
│       └── main.go           # Entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Load & validate config
│   ├── copier/
│   │   ├── copier.go         # Core copy logic
│   │   └── copier_test.go    # Unit tests
│   └── utils/
│       └── filelock.go       # File utility functions
├── config.yaml               # Default configuration
├── go.mod
├── go.sum
├── method.md                 # Original documentation
└── README.md
```

## 🚀 Cài đặt

### Yêu cầu
- Go 1.21 trở lên

### Build từ source

```bash
# Clone repo
git clone <repo-url>
cd copy-image

# Download dependencies
go mod tidy

# Build
go build -o copyimage.exe ./cmd/copyimage
```

## 📖 Cách sử dụng

### Chế độ Interactive (mặc định)

```bash
# Chạy với config file mặc định
./copyimage.exe

# Hoặc chỉ định config file
./copyimage.exe --config my-config.yaml
```

Chương trình sẽ hiển thị menu:
```
┌─────────────────────────────────────┐
│         LỰA CHỌN THAO TÁC           │
├─────────────────────────────────────┤
│  0: Không copy (thoát)              │
│  1: Copy và ghi đè files cũ         │
│  2: Copy và bỏ qua files đã tồn tại │
└─────────────────────────────────────┘
```

### Chế độ Command Line

```bash
# Copy với các options
./copyimage.exe \
  --source "\\192.1.1.1\share\images" \
  --dest "\\192.1.1.20\backup\images" \
  --overwrite \
  --workers 15 \
  --ext ".jpg,.png,.gif" \
  --interactive=false

# Dry-run mode (xem trước, không copy thật)
./copyimage.exe --dry-run --interactive=false

# Xem version
./copyimage.exe --version
```

### CLI Flags

| Flag | Mô tả | Mặc định |
|------|-------|----------|
| `--source` | Đường dẫn thư mục nguồn | (từ config) |
| `--dest` | Đường dẫn thư mục đích | (từ config) |
| `--overwrite` | Ghi đè file đã tồn tại | false |
| `--workers` | Số lượng worker song song | 10 |
| `--config` | Đường dẫn file config | config.yaml |
| `--dry-run` | Chế độ xem trước | false |
| `--ext` | Danh sách extension (phân cách bởi dấu phẩy) | (tất cả) |
| `--interactive` | Chế độ tương tác | true |
| `--version` | Hiển thị version | - |

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

## 📊 Kết quả mẫu

```
╔═══════════════════════════════════════════════════════════╗
║          📷 Bulk Image Copy Tool - v1.0.0                 ║
╚═══════════════════════════════════════════════════════════╝

🔍 Đang quét thư mục nguồn...
📁 Tìm thấy 100 file(s)

🚀 Bắt đầu copy files...
Copying files... [=================>          ] 75/100 7.5 it/s

========== KẾT QUẢ ==========
Tổng số files: 100
Thành công:    95 ✓
Thất bại:      3 ✗
Bỏ qua:        2 ⊘
Thời gian:     5.20s
==============================
```

## 🧪 Testing

```bash
# Run tests
go test ./...

# Run tests với coverage
go test -cover ./...

# Run tests verbose
go test -v ./...
```

## 📝 License

MIT License

## 🤝 Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
