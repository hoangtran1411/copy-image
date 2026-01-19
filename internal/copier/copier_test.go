package copier

import (
	"os"
	"path/filepath"
	"testing"

	"copy-image/internal/config"
)

func TestNew(t *testing.T) {
	cfg := config.DefaultConfig()
	c := New(cfg)

	if c == nil {
		t.Fatal("Expected Copier instance, got nil")
	}
	if c.config != cfg {
		t.Error("Expected config to be set correctly")
	}
}

func TestCopyFile(t *testing.T) {
	// Create temp directories
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create a test file
	testContent := []byte("Hello, World!")
	srcFile := filepath.Join(srcDir, "test.txt")
	if err := os.WriteFile(srcFile, testContent, 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   true,
		MaxRetries:  1,
	}

	c := New(cfg)

	// Test copy
	err := c.CopyFile(srcFile, true)
	if err != nil {
		t.Errorf("CopyFile failed: %v", err)
	}

	// Verify file was copied
	dstFile := filepath.Join(dstDir, "test.txt")
	content, err := os.ReadFile(dstFile)
	if err != nil {
		t.Errorf("Failed to read destination file: %v", err)
	}

	if string(content) != string(testContent) {
		t.Errorf("Expected content %q, got %q", testContent, content)
	}
}

func TestCopyFileNoOverwrite(t *testing.T) {
	// Create temp directories
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create source file
	srcContent := []byte("Source content")
	srcFile := filepath.Join(srcDir, "test.txt")
	if err := os.WriteFile(srcFile, srcContent, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create existing destination file with different content
	dstContent := []byte("Existing content")
	dstFile := filepath.Join(dstDir, "test.txt")
	if err := os.WriteFile(dstFile, dstContent, 0644); err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   false,
		MaxRetries:  1,
	}

	c := New(cfg)

	// Test copy without overwrite
	err := c.CopyFile(srcFile, false)
	if err != nil {
		t.Errorf("CopyFile failed: %v", err)
	}

	// Verify original content was preserved
	content, err := os.ReadFile(dstFile)
	if err != nil {
		t.Errorf("Failed to read destination file: %v", err)
	}

	if string(content) != string(dstContent) {
		t.Errorf("Expected original content %q, got %q", dstContent, content)
	}
}

