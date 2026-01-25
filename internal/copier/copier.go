package copier

import (
	"context"
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

// CopyResult represents the result of a single file copy operation.
// It tracks whether the copy succeeded, was skipped, or failed with an error.
type CopyResult struct {
	FileName string
	Success  bool
	Skipped  bool
	Error    error
}

// CopySummary represents the aggregate results of a batch copy operation.
// It provides statistics for reporting progress to users.
type CopySummary struct {
	TotalFiles  int
	Successful  int
	Failed      int
	Skipped     int
	Duration    time.Duration
	FailedFiles []string
}

// ProgressCallback is a function type for reporting copy progress.
// It receives the current count, total count, current filename, and status.
type ProgressCallback func(current int, total int, fileName string, status string)

// Copier handles file copying operations with support for parallel execution,
// retry logic, and progress reporting.
type Copier struct {
	config  *config.Config
	results []CopyResult
}

// New creates a new Copier instance with the given configuration.
// The copier is stateless between copy operations, so the same instance
// can be reused for multiple copy batches.
func New(cfg *config.Config) *Copier {
	return &Copier{
		config:  cfg,
		results: make([]CopyResult, 0),
	}
}

// GetFiles retrieves all files from the source directory that match
// the extension filter (if configured). Only regular files are returned;
// directories are not included.
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

		// Skip files that don't match the extension filter
		if c.config.HasExtensionFilter() && !c.config.IsExtensionAllowed(ext) {
			continue
		}

		files = append(files, filepath.Join(c.config.Source, fileName))
	}

	return files, nil
}

// CopyFile copies a single file from source to the configured destination.
// If overwrite is false and the destination file exists, the copy is skipped.
// The function ensures the destination directory exists before copying.
func (c *Copier) CopyFile(ctx context.Context, sourcePath string, overwrite bool) error {
	// Check for cancellation before starting
	if err := ctx.Err(); err != nil {
		return err
	}

	fileName := filepath.Base(sourcePath)
	destPath := filepath.Join(c.config.Destination, fileName)

	// Skip if file exists and we're not overwriting
	if utils.FileExists(destPath) && !overwrite {
		return nil
	}

	// Check if source file is locked by another process
	if utils.IsFileLocked(sourcePath) {
		return fmt.Errorf("file is locked by another process")
	}

	// Ensure destination directory exists
	if err := utils.EnsureDir(c.config.Destination); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Open source file for reading
	srcFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() { _ = srcFile.Close() }()

	// Create destination file
	dstFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		// Capture close errors - they may indicate write failures
		if cerr := dstFile.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("failed to close destination file: %w", cerr)
		}
	}()

	// Copy content using buffered I/O
	// Only CopyBuffer allows cancellation if we implement a custom reader,
	// but standard Copy respects context if passed to a wrapper, or we just check before.
	// For now, we stick to io.Copy but at least we checked context at start.
	// A more advanced version would use a cancelable reader.
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	// Sync to ensure data is flushed to disk
	// This is important for data integrity, especially on network drives
	if err := dstFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync file: %w", err)
	}

	return nil
}

// CopyFileWithRetry attempts to copy a file with automatic retries on failure.
// It uses exponential backoff between retries to handle transient errors
// like network hiccups or temporary file locks.
func (c *Copier) CopyFileWithRetry(ctx context.Context, sourcePath string) CopyResult {
	fileName := filepath.Base(sourcePath)
	destPath := filepath.Join(c.config.Destination, fileName)

	// Check if we should skip this file
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
		// Check context before each attempt
		if err := ctx.Err(); err != nil {
			return CopyResult{
				FileName: fileName,
				Success:  false,
				Skipped:  false,
				Error:    err,
			}
		}

		err := c.CopyFile(ctx, sourcePath, c.config.Overwrite)
		if err == nil {
			return CopyResult{
				FileName: fileName,
				Success:  true,
				Skipped:  false,
				Error:    nil,
			}
		}
		lastErr = err

		// Exponential backoff
		if attempt < c.config.MaxRetries {
			select {
			case <-ctx.Done():
				return CopyResult{
					FileName: fileName,
					Success:  false,
					Skipped:  false,
					Error:    ctx.Err(),
				}
			case <-time.After(time.Duration(attempt+1) * 100 * time.Millisecond):
				// Continue to next attempt
			}
		}
	}

	return CopyResult{
		FileName: fileName,
		Success:  false,
		Skipped:  false,
		Error:    lastErr,
	}
}

