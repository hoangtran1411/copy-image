---
description: Cách thực hiện commit và release chuẩn chuyên nghiệp sử dụng Git Release Management Skill
---

Sử dụng workflow này để đảm bảo lịch sử dự án luôn sạch đẹp và các bản phát hành có đầy đủ thông tin.

### Bước 1: Commit công việc hàng ngày
Khi bạn hoàn thành một thay đổi nhỏ (fix bug, thêm feature, sửa docs), hãy dùng lệnh:
```bash
git add .
git commit -m "<type>(<scope>): <description>"
```
*Gợi ý:* Bạn có thể bảo AI: "**Hãy commit các thay đổi vừa rồi theo chuẩn Conventional Commits**".

### Bước 2: Chuẩn bị Release (Khi code đã ổn định)
1. Kiểm tra lại các tính năng đã hoàn thiện.
2. Cập nhật version trong code (ví dụ: `CurrentVersion` trong `updater.go`).
3. Commit việc nâng version:
   ```bash
   git add .
   git commit -m "chore: bump version to vX.Y.Z"
   ```

### Bước 3: Tạo Tag Release chuyên nghiệp
Sử dụng lệnh `git tag` với thông điệp đầy đủ:
```bash
git tag -a vX.Y.Z -m "vX.Y.Z - [Tiêu đề Release]

🚀 [Tính năng mới]
- ...
🛠️ [Sửa lỗi & Cải tiến]
- ...
"
```
*Gợi ý:* Bạn có thể bảo AI: "**Hãy tạo release tag v2.2.0 cho những gì chúng ta đã làm từ bản v2.1.3 đến nay, trình bày đẹp mắt theo Skill đã có**".

### Bước 4: Đẩy lên GitHub
```bash
git push origin main --tags
```

---
**Lưu ý:** Bạn có thể yêu cầu AI thực hiện riêng lẻ Bước 1 bất cứ lúc nào, và chỉ thực hiện Bước 3-4 khi bạn thực sự muốn ra mắt phiên bản mới.
