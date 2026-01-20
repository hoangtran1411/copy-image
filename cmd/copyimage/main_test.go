package main

import (
	"os"
	"path/filepath"
	"testing"

	"copy-image/internal/config"
)

func TestParseExtensions(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "single extension with dot",
			input:    ".jpg",
			expected: []string{".jpg"},
		},
		{
			name:     "single extension without dot",
			input:    "jpg",
			expected: []string{".jpg"},
		},
		{
			name:     "multiple extensions",
			input:    ".jpg,.png,.gif",
			expected: []string{".jpg", ".png", ".gif"},
		},
		{
			name:     "extensions with spaces",
			input:    ".jpg, .png, .gif",
			expected: []string{".jpg", ".png", ".gif"},
		},
		{
			name:     "mixed with and without dots",
			input:    "jpg,.png,gif",
			expected: []string{".jpg", ".png", ".gif"},
		},
		{
			name:     "uppercase extensions",
			input:    ".JPG,.PNG",
			expected: []string{".jpg", ".png"},
		},
		{
			name:     "extra commas",
			input:    ".jpg,,,.png",
			expected: []string{".jpg", ".png"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseExtensions(tt.input)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d extensions, got %d", len(tt.expected), len(result))
				return
			}

			for i, ext := range result {
				if ext != tt.expected[i] {
					t.Errorf("Expected extension[%d]=%s, got %s", i, tt.expected[i], ext)
				}
			}
		})
	}
}

func TestPrintBanner(t *testing.T) {
	// Just ensure it doesn't panic
	printBanner()
}

func TestVersion(t *testing.T) {
	if version == "" {
		t.Error("Expected version to be set")
	}
	if version != "1.0.0" {
		t.Errorf("Expected version='1.0.0', got %s", version)
	}
}

func TestLoadConfigDefault(t *testing.T) {
	cfg := loadConfig("", "", "", false, 10, false, "")

	if cfg == nil {
		t.Fatal("Expected non-nil config")
	}
}

func TestLoadConfigWithCLIOverrides(t *testing.T) {
	cfg := loadConfig("", "/src/path", "/dst/path", true, 15, true, ".jpg,.png")

	if cfg.Source != "/src/path" {
		t.Errorf("Expected Source='/src/path', got %s", cfg.Source)
	}
	if cfg.Destination != "/dst/path" {
		t.Errorf("Expected Destination='/dst/path', got %s", cfg.Destination)
	}
	if cfg.Overwrite != true {
		t.Error("Expected Overwrite=true")
	}
	if cfg.Workers != 15 {
		t.Errorf("Expected Workers=15, got %d", cfg.Workers)
	}
	if cfg.DryRun != true {
		t.Error("Expected DryRun=true")
	}
	if len(cfg.Extensions) != 2 {
		t.Errorf("Expected 2 extensions, got %d", len(cfg.Extensions))
	}
}

func TestLoadConfigFromFile(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `
source: "/file/source"
destination: "/file/dest"
workers: 8
overwrite: false
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg := loadConfig(configPath, "", "", false, 10, false, "")

	if cfg.Source != "/file/source" {
		t.Errorf("Expected Source='/file/source', got %s", cfg.Source)
	}
	if cfg.Destination != "/file/dest" {
		t.Errorf("Expected Destination='/file/dest', got %s", cfg.Destination)
	}
	if cfg.Workers != 8 {
		t.Errorf("Expected Workers=8, got %d", cfg.Workers)
	}
}

func TestLoadConfigFileWithCLIOverride(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `
source: "/file/source"
destination: "/file/dest"
workers: 8
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// CLI should override file config
	cfg := loadConfig(configPath, "/cli/source", "", false, 10, false, "")

	if cfg.Source != "/cli/source" {
		t.Errorf("Expected CLI Source='/cli/source' to override, got %s", cfg.Source)
	}
	if cfg.Destination != "/file/dest" {
		t.Errorf("Expected Destination from file='/file/dest', got %s", cfg.Destination)
	}
}

func TestLoadConfigNonExistentFile(t *testing.T) {
	cfg := loadConfig("/non/existent/config.yaml", "/src", "/dst", false, 10, false, "")

	// Should return default config with CLI overrides
	if cfg.Source != "/src" {
		t.Errorf("Expected Source='/src', got %s", cfg.Source)
	}
}

func TestPrintConfig(t *testing.T) {
	cfg := &config.Config{
		Source:      "/test/source",
		Destination: "/test/dest",
		Workers:     10,
		Overwrite:   true,
		DryRun:      false,
		Extensions:  []string{},
	}

	// Just ensure it doesn't panic
	printConfig(cfg)
}

