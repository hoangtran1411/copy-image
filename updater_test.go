package main

import (
	"testing"
)

// TestCompareVersions verifies that semantic version comparison works correctly.
// This is critical for the auto-update feature to properly determine if an update is available.
func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected bool
	}{
		// v1 > v2 cases
		{
			name:     "major version higher",
			v1:       "v2.0.0",
			v2:       "v1.0.0",
			expected: true,
		},
		{
			name:     "minor version higher",
			v1:       "v1.2.0",
			v2:       "v1.1.0",
			expected: true,
		},
		{
			name:     "patch version higher",
			v1:       "v1.0.2",
			v2:       "v1.0.1",
			expected: true,
		},
		{
			name:     "version without v prefix",
			v1:       "2.0.0",
			v2:       "1.0.0",
			expected: true,
		},
		// v1 <= v2 cases
		{
			name:     "equal versions",
			v1:       "v1.0.0",
			v2:       "v1.0.0",
			expected: false,
		},
		{
			name:     "v1 older than v2",
			v1:       "v1.0.0",
			v2:       "v2.0.0",
			expected: false,
		},
		{
			name:     "minor version lower",
			v1:       "v1.0.0",
			v2:       "v1.1.0",
			expected: false,
		},
		{
			name:     "patch version lower",
			v1:       "v1.0.0",
			v2:       "v1.0.1",
			expected: false,
		},
		// Edge cases
		{
			name:     "partial version v1",
			v1:       "v1.2",
			v2:       "v1.1.0",
			expected: true,
		},
		{
			name:     "partial version v2",
			v1:       "v1.1.0",
			v2:       "v1.2",
			expected: false,
		},
		{
			name:     "mixed format",
			v1:       "v2.0.0",
			v2:       "1.9.9",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareVersions(tt.v1, tt.v2)
			if result != tt.expected {
				t.Errorf("CompareVersions(%q, %q) = %v, want %v", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

// TestParseVersion verifies that version strings are correctly parsed into components.
func TestParseVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected [3]int
	}{
		{
			name:     "full version",
			version:  "1.2.3",
			expected: [3]int{1, 2, 3},
		},
		{
			name:     "two parts",
			version:  "1.2",
			expected: [3]int{1, 2, 0},
		},
		{
			name:     "one part",
			version:  "1",
			expected: [3]int{1, 0, 0},
		},
		{
			name:     "empty string",
			version:  "",
			expected: [3]int{0, 0, 0},
		},
		{
			name:     "with extra parts",
			version:  "1.2.3.4",
			expected: [3]int{1, 2, 3},
		},
		{
			name:     "invalid characters default to 0",
			version:  "a.b.c",
			expected: [3]int{0, 0, 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseVersion(tt.version)
			if result != tt.expected {
				t.Errorf("parseVersion(%q) = %v, want %v", tt.version, result, tt.expected)
			}
		})
	}
}

// TestGetCurrentVersion ensures the version string is returned correctly.
func TestGetCurrentVersion(t *testing.T) {
	app := &App{}

	version := app.GetCurrentVersion()

	// Version should not be empty
	if version == "" {
		t.Error("GetCurrentVersion() returned empty string")
	}

	// Version should start with 'v' by convention
	if version[0] != 'v' {
		t.Errorf("GetCurrentVersion() = %q, expected to start with 'v'", version)
	}
}

// TestNewApp verifies that NewApp creates a valid App instance.
func TestNewApp(t *testing.T) {
	app := NewApp()

	if app == nil {
		t.Fatal("NewApp() returned nil")
	}

	// Initially, context should be nil (set during startup)
	if app.ctx != nil {
		t.Error("Expected ctx to be nil before startup")
	}
}