// CopyFilesParallel copies multiple files concurrently using a worker pool.
// This version is for CLI mode - it uses a terminal progress bar.
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

	// Create terminal progress bar for CLI mode
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
			semaphore <- struct{}{}        // Acquire worker slot
			defer func() { <-semaphore }() // Release worker slot

			if c.config.DryRun {
				fmt.Printf("  [DRY-RUN] Would copy: %s\n", filepath.Base(f))
				atomic.AddInt32(&successful, 1)
			} else {
				// CLI mode doesn't have a cancellation context yet, using Background
				result := c.CopyFileWithRetry(context.Background(), f)

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

// CopyFilesParallelWithEvents copies files concurrently with progress callbacks.
// This version is designed for GUI mode (Wails) - instead of printing to terminal,
// it calls the provided callback function to report progress.
//
// The context parameter allows cancellation of the operation. When cancelled,
// in-progress copies will complete but no new copies will start.
func (c *Copier) CopyFilesParallelWithEvents(ctx context.Context, files []string, onProgress ProgressCallback) CopySummary {
	startTime := time.Now()

	var (
		successful int32
		failed     int32
		skipped    int32
		processed  int32
		wg         sync.WaitGroup
		failedMu   sync.Mutex
	)

	failedFiles := make([]string, 0)
	semaphore := make(chan struct{}, c.config.Workers)
	total := len(files)

	for _, file := range files {
		// Check for cancellation before starting new work
		select {
		case <-ctx.Done():
			// Context cancelled - stop processing new files
			break
		default:
			// Continue processing
		}

		wg.Add(1)
		go func(f string) {
			defer wg.Done()

			// Acquire worker slot (or wait for one to become available)
			select {
			case semaphore <- struct{}{}:
				defer func() { <-semaphore }()
			case <-ctx.Done():
				// Cancelled while waiting for a worker slot
				return
			}

			fileName := filepath.Base(f)
			var status string

			if c.config.DryRun {
				status = "success"
				atomic.AddInt32(&successful, 1)
			} else {
				result := c.CopyFileWithRetry(ctx, f)

				if result.Success {
					status = "success"
					atomic.AddInt32(&successful, 1)
				} else if result.Skipped {
					status = "skipped"
					atomic.AddInt32(&skipped, 1)
				} else {
					status = "failed"
					atomic.AddInt32(&failed, 1)
					failedMu.Lock()
					failedFiles = append(failedFiles, fmt.Sprintf("%s: %v", result.FileName, result.Error))
					failedMu.Unlock()
				}
			}

			// Report progress via callback
			current := int(atomic.AddInt32(&processed, 1))
			if onProgress != nil {
				onProgress(current, total, fileName, status)
			}
		}(file)
	}

	wg.Wait()

	return CopySummary{
		TotalFiles:  total,
		Successful:  int(successful),
		Failed:      int(failed),
		Skipped:     int(skipped),
		Duration:    time.Since(startTime),
		FailedFiles: failedFiles,
	}
}

// PrintSummary prints a formatted summary of the copy operation to stdout.
// This is used in CLI mode to display results after a batch copy completes.
func (s *CopySummary) PrintSummary() {
	fmt.Println("\n========== RESULTS ==========")
	fmt.Printf("Total files: %d\n", s.TotalFiles)
	fmt.Printf("Successful:  %d ✓\n", s.Successful)
	fmt.Printf("Failed:      %d ✗\n", s.Failed)
	fmt.Printf("Skipped:     %d ⊘\n", s.Skipped)
	fmt.Printf("Duration:    %.2fs\n", s.Duration.Seconds())
	fmt.Println("==============================")

	if len(s.FailedFiles) > 0 {
		fmt.Println("\n===== FAILED FILES =====")
		for _, f := range s.FailedFiles {
			fmt.Printf("  ✗ %s\n", f)
		}
		fmt.Println("========================")
	}
}
