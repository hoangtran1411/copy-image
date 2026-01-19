package utils

import (
	"os"
)

// IsFileLocked checks if a file is currently locked by another process
// Returns true if the file is locked or cannot be accessed, false otherwise
func IsFileLocked(filePath string) bool {
	file, err := os.OpenFile(filePath, os.O_RDWR, 0666)
	if err != nil {
		return true // File is locked or doesn't exist
	}
	defer func() { _ = file.Close() }()
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
