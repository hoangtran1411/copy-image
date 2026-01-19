package main_test

import (
	"testing"
)

// Helper function example
func add(a, b int) int {
	return a + b
}

// Table-driven test example
func TestAdd(t *testing.T) {
	// Define test cases struct
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Positive numbers",
			args: args{a: 2, b: 3},
			want: 5,
		},
		{
			name: "Negative numbers",
			args: args{a: -1, b: -1},
			want: -2,
		},
		{
			name: "Mixed numbers",
			args: args{a: -5, b: 5},
			want: 0,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := add(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("add() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Example of testing with setup/teardown
func TestWithSetup(t *testing.T) {
	// Setup code usually goes here or in TestMain
	t.Log("Setup done")

	t.Run("Sub-test 1", func(t *testing.T) {
		// Test logic
	})

	// Teardown code
	t.Cleanup(func() {
		t.Log("Cleanup done")
	})
}
