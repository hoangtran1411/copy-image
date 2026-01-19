# CopyImageDMDGL - Method Documentation

> Tài liệu này tổng hợp logic và phương thức chính của ứng dụng CopyImageDMDGL (C#) để hỗ trợ việc chuyển đổi sang Go.

---

## 📋 Mục đích ứng dụng

Ứng dụng console **sao chép hình ảnh hàng loạt** từ thư mục network share nguồn sang thư mục đích.

**Use case thực tế:**
- Copy hình mẫu sản phẩm từ server lưu trữ tạm (`HÌNH CHƯA TẢI`) sang server chính (`HinhAnh`)
- Hỗ trợ ghi đè file nếu cần
- Xử lý song song để tăng tốc độ copy

---

## ⚙️ Cấu hình (Constants)

```
SOURCE_PATH      = "\\192.1.1.1\DM_DON_GIA_LUONG\ROUTING 2023 + HÌNH MẪU\HINHMAUSP\HÌNH CHƯA TẢI"
DESTINATION_PATH = "\\192.1.1.20\dmdgl$\HinhAnh"
```

**Gợi ý cho Go:**
- Sử dụng file config (JSON, YAML, TOML) hoặc environment variables
- Hỗ trợ command-line flags: `--source`, `--dest`, `--overwrite`

---

## 🔄 Luồng xử lý chính (Main Flow)

```
START
  │
  ▼
┌─────────────────────────────────────┐
│ 1. Hiển thị menu lựa chọn           │
│    - 0: Không copy (thoát)          │
│    - 1: Copy và ghi đè              │
└─────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────┐
│ 2. Validate input                   │
│    - Lặp cho đến khi nhập đúng 0/1  │
└─────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────┐
│ 3. Kiểm tra thư mục nguồn tồn tại   │
│    - Nếu không tồn tại → thông báo  │
└─────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────┐
│ 4. Lấy danh sách files              │
│    - Nếu rỗng → thông báo           │
└─────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────┐
│ 5. Copy song song (Parallel)        │
│    - Với mỗi file:                  │
│      • Tạo đường dẫn đích           │
│      • Copy file (ghi đè nếu chọn)  │
│      • Xử lý exception              │
│      • Log kết quả                  │
└─────────────────────────────────────┘
  │
  ▼
┌─────────────────────────────────────┐
│ 6. Hoàn thành                       │
│    - Hiển thị thông báo             │
│    - Đợi người dùng nhấn phím       │
└─────────────────────────────────────┘
  │
  ▼
END
```

---

## 📦 Các Method chính

### 1. `Main()` - Entry Point

**Mục đích:** Điều phối toàn bộ luồng xử lý

**Logic:**
```
1. In menu console
2. Đọc input người dùng → validate (chỉ chấp nhận 0 hoặc 1)
3. Nếu chọn 0 → thoát
4. Nếu chọn 1:
   a. Kiểm tra thư mục nguồn
   b. Lấy danh sách files
   c. Copy song song với option ghi đè = true
5. In kết quả và đợi phím bấm
```

**Input:** Không có tham số
**Output:** Console output

---

### 2. `IsFileLocked(filePath string) bool`

**Mục đích:** Kiểm tra xem file có đang bị lock (đang được mở bởi process khác) không

**Logic:**
```
1. Thử mở file với mode ReadWrite và FileShare.None
2. Nếu mở được → file không bị lock → return false
3. Nếu IOException → file đang bị lock → return true
```

**Input:** `filePath` - Đường dẫn tuyệt đối đến file cần kiểm tra
**Output:** `bool` - `true` nếu file đang bị lock, `false` nếu không

**Code C# gốc:**
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

**Gợi ý cho Go:**
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

**Mục đích:** Copy một file từ nguồn sang đích

**Logic:**
```
1. Lấy tên file từ đường dẫn nguồn
2. Tạo đường dẫn đích = destPath + fileName
3. Copy file:
   - Nếu overwrite = true → ghi đè nếu tồn tại
   - Nếu overwrite = false → bỏ qua nếu tồn tại
4. Xử lý exception:
   - File đang bị lock → log và bỏ qua
   - Lỗi khác → log lỗi
```

**Input:**
- `sourcePath` - Đường dẫn file nguồn
- `destPath` - Thư mục đích
- `overwrite` - Có ghi đè không

**Output:** `error` hoặc `nil`

---

### 4. `CopyFilesParallel(files []string, destPath string, overwrite bool)`

**Mục đích:** Copy nhiều files song song để tăng hiệu suất

**Logic (C# dùng Parallel.ForEach):**
```
1. Với mỗi file trong danh sách (song song):
   a. Gọi CopyFile()
   b. Log kết quả: ✓ thành công hoặc ✗ thất bại
```

**Gợi ý cho Go (dùng goroutines + WaitGroup):**
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

## 🛡️ Xử lý lỗi (Error Handling)

| Loại lỗi | Xử lý |
|----------|-------|
| Thư mục nguồn không tồn tại | Log thông báo và thoát |
| Không có file nào trong thư mục | Log thông báo và thoát |
| File đang bị lock | Bỏ qua, log với prefix ✗ |
| IOException khác | Log chi tiết lỗi với prefix ✗ |
| Exception chung | Log và tiếp tục với file khác |

---

## 🚀 Gợi ý cải tiến cho Go

### 1. **CLI với Cobra/Flag**
```
copyimage --source "/path/to/source" --dest "/path/to/dest" --overwrite --workers 10
```

### 2. **Progress Bar**
Sử dụng thư viện như `github.com/schollz/progressbar/v3`

### 3. **Logging có cấu trúc**
Dùng `log/slog` (Go 1.21+) hoặc `zerolog`/`zap`

### 4. **Retry mechanism**
Thử lại khi copy thất bại (tối đa 3 lần)

### 5. **Dry-run mode**
Option `--dry-run` để xem trước file sẽ được copy

### 6. **Filter files**
Option `--ext .jpg,.png` để chỉ copy một số loại file

### 7. **Worker Pool**
Kiểm soát số lượng goroutines đồng thời để tránh quá tải

### 8. **Report/Summary**
```
========== KẾT QUẢ ==========
Tổng số files: 100
Thành công:    95
Thất bại:      3
Bỏ qua:        2
Thời gian:     5.2s
=============================
```

### 9. **Config file**
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

## 📁 Cấu trúc project Go đề xuất

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

## ✅ Checklist chuyển đổi

- [ ] Tạo project Go mới với `go mod init`
- [ ] Implement `config` package (load từ file/flags/env)
- [ ] Implement `isFileLocked()` function
- [ ] Implement `copyFile()` function
- [ ] Implement `copyFilesParallel()` với worker pool
- [ ] Thêm CLI flags (cobra hoặc flag package)
- [ ] Thêm progress bar
- [ ] Thêm summary report
- [ ] Viết unit tests
- [ ] Build và test trên Windows với UNC paths

---

*Tài liệu được tạo: 2026-01-19*
