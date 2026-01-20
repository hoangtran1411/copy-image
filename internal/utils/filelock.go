package utils

import (
	"os"
)

// IsFileLocked checks if a file is currently locked by another process
// Returns true if the file is locked for reading, false otherwise.
// Note: We only check for read access because we only need to read the file to copy it.
// Checking O_RDWR (Read/Write) causes "locked" errors if the file is Read-Only or
// if the user doesn't have Write permissions (common on network shares).
func IsFileLocked(filePath string) bool {
	// Try to open for READ ONLY.
	// If we can read it, we can copy it.
	file, err := os.Open(filePath)
	if err != nil {
		// Only consider it locked if we can't even read it.
		// Detailed error checking could distinguish "locked" vs "permission denied",
		// but for now, if we can't read it, we can't copy it anyway.
		return true
	}
	_ = file.Close()
	return false
}

// FileExists checks if a file exists at the given path
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks if a directory exists at the given path
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(path string) error {
	if !DirExists(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
