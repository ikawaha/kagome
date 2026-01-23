package main

import "testing"

// TestWouldOverflowTokenAllocation tests the overflow detection logic
func TestWouldOverflowTokenAllocation(t *testing.T) {
	tests := []struct {
		name     string
		n        int
		wantFail bool
	}{
		{
			name:     "negative number",
			n:        -1,
			wantFail: true,
		},
		{
			name:     "zero tokens",
			n:        0,
			wantFail: false,
		},
		{
			name:     "normal amount",
			n:        1000,
			wantFail: false,
		},
		{
			name:     "large but safe",
			n:        1000000,
			wantFail: false,
		},
		{
			name:     "maximum int",
			n:        int(^uint(0) >> 1), // MaxInt
			wantFail: true,              // Should overflow on most platforms
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := wouldOverflowTokenAllocation(tt.n)
			if got != tt.wantFail {
				t.Errorf("wouldOverflowTokenAllocation(%d) = %v, want %v", tt.n, got, tt.wantFail)
			}
		})
	}
}
