package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/smlx/tz/internal/location"
	"github.com/smlx/tz/internal/parser"
)

// formatOutput formats the time output including timezone offset calculations.
func formatOutput(name string, t time.Time) string {
	zone, offsetSec := t.Zone()
	hours := offsetSec / 3600
	mins := (offsetSec % 3600) / 60
	if mins < 0 {
		mins = -mins
	}
	var offset string
	if mins == 0 {
		offset = fmt.Sprintf("UTC%+d", hours)
	} else {
		offset = fmt.Sprintf("UTC%+d:%02d", hours, mins)
	}

	const format = "Mon, 02 Jan 2006 15:04:05 MST"
	formattedTime := t.Format(format)

	return fmt.Sprintf("%s (%s %s %s)\n%s", name, t.Location().String(), zone, offset, formattedTime)
}

// Convert performs timezone conversion based on target, source, and time specification.
func Convert(target, source string, timeSpec []string, now time.Time) (string, error) {
	// resolve source
	var sourceLoc *time.Location
	var sourceName string
	var err error
	if source == "" || source == "@" {
		sourceLoc = time.Local
		sourceName = "Local"
	} else {
		sourceLoc, sourceName, err = location.Find(source)
		if err != nil {
			return "", fmt.Errorf("failed to find source location: %w", err)
		}
	}

	// resolve target
	var targetLoc *time.Location
	var targetName string
	if target == "" || target == "@" {
		targetLoc = time.Local
		targetName = "Local"
	} else {
		targetLoc, targetName, err = location.Find(target)
		if err != nil {
			return "", fmt.Errorf("failed to find target location: %w", err)
		}
	}

	// parse time specification
	specStr := strings.Join(timeSpec, " ")
	baseTime := now.In(sourceLoc)
	parsedTime, err := parser.Evaluate(specStr, baseTime)
	if err != nil {
		return "", fmt.Errorf("failed to parse time specification: %w", err)
	}

	// convert time spec to target time zone
	targetTime := parsedTime.In(targetLoc)

	return fmt.Sprintf("%s\n\n%s", formatOutput(sourceName, parsedTime), formatOutput(targetName, targetTime)), nil
}