func TestPrintConfigWithExtensions(t *testing.T) {
	cfg := &config.Config{
		Source:      "/test/source",
		Destination: "/test/dest",
		Workers:     10,
		Overwrite:   true,
		DryRun:      false,
		Extensions:  []string{".jpg", ".png"},
	}

	// Just ensure it doesn't panic with extensions
	printConfig(cfg)
}

func TestLoadConfigWorkersNotChanged(t *testing.T) {
	// When workers is default (10), it should not override config file value
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `
source: "/file/source"
destination: "/file/dest"
workers: 5
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// workers=10 is default, should NOT override
	cfg := loadConfig(configPath, "", "", false, 10, false, "")

	if cfg.Workers != 5 {
		t.Errorf("Expected Workers=5 from file (not overridden), got %d", cfg.Workers)
	}
}

func TestLoadConfigWorkersChanged(t *testing.T) {
	// When workers is NOT default, it should override
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `
source: "/file/source"
destination: "/file/dest"
workers: 5
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// workers=20 is NOT default, should override
	cfg := loadConfig(configPath, "", "", false, 20, false, "")

	if cfg.Workers != 20 {
		t.Errorf("Expected Workers=20 (CLI override), got %d", cfg.Workers)
	}
}

func TestLoadConfigEmptyConfigFile(t *testing.T) {
	// Empty config file name should use defaults
	cfg := loadConfig("", "", "", false, 10, false, "")

	if cfg == nil {
		t.Fatal("Expected non-nil config")
	}
	// Should have default values
	if cfg.Workers != 10 {
		t.Errorf("Expected default Workers=10, got %d", cfg.Workers)
	}
}

func TestLoadConfigAllCLIFlags(t *testing.T) {
	cfg := loadConfig("", "/source", "/dest", true, 25, true, ".jpg,.png,.gif")

	if cfg.Source != "/source" {
		t.Errorf("Expected Source='/source', got %s", cfg.Source)
	}
	if cfg.Destination != "/dest" {
		t.Errorf("Expected Destination='/dest', got %s", cfg.Destination)
	}
	if !cfg.Overwrite {
		t.Error("Expected Overwrite=true")
	}
	if cfg.Workers != 25 {
		t.Errorf("Expected Workers=25, got %d", cfg.Workers)
	}
	if !cfg.DryRun {
		t.Error("Expected DryRun=true")
	}
	if len(cfg.Extensions) != 3 {
		t.Errorf("Expected 3 extensions, got %d", len(cfg.Extensions))
	}
}

func TestLoadConfigPartialCLIFlags(t *testing.T) {
	// Only source and dest provided
	cfg := loadConfig("", "/partial/source", "/partial/dest", false, 10, false, "")

	if cfg.Source != "/partial/source" {
		t.Errorf("Expected Source='/partial/source', got %s", cfg.Source)
	}
	if cfg.Destination != "/partial/dest" {
		t.Errorf("Expected Destination='/partial/dest', got %s", cfg.Destination)
	}
	// Other flags should remain default
	if cfg.Overwrite {
		t.Error("Expected Overwrite=false")
	}
	if cfg.DryRun {
		t.Error("Expected DryRun=false")
	}
}

func TestLoadConfigOverwriteFalseNoOverride(t *testing.T) {
	// When overwrite CLI flag is false, it should not override true in config
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `
source: "/file/source"
destination: "/file/dest"
overwrite: true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg := loadConfig(configPath, "", "", false, 10, false, "")

	// File says overwrite=true, CLI says false, but false is default so file wins
	if cfg.Overwrite != true {
		t.Errorf("Expected Overwrite=true from file, got %v", cfg.Overwrite)
	}
}

func TestLoadConfigDryRunFalseNoOverride(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `
source: "/file/source"
destination: "/file/dest"
dry_run: true
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg := loadConfig(configPath, "", "", false, 10, false, "")

	// File says dry_run=true, CLI says false (default), so file wins
	if cfg.DryRun != true {
		t.Errorf("Expected DryRun=true from file, got %v", cfg.DryRun)
	}
}

func TestLoadConfigExtensionsFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `
source: "/file/source"
destination: "/file/dest"
extensions:
  - .webp
  - .avif
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg := loadConfig(configPath, "", "", false, 10, false, "")

	if len(cfg.Extensions) != 2 {
		t.Errorf("Expected 2 extensions from file, got %d", len(cfg.Extensions))
	}
}

func TestLoadConfigExtensionsCLIOverridesFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `
source: "/file/source"
destination: "/file/dest"
extensions:
  - .webp
  - .avif
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// CLI extensions should override file extensions
	cfg := loadConfig(configPath, "", "", false, 10, false, ".bmp")

	if len(cfg.Extensions) != 1 {
		t.Errorf("Expected 1 extension from CLI, got %d", len(cfg.Extensions))
	}
	if cfg.Extensions[0] != ".bmp" {
		t.Errorf("Expected extension '.bmp', got '%s'", cfg.Extensions[0])
	}
}
