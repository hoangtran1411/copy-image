package copier

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"copy-image/internal/config"
	"copy-image/internal/utils"

	"github.com/schollz/progressbar/v3"
)

// CopyResult represents the result of a copy operation
type CopyResult struct {
	FileName string
	Success  bool
	Skipped  bool
	Error    error
}

// CopySummary represents the summary of all copy operations
type CopySummary struct {
	TotalFiles  int
	Successful  int
	Failed      int
	Skipped     int
	Duration    time.Duration
	FailedFiles []string
}

// Copier handles file copying operations
type Copier struct {
	config  *config.Config
	results []CopyResult
}

// New creates a new Copier instance
func New(cfg *config.Config) *Copier {
	return &Copier{
		config:  cfg,
		results: make([]CopyResult, 0),
	}
}

// GetFiles retrieves all files from the source directory
func (c *Copier) GetFiles() ([]string, error) {
	if !utils.DirExists(c.config.Source) {
		return nil, fmt.Errorf("source directory does not exist: %s", c.config.Source)
	}

	var files []string
	entries, err := os.ReadDir(c.config.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fileName := entry.Name()
		ext := strings.ToLower(filepath.Ext(fileName))

		// Check extension filter
		if c.config.HasExtensionFilter() && !c.config.IsExtensionAllowed(ext) {
			continue
		}

		files = append(files, filepath.Join(c.config.Source, fileName))
	}

	return files, nil
}

// CopyFile copies a single file from source to destination
func (c *Copier) CopyFile(sourcePath string, overwrite bool) error {
	fileName := filepath.Base(sourcePath)
	destPath := filepath.Join(c.config.Destination, fileName)

	// Check if destination file exists
	if utils.FileExists(destPath) && !overwrite {
		return nil // Skip if file exists and overwrite is false
	}

	// Check if source file is locked
	if utils.IsFileLocked(sourcePath) {
		return fmt.Errorf("file is locked by another process")
	}

	// Ensure destination directory exists
	if err := utils.EnsureDir(c.config.Destination); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Open source file
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy content
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Sync to ensure data is written to disk
	if err := dstFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

// CopyFileWithRetry copies a file with retry mechanism
func (c *Copier) CopyFileWithRetry(sourcePath string) CopyResult {
	fileName := filepath.Base(sourcePath)
	destPath := filepath.Join(c.config.Destination, fileName)

	// Check if should skip
	if utils.FileExists(destPath) && !c.config.Overwrite {
		return CopyResult{
			FileName: fileName,
			Success:  false,
			Skipped:  true,
			Error:    nil,
		}
	}

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		err := c.CopyFile(sourcePath, c.config.Overwrite)
		if err == nil {
			return CopyResult{
				FileName: fileName,
				Success:  true,
				Skipped:  false,
				Error:    nil,
			}
		}
		lastErr = err

		// Wait before retry (exponential backoff)
		if attempt < c.config.MaxRetries {
			time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
		}
	}

	return CopyResult{
		FileName: fileName,
		Success:  false,
		Skipped:  false,
		Error:    lastErr,
	}
}

// CopyFilesParallel copies multiple files in parallel using worker pool
func (c *Copier) CopyFilesParallel(files []string) CopySummary {
	startTime := time.Now()

	var (
		successful int32
		failed     int32
		skipped    int32
		wg         sync.WaitGroup
		failedMu   sync.Mutex
	)

	failedFiles := make([]string, 0)
	semaphore := make(chan struct{}, c.config.Workers)

	// Create progress bar
	bar := progressbar.NewOptions(len(files),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetDescription("[cyan]Copying files...[reset]"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire
			defer func() { <-semaphore }() // Release

			if c.config.DryRun {
				fmt.Printf("  [DRY-RUN] Would copy: %s\n", filepath.Base(f))
				atomic.AddInt32(&successful, 1)
			} else {
				result := c.CopyFileWithRetry(f)

				if result.Success {
					atomic.AddInt32(&successful, 1)
				} else if result.Skipped {
					atomic.AddInt32(&skipped, 1)
				} else {
					atomic.AddInt32(&failed, 1)
					failedMu.Lock()
					failedFiles = append(failedFiles, fmt.Sprintf("%s: %v", result.FileName, result.Error))
					failedMu.Unlock()
				}
			}

			_ = bar.Add(1)
		}(file)
	}

	wg.Wait()
	_ = bar.Finish()
	fmt.Println() // New line after progress bar

	return CopySummary{
		TotalFiles:  len(files),
		Successful:  int(successful),
		Failed:      int(failed),
		Skipped:     int(skipped),
		Duration:    time.Since(startTime),
		FailedFiles: failedFiles,
	}
}

// PrintSummary prints a formatted summary of the copy operation
func (s *CopySummary) PrintSummary() {
	fmt.Println("\n========== KẾT QUẢ ==========")
	fmt.Printf("Tổng số files: %d\n", s.TotalFiles)
	fmt.Printf("Thành công:    %d ✓\n", s.Successful)
	fmt.Printf("Thất bại:      %d ✗\n", s.Failed)
	fmt.Printf("Bỏ qua:        %d ⊘\n", s.Skipped)
	fmt.Printf("Thời gian:     %.2fs\n", s.Duration.Seconds())
	fmt.Println("==============================")

	if len(s.FailedFiles) > 0 {
		fmt.Println("\n===== FILES THẤT BẠI =====")
		for _, f := range s.FailedFiles {
			fmt.Printf("  ✗ %s\n", f)
		}
		fmt.Println("==========================")
	}
}
