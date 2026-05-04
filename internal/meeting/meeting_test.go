package meeting

import (
	"strings"
	"testing"
	"time"
)

func TestPlan(t *testing.T) {
	now := time.Date(2023, 10, 25, 10, 30, 0, 0, time.UTC)
	locations := []string{"Zurich", "Tokyo"}

	res, err := Plan(locations, now)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(res, "UTC") {
		t.Errorf("expected UTC in output, got: %s", res)
	}
	if !strings.Contains(res, "Local") {
		t.Errorf("expected Local in output, got: %s", res)
	}
	if !strings.Contains(res, "Europe/Zurich") {
		t.Errorf("expected Europe/Zurich header in output, got: %s", res)
	}
	if !strings.Contains(res, "CEST UTC+2") && !strings.Contains(res, "CET UTC+1") {
		t.Errorf("expected Zurich timezone header in output, got: %s", res)
	}
	if !strings.Contains(res, "Asia/Tokyo") {
		t.Errorf("expected Asia/Tokyo header in output, got: %s", res)
	}
	if !strings.Contains(res, "JST UTC+9") {
		t.Errorf("expected Tokyo timezone header in output, got: %s", res)
	}
	if !strings.Contains(res, "<-- current time") {
		t.Errorf("expected <-- current time in output, got: %s", res)
	}
}
