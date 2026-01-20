//go:build windows

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// CurrentVersion holds the application version.
// This should be updated when releasing new versions.
// For production builds, use ldflags to inject the version at build time:
// go build -ldflags "-X main.CurrentVersion=v2.1.3"
var CurrentVersion = "v2.1.4"

// GitHubOwner and GitHubRepo identify the repository for update checks.
// These constants define where to look for new releases on GitHub.
const (
	GitHubOwner = "hoangtran1411"
	GitHubRepo  = "copy-image"
)

// UpdateInfo holds information about available updates.
// This struct is returned to the frontend to display update notifications.
type UpdateInfo struct {
	Available   bool   `json:"available"`
	CurrentVer  string `json:"currentVersion"`
	LatestVer   string `json:"latestVersion"`
	DownloadURL string `json:"downloadUrl"`
	ReleaseURL  string `json:"releaseUrl"`
}

// GitHubRelease represents the relevant fields from GitHub's release API response.
// We only parse the fields we need to minimize processing overhead.
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// GetCurrentVersion returns the current app version.
// The frontend displays this in the header to help users identify their version.
func (a *App) GetCurrentVersion() string {
	return CurrentVersion
}

// CheckForUpdate queries GitHub API to check if a newer version is available.
// This runs asynchronously on app startup so it doesn't block the UI.
// Returns update info including download URL if an update is available.
func (a *App) CheckForUpdate() UpdateInfo {
	info := UpdateInfo{
		Available:  false,
		CurrentVer: CurrentVersion,
	}

	// Construct the GitHub API URL for the latest release.
	// Using the releases/latest endpoint gives us the most recent non-prerelease version.
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", GitHubOwner, GitHubRepo)

	resp, err := http.Get(url)
	if err != nil {
		// Network errors are silently ignored - the app should work offline.
		return info
	}
	defer resp.Body.Close()

	// Non-200 responses indicate API issues or rate limiting.
	// We fail gracefully by returning no update available.
	if resp.StatusCode != http.StatusOK {
		return info
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return info
	}

	info.LatestVer = release.TagName
	info.ReleaseURL = release.HTMLURL

	// Find the Windows executable in the release assets.
	// We specifically look for the "desktop-windows-amd64" version to avoid
	// accidentally downloading the CLI version within the Desktop app.
	for _, asset := range release.Assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, "desktop") && strings.HasSuffix(name, ".exe") {
			info.DownloadURL = asset.BrowserDownloadURL
			break
		}
	}

	// Fallback for older naming conventions or if 'desktop' is not found
	if info.DownloadURL == "" {
		for _, asset := range release.Assets {
			if strings.HasSuffix(strings.ToLower(asset.Name), ".exe") {
				info.DownloadURL = asset.BrowserDownloadURL
				break
			}
		}
	}

	// Compare versions using semantic versioning.
	// Only mark as available if the remote version is strictly newer.
	if info.LatestVer != "" && CompareVersions(info.LatestVer, CurrentVersion) {
		info.Available = true
	}

	return info
}

// CompareVersions determines if v1 is newer than v2 using semantic versioning.
// Returns true if v1 > v2, false otherwise.
// This handles version strings like "v1.2.3" or "1.2.3".
func CompareVersions(v1, v2 string) bool {
	// Remove the 'v' prefix if present for consistent parsing.
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	parts1 := parseVersion(v1)
	parts2 := parseVersion(v2)

	// Compare major, minor, patch in order of significance.
	// Return as soon as we find a difference.
	for i := 0; i < 3; i++ {
		if parts1[i] > parts2[i] {
			return true
		}
		if parts1[i] < parts2[i] {
			return false
		}
	}

	// Versions are equal
	return false
}

// parseVersion splits a version string into [major, minor, patch] integers.
// Missing parts default to 0 (e.g., "1.2" becomes [1, 2, 0]).
func parseVersion(v string) [3]int {
	var result [3]int
	parts := strings.Split(v, ".")

	for i := 0; i < len(parts) && i < 3; i++ {
		// Use Sscanf for safe integer parsing - invalid inputs become 0.
		fmt.Sscanf(parts[i], "%d", &result[i])
	}

	return result
}

// PerformUpdate downloads and installs a new version of the application.
// This is a complex operation that:
// 1. Downloads the new executable to a temp file
// 2. Creates a batch script to replace the running executable
// 3. Exits the current app and lets the batch script do the swap
//
// We use a batch script because Windows locks running executables,
// so we can't directly overwrite the file while it's running.
func (a *App) PerformUpdate(downloadURL string) (bool, error) {
	if downloadURL == "" {
		return false, fmt.Errorf("no download URL provided")
	}

	// Get the path to the currently running executable.
	// This is the file we'll replace with the new version.
	exePath, err := os.Executable()
	if err != nil {
		return false, fmt.Errorf("failed to get executable path: %w", err)
	}
	exePath, _ = filepath.Abs(exePath)

	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "copyimage_update.exe")

	// Notify the frontend that download is starting.
	runtime.EventsEmit(a.ctx, "update:progress", "Downloading update...")

	// Download the new version.
	resp, err := http.Get(downloadURL)
	if err != nil {
		return false, fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create the temp file for the download.
	out, err := os.Create(tempFile)
	if err != nil {
		return false, fmt.Errorf("failed to create temp file: %w", err)
	}

	// Copy the downloaded content to the temp file.
	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return false, fmt.Errorf("failed to save update: %w", err)
	}

	runtime.EventsEmit(a.ctx, "update:progress", "Installing update...")

	// Create a batch script that will:
	// 1. Wait for this process to exit (timeout)
	// 2. Delete the old executable
	// 3. Move the new executable to the original location
	// 4. Start the new executable
	// 5. Delete itself
	//
	// This approach is necessary on Windows because you can't replace
	// a running executable directly.
	batchPath := filepath.Join(tempDir, "update_copyimage.bat")
	// Optimized batch script for Windows:
	// - timeout: waits for the app to close
	// - del: removes old exe
	// - move: installs new exe
	// - start: launches the updated app
	// - (goto) trick: safely deletes the script itself after execution
	batchContent := fmt.Sprintf(`@echo off
timeout /t 2 /nobreak >nul
del /f /q "%s"
move /y "%s" "%s"
start "" "%s"
(goto) 2>nul & del "%%~f0"
`, exePath, tempFile, exePath, exePath)

	// Write batch script with standard permissions for Windows (0666)
	if err := os.WriteFile(batchPath, []byte(batchContent), 0666); err != nil {
		return false, fmt.Errorf("failed to create update script: %w", err)
	}

	// Run the batch script as a detached process to ensure it continues
	// running after the main application exits.
	cmd := exec.Command("cmd", "/c", batchPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x00000008, // DETACHED_PROCESS
	}
	
	if err := cmd.Start(); err != nil {
		return false, fmt.Errorf("failed to start update script: %w", err)
	}

	// Exit the application to allow the batch script to replace the executable.
	runtime.Quit(a.ctx)

	return true, nil
}
