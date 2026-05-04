package converter

import (
	"testing"
	"time"
)

func TestConvert(t *testing.T) {
	// Mock local time to UTC for predictable testing
	origLocal := time.Local
	t.Cleanup(func() {
		time.Local = origLocal
	})
	time.Local, _ = time.LoadLocation("UTC")
	now := time.Date(2023, 10, 25, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		target   string
		source   string
		timeSpec []string
		expected string
	}{
		{
			name:     "local to target now",
			target:   "Zurich",
			source:   "@",
			timeSpec: []string{},
			expected: "Zürich (Europe/Zurich CEST UTC+2)\nWed, 25 Oct 2023 12:00:00 CEST",
		},
		{
			name:     "source to target now",
			target:   "Tokyo",
			source:   "Bangalore",
			timeSpec: []string{},
			// Bangalore (IST) is UTC+5:30. At 10:00 UTC, it's 15:30 IST.
			// Tokyo (JST) is UTC+9. At 10:00 UTC, it's 19:00 JST.
			expected: "Tokyo (Asia/Tokyo JST UTC+9)\nWed, 25 Oct 2023 19:00:00 JST",
		},
		{
			name:     "source to target with time",
			target:   "Zurich",
			source:   "Sydney",
			timeSpec: []string{"5am"},
			// base is 10:00 UTC -> 21:00 AEDT in Sydney.
			// Next 5am in Sydney is 2023-10-26 05:00:00 AEDT -> 2023-10-25 18:00:00 UTC
			// 18:00 UTC in Zurich (CEST, UTC+2) is 20:00 CEST.
			expected: "Zürich (Europe/Zurich CEST UTC+2)\nWed, 25 Oct 2023 20:00:00 CEST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := Convert(tt.target, tt.source, tt.timeSpec, now)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res != tt.expected {
				t.Errorf("Convert() = %v, want %v", res, tt.expected)
			}
		})
	}
}
