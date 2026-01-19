package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsFileLocked(t *testing.T) {
	// Create a temp file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// File should not be locked
	if IsFileLocked(testFile) {
		t.Error("Expected file to not be locked")
	}
}

func TestIsFileLockedNonExistent(t *testing.T) {
	// Non-existent file should return true (like locked)
	if !IsFileLocked("/non/existent/file.txt") {
		t.Error("Expected true for non-existent file")
	}
}

func TestFileExists(t *testing.T) {
	// Create a temp file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "exists.txt")

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// File exists
	if !FileExists(testFile) {
		t.Error("Expected FileExists to return true for existing file")
	}

	// File doesn't exist
	if FileExists(filepath.Join(tmpDir, "nonexistent.txt")) {
		t.Error("Expected FileExists to return false for non-existent file")
	}
}

func TestDirExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Directory exists
	if !DirExists(tmpDir) {
		t.Error("Expected DirExists to return true for existing directory")
	}

	// Directory doesn't exist
	if DirExists(filepath.Join(tmpDir, "nonexistent")) {
		t.Error("Expected DirExists to return false for non-existent directory")
	}

	// File path (not directory)
	testFile := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if DirExists(testFile) {
		t.Error("Expected DirExists to return false for a file")
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create nested directory
	newDir := filepath.Join(tmpDir, "level1", "level2", "level3")
	if err := EnsureDir(newDir); err != nil {
		t.Errorf("EnsureDir failed: %v", err)
	}

	// Verify directory was created
	if !DirExists(newDir) {
		t.Error("Expected directory to be created")
	}

	// Call again on existing directory (should not fail)
	if err := EnsureDir(newDir); err != nil {
		t.Errorf("EnsureDir on existing directory failed: %v", err)
	}
}

func TestEnsureDirOnExisting(t *testing.T) {
	tmpDir := t.TempDir()

	// Call EnsureDir on existing directory
	if err := EnsureDir(tmpDir); err != nil {
		t.Errorf("EnsureDir on existing directory failed: %v", err)
	}
}
