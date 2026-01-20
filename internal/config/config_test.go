package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Source != "" {
		t.Errorf("Expected empty Source, got %s", cfg.Source)
	}
	if cfg.Destination != "" {
		t.Errorf("Expected empty Destination, got %s", cfg.Destination)
	}
	if cfg.Workers != 10 {
		t.Errorf("Expected Workers=10, got %d", cfg.Workers)
	}
	if cfg.Overwrite != false {
		t.Error("Expected Overwrite=false")
	}
	if cfg.MaxRetries != 3 {
		t.Errorf("Expected MaxRetries=3, got %d", cfg.MaxRetries)
	}
	if cfg.DryRun != false {
		t.Error("Expected DryRun=false")
	}
	if len(cfg.Extensions) != 0 {
		t.Errorf("Expected empty Extensions, got %v", cfg.Extensions)
	}
}

func TestLoadFromFile(t *testing.T) {
	// Create temp config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test-config.yaml")

	configContent := `
source: "/test/source"
destination: "/test/dest"
workers: 5
overwrite: true
max_retries: 2
dry_run: true
extensions:
  - .jpg
  - .png
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	cfg, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	if cfg.Source != "/test/source" {
		t.Errorf("Expected Source='/test/source', got %s", cfg.Source)
	}
	if cfg.Destination != "/test/dest" {
		t.Errorf("Expected Destination='/test/dest', got %s", cfg.Destination)
	}
	if cfg.Workers != 5 {
		t.Errorf("Expected Workers=5, got %d", cfg.Workers)
	}
	if cfg.Overwrite != true {
		t.Error("Expected Overwrite=true")
	}
	if cfg.MaxRetries != 2 {
		t.Errorf("Expected MaxRetries=2, got %d", cfg.MaxRetries)
	}
	if cfg.DryRun != true {
		t.Error("Expected DryRun=true")
	}
	if len(cfg.Extensions) != 2 {
		t.Errorf("Expected 2 extensions, got %d", len(cfg.Extensions))
	}
}

func TestLoadFromFileNotFound(t *testing.T) {
	_, err := LoadFromFile("/non/existent/config.yaml")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestLoadFromFileInvalidYaml(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "invalid.yaml")

	// Write invalid YAML
	invalidContent := `
source: [invalid
destination: unclosed bracket
`
	if err := os.WriteFile(configPath, []byte(invalidContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	_, err := LoadFromFile(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &Config{
				Source:      "/path/to/source",
				Destination: "/path/to/dest",
				Workers:     10,
			},
			expectError: false,
		},
		{
			name: "missing source",
			config: &Config{
				Source:      "",
				Destination: "/path/to/dest",
			},
			expectError: true,
		},
		{
			name: "missing destination",
			config: &Config{
				Source:      "/path/to/source",
				Destination: "",
			},
			expectError: true,
		},
		{
			name: "workers too low - auto fix",
			config: &Config{
				Source:      "/path/to/source",
				Destination: "/path/to/dest",
				Workers:     0,
			},
			expectError: false,
		},
		{
			name: "workers too high - auto fix",
			config: &Config{
				Source:      "/path/to/source",
				Destination: "/path/to/dest",
				Workers:     100,
			},
			expectError: false,
		},
		{
			name: "negative retries - auto fix",
			config: &Config{
				Source:      "/path/to/source",
				Destination: "/path/to/dest",
				Workers:     10,
				MaxRetries:  -1,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateWorkersAutoFix(t *testing.T) {
	cfg := &Config{
		Source:      "/path/to/source",
		Destination: "/path/to/dest",
		Workers:     0,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfg.Workers != 1 {
		t.Errorf("Expected Workers to be fixed to 1, got %d", cfg.Workers)
	}

	cfg.Workers = 100
	if err := cfg.Validate(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfg.Workers != 50 {
		t.Errorf("Expected Workers to be fixed to 50, got %d", cfg.Workers)
	}
}

func TestValidateMaxRetriesAutoFix(t *testing.T) {
	cfg := &Config{
		Source:      "/path/to/source",
		Destination: "/path/to/dest",
		Workers:     10,
		MaxRetries:  -5,
	}
	if err := cfg.Validate(); err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if cfg.MaxRetries != 0 {
		t.Errorf("Expected MaxRetries to be fixed to 0, got %d", cfg.MaxRetries)
	}
}

func TestHasExtensionFilter(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.HasExtensionFilter() {
		t.Error("Expected HasExtensionFilter=false for empty extensions")
	}

	cfg.Extensions = []string{".jpg", ".png"}
	if !cfg.HasExtensionFilter() {
		t.Error("Expected HasExtensionFilter=true for non-empty extensions")
	}
}

func TestIsExtensionAllowed(t *testing.T) {
	cfg := &Config{
		Extensions: []string{".jpg", ".png", ".gif"},
	}

	tests := []struct {
		ext      string
		expected bool
	}{
		{".jpg", true},
		{".png", true},
		{".gif", true},
		{".pdf", false},
		{".doc", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			result := cfg.IsExtensionAllowed(tt.ext)
			if result != tt.expected {
				t.Errorf("IsExtensionAllowed(%s) = %v, expected %v", tt.ext, result, tt.expected)
			}
		})
	}
}

func TestIsExtensionAllowedNoFilter(t *testing.T) {
	cfg := &Config{
		Extensions: []string{},
	}

	// All extensions should be allowed when no filter is set
	if !cfg.IsExtensionAllowed(".anything") {
		t.Error("Expected all extensions to be allowed when no filter")
	}
	if !cfg.IsExtensionAllowed(".random") {
		t.Error("Expected all extensions to be allowed when no filter")
	}
}

// TestIsExtensionAllowedCaseInsensitive verifies that extension matching is case-insensitive.
func TestIsExtensionAllowedCaseInsensitive(t *testing.T) {
	cfg := &Config{
		Extensions: []string{".jpg", ".PNG"},
	}

	tests := []struct {
		ext      string
		expected bool
	}{
		{".jpg", true},
		{".JPG", true},
		{".Jpg", true},
		{".png", true},
		{".PNG", true},
		{".gif", false},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			result := cfg.IsExtensionAllowed(tt.ext)
			if result != tt.expected {
				t.Errorf("IsExtensionAllowed(%s) = %v, expected %v", tt.ext, result, tt.expected)
			}
		})
	}
}

// TestSaveToFile verifies that configuration can be persisted and reloaded.
func TestSaveToFile(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "saved-config.yaml")

	// Create config with custom values
	cfg := &Config{
		Source:      "/test/source",
		Destination: "/test/dest",
		Workers:     15,
		Overwrite:   true,
		Extensions:  []string{".jpg", ".png"},
		MaxRetries:  5,
		DryRun:      true,
	}

	// Save to file
	if err := cfg.SaveToFile(configPath); err != nil {
		t.Fatalf("SaveToFile failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Reload and verify
	loaded, err := LoadFromFile(configPath)
	if err != nil {
		t.Fatalf("LoadFromFile failed: %v", err)
	}

	if loaded.Source != cfg.Source {
		t.Errorf("Source mismatch: got %s, want %s", loaded.Source, cfg.Source)
	}
	if loaded.Workers != cfg.Workers {
		t.Errorf("Workers mismatch: got %d, want %d", loaded.Workers, cfg.Workers)
	}
	if loaded.DryRun != cfg.DryRun {
		t.Errorf("DryRun mismatch: got %v, want %v", loaded.DryRun, cfg.DryRun)
	}
}

// TestSaveToFileInvalidPath verifies error handling for invalid file paths.
func TestSaveToFileInvalidPath(t *testing.T) {
	cfg := DefaultConfig()

	// Try to save to an invalid path
	err := cfg.SaveToFile("/nonexistent/directory/config.yaml")
	if err == nil {
		t.Error("Expected error for invalid path")
	}
}

// TestCopyGroups tests the Copy Groups related methods.
func TestCopyGroups(t *testing.T) {
	cfg := DefaultConfig()

	// Initially no groups
	if len(cfg.Groups) != 0 {
		t.Errorf("Expected 0 groups, got %d", len(cfg.Groups))
	}

	// Add a group
	group := CopyGroup{
		ID:      "group-1",
		Name:    "Test Group",
		Source:  "/source",
		Enabled: true,
		Destinations: []Destination{
			{ID: "dest-1", Path: "/dest1", Overwrite: true, Enabled: true},
			{ID: "dest-2", Path: "/dest2", Overwrite: false, Enabled: false},
		},
	}
	cfg.AddGroup(group)

	if len(cfg.Groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(cfg.Groups))
	}

	// Find group
	found := cfg.FindGroup("group-1")
	if found == nil {
		t.Fatal("FindGroup returned nil")
	}
	if found.Name != "Test Group" {
		t.Errorf("Expected name 'Test Group', got %s", found.Name)
	}

	// Find non-existent group
	notFound := cfg.FindGroup("nonexistent")
	if notFound != nil {
		t.Error("Expected nil for non-existent group")
	}

	// Get enabled groups
	enabled := cfg.GetEnabledGroups()
	if len(enabled) != 1 {
		t.Errorf("Expected 1 enabled group, got %d", len(enabled))
	}

	// Add disabled group
	disabledGroup := CopyGroup{
		ID:      "group-2",
		Name:    "Disabled Group",
		Enabled: false,
	}
	cfg.AddGroup(disabledGroup)

	enabled = cfg.GetEnabledGroups()
	if len(enabled) != 1 {
		t.Errorf("Expected 1 enabled group after adding disabled, got %d", len(enabled))
	}

	// Remove group
	removed := cfg.RemoveGroup("group-1")
	if !removed {
		t.Error("RemoveGroup returned false")
	}
	if len(cfg.Groups) != 1 {
		t.Errorf("Expected 1 group after removal, got %d", len(cfg.Groups))
	}

	// Try to remove non-existent group
	removed = cfg.RemoveGroup("nonexistent")
	if removed {
		t.Error("RemoveGroup should return false for non-existent group")
	}
}

// TestValidateWithGroups verifies validation works correctly with Copy Groups.
func TestValidateWithGroups(t *testing.T) {
	// Config with groups should not require source/destination
	cfg := &Config{
		Source:      "", // Empty
		Destination: "", // Empty
		Workers:     10,
		Groups: []CopyGroup{
			{ID: "group-1", Source: "/source", Enabled: true},
		},
	}

	err := cfg.Validate()
	if err != nil {
		t.Errorf("Config with groups should be valid, got error: %v", err)
	}
}

// TestDestinationStruct tests the Destination struct fields.
func TestDestinationStruct(t *testing.T) {
	dest := Destination{
		ID:        "dest-1",
		Path:      "/path/to/dest",
		Overwrite: true,
		Enabled:   true,
	}

	if dest.ID != "dest-1" {
		t.Errorf("Expected ID 'dest-1', got %s", dest.ID)
	}
	if dest.Path != "/path/to/dest" {
		t.Errorf("Expected Path '/path/to/dest', got %s", dest.Path)
	}
	if !dest.Overwrite {
		t.Error("Expected Overwrite to be true")
	}
	if !dest.Enabled {
		t.Error("Expected Enabled to be true")
	}
}
