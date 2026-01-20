//go:build windows

package main

import (
	"context"
	"fmt"

	"copy-image/internal/config"
	"copy-image/internal/copier"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct represents the main application.
// It holds the application context and manages the lifecycle of copy operations.
// The context is used for Wails runtime calls like dialogs and events.
type App struct {
	ctx    context.Context
	config *config.Config
	copier *copier.Copier

	// cancelFunc allows us to cancel ongoing copy operations.
	// This is essential for providing a responsive UI where users can stop
	// long-running tasks without waiting for completion.
	cancelFunc context.CancelFunc
}

// NewApp creates a new App application struct.
// We initialize with nil values because the actual setup happens in startup()
// after Wails has fully initialized the runtime context.
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods later.
// This is the first lifecycle hook where we have access to Wails runtime.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.config = config.DefaultConfig()

	// Attempt to load config from file on startup.
	// We silently ignore errors here because the app should still work
	// with default config if no config file exists.
	if loadedCfg, err := config.LoadFromFile("config.yaml"); err == nil {
		a.config = loadedCfg
	}
}

// GetConfig returns the current configuration.
// The frontend uses this to populate the settings UI on load.
func (a *App) GetConfig() *config.Config {
	return a.config
}

// UpdateConfig updates the application configuration.
// This is called when the user changes settings in the UI.
// We validate before accepting to prevent invalid states.
func (a *App) UpdateConfig(cfg *config.Config) error {
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}
	a.config = cfg
	return nil
}

// SaveConfig persists the current configuration to a YAML file.
// This ensures user preferences survive app restarts.
func (a *App) SaveConfig() error {
	return a.config.SaveToFile("config.yaml")
}

// SelectSourceFolder opens a native directory picker dialog for source folder.
// Using native dialogs provides a familiar experience and respects OS accessibility settings.
func (a *App) SelectSourceFolder() (string, error) {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Source Folder",
	})
	if err != nil {
		return "", fmt.Errorf("failed to open directory dialog: %w", err)
	}
	return folder, nil
}

// SelectDestFolder opens a native directory picker dialog for destination folder.
func (a *App) SelectDestFolder() (string, error) {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Destination Folder",
	})
	if err != nil {
		return "", fmt.Errorf("failed to open directory dialog: %w", err)
	}
	return folder, nil
}

// ScanFiles scans the source directory and returns a list of files to copy.
// This is separated from the copy operation so the UI can show a preview
// of how many files will be copied before the user commits.
func (a *App) ScanFiles() ([]string, error) {
	if a.config.Source == "" {
		return nil, fmt.Errorf("source path is not configured")
	}

	a.copier = copier.New(a.config)
	files, err := a.copier.GetFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to scan files: %w", err)
	}

	return files, nil
}

// ProgressEvent represents a single progress update sent to the frontend.
// We use a struct instead of multiple parameters to make the event payload
// self-documenting and easier to extend in the future.
type ProgressEvent struct {
	Current  int     `json:"current"`
	Total    int     `json:"total"`
	Percent  float64 `json:"percent"`
	FileName string  `json:"fileName"`
	Status   string  `json:"status"` // "copying", "success", "failed", "skipped"
}

// CopyResult represents the final result of a copy operation.
// This provides a summary for the UI to display completion statistics.
type CopyResult struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	TotalFiles  int      `json:"totalFiles"`
	Successful  int      `json:"successful"`
	Failed      int      `json:"failed"`
	Skipped     int      `json:"skipped"`
	FailedFiles []string `json:"failedFiles"`
	Duration    float64  `json:"duration"` // in seconds
}

// StartCopy begins the file copy operation.
// It creates a cancellable context so users can stop the operation mid-way.
// Progress updates are emitted as events to keep the UI responsive.
func (a *App) StartCopy(overwrite bool) CopyResult {
	if a.copier == nil {
		return CopyResult{
			Success: false,
			Message: "Please scan files first",
		}
	}

	// Update the overwrite setting based on user choice
	a.config.Overwrite = overwrite

	// Create a cancellable context for this copy operation.
	// This allows users to stop long-running copies without closing the app.
	ctx, cancel := context.WithCancel(a.ctx)
	a.cancelFunc = cancel
	defer func() {
		a.cancelFunc = nil
	}()

	// Get files to copy
	files, err := a.copier.GetFiles()
	if err != nil {
		return CopyResult{
			Success: false,
			Message: fmt.Sprintf("Failed to get files: %v", err),
		}
	}

	if len(files) == 0 {
		return CopyResult{
			Success: true,
			Message: "No files found to copy",
		}
	}

	// Emit initial progress
	runtime.EventsEmit(a.ctx, "copy:start", map[string]any{
		"total": len(files),
	})

	// Create a new copier with event emitting capability
	summary := a.copier.CopyFilesParallelWithEvents(ctx, files, func(current int, total int, fileName string, status string) {
		// Emit progress event to frontend
		runtime.EventsEmit(a.ctx, "copy:progress", ProgressEvent{
			Current:  current,
			Total:    total,
			Percent:  float64(current) / float64(total) * 100,
			FileName: fileName,
			Status:   status,
		})
	})

	// Build result
	result := CopyResult{
		Success:     summary.Failed == 0,
		TotalFiles:  summary.TotalFiles,
		Successful:  summary.Successful,
		Failed:      summary.Failed,
		Skipped:     summary.Skipped,
		FailedFiles: summary.FailedFiles,
		Duration:    summary.Duration.Seconds(),
	}

	if summary.Failed > 0 {
		result.Message = fmt.Sprintf("Completed with %d errors", summary.Failed)
	} else {
		result.Message = fmt.Sprintf("Successfully copied %d files", summary.Successful)
	}

	// Emit completion event
	runtime.EventsEmit(a.ctx, "copy:complete", result)

	return result
}

// CancelCopy stops an ongoing copy operation.
// This is called when the user clicks the cancel button.
// The cancellation is graceful - in-progress file copies may complete,
// but no new files will start copying.
func (a *App) CancelCopy() {
	if a.cancelFunc != nil {
		a.cancelFunc()
		runtime.EventsEmit(a.ctx, "copy:cancelled", nil)
	}
}
