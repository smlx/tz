package location

import (
	"testing"
)

func TestFind(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		expectIANA string
	}{
		{"exact match", "new york", "America/New_York"},
		{"exact match with underscores", "New_York", "America/New_York"},
		{"fuzzy match zurich", "Zürich", "Europe/Zurich"},
		{"fuzzy match zurich lowercase", "zurich", "Europe/Zurich"},
		{"fuzzy match sydney", "sydney", "Australia/Sydney"},
		{"fuzzy match richmond,ca", "richmond,ca", "America/Vancouver"},
		{"fuzzy match richmond,va", "richmond,va", "America/New_York"},
		{"fuzzy match london", "london", "Europe/London"},
		{"local time", "@", "Local"},
		{"shorthand UTC", "UTC", "UTC"},
		{"shorthand UTC+8", "UTC+8", "UTC+8"},
		{"shorthand +0800", "+0800", "UTC+8"},
		{"shorthand -07:00", "-07:00", "UTC-7"},
		{"shorthand UTC-07:00", "UTC-07:00", "UTC-7"},
		{"shorthand GMT-7", "GMT-7", "UTC-7"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loc, _, err := Find(tt.input)
			if err != nil {
				t.Fatalf("unexpected error finding %s: %v", tt.input, err)
			}
			if loc.String() != tt.expectIANA {
				t.Errorf("Find(%q) = %v, want %v", tt.input, loc.String(), tt.expectIANA)
			}
		})
	}
}
