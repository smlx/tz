package parser

import (
	"testing"
	"time"
)

func TestEvaluate(t *testing.T) {
	loc, _ := time.LoadLocation("UTC")
	base := time.Date(2023, 10, 25, 10, 0, 0, 0, loc) // Wednesday, Oct 25 2023, 10:00

	tests := []struct {
		name     string
		timeSpec string
		expected time.Time
	}{
		{"empty", "", base},
		{"at", "@", base},
		{"5am", "5am", time.Date(2023, 10, 26, 5, 0, 0, 0, loc)}, // tomorrow 5am
		{"15:00", "15:00", time.Date(2023, 10, 25, 15, 0, 0, 0, loc)}, // today 3pm
		{"5am friday", "5am friday", time.Date(2023, 10, 27, 5, 0, 0, 0, loc)}, // this coming friday
		{"friday", "friday", time.Date(2023, 10, 27, 10, 0, 0, 0, loc)}, // this coming friday at current time
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Evaluate(tt.timeSpec, base)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !res.Equal(tt.expected) {
				t.Errorf("Evaluate(%q) = %v, want %v", tt.timeSpec, res, tt.expected)
			}
		})
	}
}