func TestGetFiles(t *testing.T) {
	// Create temp directory with test files
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create test files
	testFiles := []string{"image1.jpg", "image2.png", "document.pdf"}
	for _, f := range testFiles {
		if err := os.WriteFile(filepath.Join(srcDir, f), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   true,
		Extensions:  []string{},
	}

	c := New(cfg)

	files, err := c.GetFiles()
	if err != nil {
		t.Errorf("GetFiles failed: %v", err)
	}

	if len(files) != 3 {
		t.Errorf("Expected 3 files, got %d", len(files))
	}
}

func TestGetFilesWithExtensionFilter(t *testing.T) {
	// Create temp directory with test files
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create test files
	testFiles := []string{"image1.jpg", "image2.png", "document.pdf", "photo.jpeg"}
	for _, f := range testFiles {
		if err := os.WriteFile(filepath.Join(srcDir, f), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   true,
		Extensions:  []string{".jpg", ".jpeg"},
	}

	c := New(cfg)

	files, err := c.GetFiles()
	if err != nil {
		t.Errorf("GetFiles failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 files (jpg, jpeg only), got %d", len(files))
	}
}

func TestCopyFilesParallel(t *testing.T) {
	// Create temp directories
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", "file3.txt"}
	var filePaths []string
	for _, f := range testFiles {
		path := filepath.Join(srcDir, f)
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		filePaths = append(filePaths, path)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     2,
		Overwrite:   true,
		MaxRetries:  1,
	}

	c := New(cfg)

	summary := c.CopyFilesParallel(filePaths)

	if summary.TotalFiles != 3 {
		t.Errorf("Expected TotalFiles=3, got %d", summary.TotalFiles)
	}

	if summary.Successful != 3 {
		t.Errorf("Expected Successful=3, got %d", summary.Successful)
	}

	if summary.Failed != 0 {
		t.Errorf("Expected Failed=0, got %d", summary.Failed)
	}

	// Verify all files were copied
	for _, f := range testFiles {
		dstPath := filepath.Join(dstDir, f)
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			t.Errorf("File %s was not copied", f)
		}
	}
}

func TestGetFilesNonExistentDir(t *testing.T) {
	cfg := &config.Config{
		Source:      "/non/existent/directory",
		Destination: "/some/dest",
		Workers:     1,
	}

	c := New(cfg)

	_, err := c.GetFiles()
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}

func TestCopyFileWithRetrySuccess(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create test file
	srcFile := filepath.Join(srcDir, "retry_test.txt")
	if err := os.WriteFile(srcFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   true,
		MaxRetries:  3,
	}

	c := New(cfg)

	result := c.CopyFileWithRetry(srcFile)

	if !result.Success {
		t.Error("Expected Success=true")
	}
	if result.Skipped {
		t.Error("Expected Skipped=false")
	}
	if result.Error != nil {
		t.Errorf("Expected no error, got: %v", result.Error)
	}
}

func TestCopyFileWithRetrySkipped(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create source file
	srcFile := filepath.Join(srcDir, "skip_test.txt")
	if err := os.WriteFile(srcFile, []byte("source"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create existing destination file
	dstFile := filepath.Join(dstDir, "skip_test.txt")
	if err := os.WriteFile(dstFile, []byte("existing"), 0644); err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   false, // Don't overwrite
		MaxRetries:  1,
	}

	c := New(cfg)

	result := c.CopyFileWithRetry(srcFile)

	if result.Success {
		t.Error("Expected Success=false for skipped file")
	}
	if !result.Skipped {
		t.Error("Expected Skipped=true")
	}
}

func TestCopyFilesParallelWithSkip(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create source files
	srcFile1 := filepath.Join(srcDir, "new.txt")
	srcFile2 := filepath.Join(srcDir, "existing.txt")
	if err := os.WriteFile(srcFile1, []byte("new content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	if err := os.WriteFile(srcFile2, []byte("source content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create existing destination file
	dstFile2 := filepath.Join(dstDir, "existing.txt")
	if err := os.WriteFile(dstFile2, []byte("existing content"), 0644); err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     2,
		Overwrite:   false, // Don't overwrite
		MaxRetries:  1,
	}

	c := New(cfg)

	summary := c.CopyFilesParallel([]string{srcFile1, srcFile2})

	if summary.TotalFiles != 2 {
		t.Errorf("Expected TotalFiles=2, got %d", summary.TotalFiles)
	}
	if summary.Successful != 1 {
		t.Errorf("Expected Successful=1, got %d", summary.Successful)
	}
	if summary.Skipped != 1 {
		t.Errorf("Expected Skipped=1, got %d", summary.Skipped)
	}
}

func TestCopyFilesParallelDryRun(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create source files
	files := []string{"dry1.txt", "dry2.txt"}
	var filePaths []string
	for _, f := range files {
		path := filepath.Join(srcDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		filePaths = append(filePaths, path)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     2,
		Overwrite:   true,
		DryRun:      true, // Dry run mode
		MaxRetries:  1,
	}

	c := New(cfg)

	summary := c.CopyFilesParallel(filePaths)

	if summary.TotalFiles != 2 {
		t.Errorf("Expected TotalFiles=2, got %d", summary.TotalFiles)
	}
	if summary.Successful != 2 {
		t.Errorf("Expected Successful=2 in dry-run, got %d", summary.Successful)
	}

	// Verify files were NOT actually copied in dry-run mode
	for _, f := range files {
		dstPath := filepath.Join(dstDir, f)
		if _, err := os.Stat(dstPath); !os.IsNotExist(err) {
			t.Errorf("File %s should NOT exist in dry-run mode", f)
		}
	}
}

func TestCopyFileSourceNotFound(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   true,
		MaxRetries:  1,
	}

	c := New(cfg)

	// Try to copy non-existent file
	err := c.CopyFile(filepath.Join(srcDir, "nonexistent.txt"), true)
	if err == nil {
		t.Error("Expected error for non-existent source file")
	}
}

func TestCopyFileWithRetryFailed(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   true,
		MaxRetries:  1,
	}

	c := New(cfg)

	// Try to copy non-existent file
	result := c.CopyFileWithRetry(filepath.Join(srcDir, "nonexistent.txt"))

	if result.Success {
		t.Error("Expected Success=false for failed copy")
	}
	if result.Skipped {
		t.Error("Expected Skipped=false for failed copy")
	}
	if result.Error == nil {
		t.Error("Expected error for failed copy")
	}
}

func TestCopySummaryPrintSummary(t *testing.T) {
	// Test with no failures
	summary := &CopySummary{
		TotalFiles:  100,
		Successful:  95,
		Failed:      3,
		Skipped:     2,
		Duration:    5 * 1000000000, // 5 seconds in nanoseconds
		FailedFiles: []string{"file1.txt: error1", "file2.txt: error2"},
	}

	// Just call PrintSummary to ensure it doesn't panic
	// We can't easily test console output, but we verify it runs without error
	summary.PrintSummary()
}

func TestCopySummaryPrintSummaryNoFailures(t *testing.T) {
	summary := &CopySummary{
		TotalFiles:  10,
		Successful:  10,
		Failed:      0,
		Skipped:     0,
		Duration:    1 * 1000000000,
		FailedFiles: []string{},
	}

	// Should run without panic
	summary.PrintSummary()
}

func TestGetFilesIgnoresDirectories(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create a file
	if err := os.WriteFile(filepath.Join(srcDir, "file.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Create a subdirectory
	if err := os.Mkdir(filepath.Join(srcDir, "subdir"), 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Extensions:  []string{},
	}

	c := New(cfg)

	files, err := c.GetFiles()
	if err != nil {
		t.Errorf("GetFiles failed: %v", err)
	}

	// Should only get the file, not the directory
	if len(files) != 1 {
		t.Errorf("Expected 1 file (ignoring directory), got %d", len(files))
	}
}

func TestCopyFileOverwriteExisting(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create source file with new content
	srcContent := []byte("NEW CONTENT")
	srcFile := filepath.Join(srcDir, "overwrite.txt")
	if err := os.WriteFile(srcFile, srcContent, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create existing destination file with old content
	dstFile := filepath.Join(dstDir, "overwrite.txt")
	if err := os.WriteFile(dstFile, []byte("OLD CONTENT"), 0644); err != nil {
		t.Fatalf("Failed to create destination file: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   true,
		MaxRetries:  1,
	}

	c := New(cfg)

	err := c.CopyFile(srcFile, true)
	if err != nil {
		t.Errorf("CopyFile failed: %v", err)
	}

	// Verify content was overwritten
	content, err := os.ReadFile(dstFile)
	if err != nil {
		t.Errorf("Failed to read destination file: %v", err)
	}

	if string(content) != string(srcContent) {
		t.Errorf("Expected content %q, got %q", srcContent, content)
	}
}

func TestCopyFilesParallelWithMultipleWorkers(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create many test files
	numFiles := 20
	var filePaths []string
	for i := 0; i < numFiles; i++ {
		fileName := filepath.Join(srcDir, "file"+string(rune('A'+i))+".txt")
		if err := os.WriteFile(fileName, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		filePaths = append(filePaths, fileName)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     5, // Multiple workers
		Overwrite:   true,
		MaxRetries:  1,
	}

	c := New(cfg)

	summary := c.CopyFilesParallel(filePaths)

	if summary.TotalFiles != numFiles {
		t.Errorf("Expected TotalFiles=%d, got %d", numFiles, summary.TotalFiles)
	}
	if summary.Successful != numFiles {
		t.Errorf("Expected Successful=%d, got %d", numFiles, summary.Successful)
	}
	if summary.Failed != 0 {
		t.Errorf("Expected Failed=0, got %d", summary.Failed)
	}
}

func TestCopyResultFields(t *testing.T) {
	result := CopyResult{
		FileName: "test.txt",
		Success:  true,
		Skipped:  false,
		Error:    nil,
	}

	if result.FileName != "test.txt" {
		t.Errorf("Expected FileName='test.txt', got %s", result.FileName)
	}
	if !result.Success {
		t.Error("Expected Success=true")
	}
	if result.Skipped {
		t.Error("Expected Skipped=false")
	}
}

func TestCopySummaryFields(t *testing.T) {
	summary := CopySummary{
		TotalFiles:  50,
		Successful:  45,
		Failed:      3,
		Skipped:     2,
		Duration:    2 * 1000000000,
		FailedFiles: []string{"a.txt", "b.txt"},
	}

	if summary.TotalFiles != 50 {
		t.Errorf("Expected TotalFiles=50, got %d", summary.TotalFiles)
	}
	if summary.Successful != 45 {
		t.Errorf("Expected Successful=45, got %d", summary.Successful)
	}
	if summary.Failed != 3 {
		t.Errorf("Expected Failed=3, got %d", summary.Failed)
	}
	if summary.Skipped != 2 {
		t.Errorf("Expected Skipped=2, got %d", summary.Skipped)
	}
	if len(summary.FailedFiles) != 2 {
		t.Errorf("Expected 2 failed files, got %d", len(summary.FailedFiles))
	}
}

func TestCopyFilesParallelEmptyList(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     2,
		Overwrite:   true,
		MaxRetries:  1,
	}

	c := New(cfg)

	// Empty file list
	summary := c.CopyFilesParallel([]string{})

	if summary.TotalFiles != 0 {
		t.Errorf("Expected TotalFiles=0, got %d", summary.TotalFiles)
	}
	if summary.Successful != 0 {
		t.Errorf("Expected Successful=0, got %d", summary.Successful)
	}
	if summary.Failed != 0 {
		t.Errorf("Expected Failed=0, got %d", summary.Failed)
	}
}

func TestCopyFileToNonExistentDestDir(t *testing.T) {
	srcDir := t.TempDir()
	// Destination is a nested directory that doesn't exist yet
	dstDir := filepath.Join(t.TempDir(), "nested", "deep", "dir")

	srcFile := filepath.Join(srcDir, "test.txt")
	if err := os.WriteFile(srcFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   true,
		MaxRetries:  1,
	}

	c := New(cfg)

	// Should create destination directory and copy
	err := c.CopyFile(srcFile, true)
	if err != nil {
		t.Errorf("CopyFile failed: %v", err)
	}

	// Verify file was copied
	dstFile := filepath.Join(dstDir, "test.txt")
	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Error("File was not copied to the new directory")
	}
}

func TestGetFilesEmptyDirectory(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Extensions:  []string{},
	}

	c := New(cfg)

	files, err := c.GetFiles()
	if err != nil {
		t.Errorf("GetFiles failed: %v", err)
	}

	if len(files) != 0 {
		t.Errorf("Expected 0 files from empty directory, got %d", len(files))
	}
}

func TestCopyFileWithRetryMultipleAttempts(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create a test file
	srcFile := filepath.Join(srcDir, "multiretry.txt")
	if err := os.WriteFile(srcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   true,
		MaxRetries:  3, // Multiple retries
	}

	c := New(cfg)

	result := c.CopyFileWithRetry(srcFile)

	if !result.Success {
		t.Error("Expected Success=true")
	}
	if result.FileName != "multiretry.txt" {
		t.Errorf("Expected FileName='multiretry.txt', got %s", result.FileName)
	}
}

func TestCopierWithZeroRetries(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	srcFile := filepath.Join(srcDir, "zero_retry.txt")
	if err := os.WriteFile(srcFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   true,
		MaxRetries:  0, // No retries
	}

	c := New(cfg)

	result := c.CopyFileWithRetry(srcFile)

	if !result.Success {
		t.Error("Expected Success=true even with 0 retries")
	}
}

func TestGetFilesOnlyFiltered(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create files with various extensions
	files := map[string]string{
		"photo.jpg":    "jpg",
		"document.pdf": "pdf",
		"data.xlsx":    "xlsx",
		"image.gif":    "gif",
	}

	for name := range files {
		if err := os.WriteFile(filepath.Join(srcDir, name), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Extensions:  []string{".gif"}, // Only .gif
	}

	c := New(cfg)

	result, err := c.GetFiles()
	if err != nil {
		t.Errorf("GetFiles failed: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 file (.gif only), got %d", len(result))
	}
}

func TestCopyFilesParallelWithFailed(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create one real file
	realFile := filepath.Join(srcDir, "real.txt")
	if err := os.WriteFile(realFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     2,
		Overwrite:   true,
		MaxRetries:  0,
	}

	c := New(cfg)

	// Include one real file and one non-existent file
	fakeFile := filepath.Join(srcDir, "nonexistent.txt")
	summary := c.CopyFilesParallel([]string{realFile, fakeFile})

	if summary.TotalFiles != 2 {
		t.Errorf("Expected TotalFiles=2, got %d", summary.TotalFiles)
	}
	if summary.Successful != 1 {
		t.Errorf("Expected Successful=1, got %d", summary.Successful)
	}
	if summary.Failed != 1 {
		t.Errorf("Expected Failed=1, got %d", summary.Failed)
	}
	if len(summary.FailedFiles) != 1 {
		t.Errorf("Expected 1 failed file, got %d", len(summary.FailedFiles))
	}
}

func TestCopyFileLargeContent(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create a larger file (1MB)
	largeContent := make([]byte, 1024*1024)
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	srcFile := filepath.Join(srcDir, "large.bin")
	if err := os.WriteFile(srcFile, largeContent, 0644); err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Overwrite:   true,
		MaxRetries:  1,
	}

	c := New(cfg)

	err := c.CopyFile(srcFile, true)
	if err != nil {
		t.Errorf("CopyFile failed: %v", err)
	}

	// Verify content
	dstFile := filepath.Join(dstDir, "large.bin")
	content, err := os.ReadFile(dstFile)
	if err != nil {
		t.Errorf("Failed to read destination file: %v", err)
	}

	if len(content) != len(largeContent) {
		t.Errorf("Expected %d bytes, got %d", len(largeContent), len(content))
	}
}

func TestCopySummaryDuration(t *testing.T) {
	summary := CopySummary{
		TotalFiles:  10,
		Successful:  10,
		Failed:      0,
		Skipped:     0,
		Duration:    5500000000, // 5.5 seconds in nanoseconds
		FailedFiles: []string{},
	}

	// Test Duration.Seconds() calculation
	seconds := summary.Duration.Seconds()
	if seconds < 5.4 || seconds > 5.6 {
		t.Errorf("Expected Duration ~5.5s, got %.2fs", seconds)
	}
}

func TestGetFilesWithMixedExtensions(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := t.TempDir()

	// Create files with various extensions including uppercase
	testFiles := []string{
		"photo.JPG",      // uppercase
		"image.jpg",      // lowercase
		"document.PDF",   // should be excluded
		"photo2.JPEG",    // uppercase
		"picture.jpeg",   // lowercase
	}

	for _, f := range testFiles {
		if err := os.WriteFile(filepath.Join(srcDir, f), []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	cfg := &config.Config{
		Source:      srcDir,
		Destination: dstDir,
		Workers:     1,
		Extensions:  []string{".jpg", ".jpeg"},
	}

	c := New(cfg)

	files, err := c.GetFiles()
	if err != nil {
		t.Errorf("GetFiles failed: %v", err)
	}

	// Should get all .jpg and .jpeg files (4 total)
	if len(files) != 4 {
		t.Errorf("Expected 4 files (.jpg and .jpeg), got %d", len(files))
	}
}


