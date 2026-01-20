---
description: Standard procedure for professional commits and releases using Git Release Management Skill
---

Use this workflow to ensure a clean project history and well-documented releases.

### Step 1: Daily Commits
When you complete a small change (bug fix, feature, docs), use the following commands:
```bash
git add .
git commit -m "<type>(<scope>): <description>"
```
*Tip:* You can ask AI: "**Commit recent changes following Conventional Commits standards**".

### Step 2: Prepare Release (When code is stable)
1. Verify that all features are complete and tested.
2. Update the version in the code (e.g., `CurrentVersion` in `updater.go`).
3. Commit the version bump:
   ```bash
   git add .
   git commit -m "chore: bump version to vX.Y.Z"
   ```

### Step 3: Create Professional Release Tag
Use the `git tag` command with a comprehensive message:
```bash
git tag -a vX.Y.Z -m "vX.Y.Z - [Release Title]

🚀 [New Features]
- ...
🛠️ [Fixes & Improvements]
- ...
"
```
*Tip:* You can ask AI: "**Create release tag v2.2.0 for changes since v2.1.3, formatted beautifully according to the Skill**".

### Step 4: Push to GitHub
```bash
git push origin main --tags
```

---
**Note:** You can ask AI to perform **Step 1** separately at any time. Only proceed to **Steps 3-4** when you are truly ready to release a new version.
