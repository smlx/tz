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
	// Check if all arguments together form a valid time specification, which
	// means target and source were omitted (defaulting to local time).
	var fullArgs []string
	if target != "" {
		fullArgs = append(fullArgs, target)
	}
	if source != "" {
		fullArgs = append(fullArgs, source)
	}
	fullArgs = append(fullArgs, timeSpec...)
	if len(fullArgs) > 0 {
		fullSpec := strings.Join(fullArgs, " ")
		if _, parseErr := parser.Evaluate(fullSpec, now.In(time.Local)); parseErr == nil {
			target = "@"
			source = "@"
			timeSpec = fullArgs
		}
	}

	// resolve target
	var targetLoc *time.Location
	var targetName string
	var err error
	if target == "" || target == "@" {
		targetLoc = time.Local
		targetName = "Local"
	} else {
		targetLoc, targetName, err = location.Find(target)
		if err != nil {
			return "", fmt.Errorf("failed to find target location: %w", err)
		}
	}

	// resolve source
	var sourceLoc *time.Location
	var sourceName string
	if source == "" || source == "@" {
		sourceLoc = time.Local
		sourceName = "Local"
	} else {
		// Check if source + timeSpec forms a valid time specification first,
		// before matching source as a location name.
		fullSpec := strings.Join(append([]string{source}, timeSpec...), " ")
		if _, parseErr := parser.Evaluate(fullSpec, now.In(time.Local)); parseErr == nil {
			sourceLoc = time.Local
			sourceName = "Local"
			timeSpec = append([]string{source}, timeSpec...)
			source = "@"
		} else {
			sourceLoc, sourceName, err = location.Find(source)
			if err != nil {
				return "", fmt.Errorf("failed to find source location: %w", err)
			}
		}
	}

	// parse time specification
	specStr := strings.Join(timeSpec, " ")
	baseTime := now.In(sourceLoc)
	parsedTime, err := parser.Evaluate(specStr, baseTime)
	if err != nil {
		return "", fmt.Errorf("failed to parse time specification: %w", err)
	}
	if (source == "" || source == "@") && parsedTime.Location() != sourceLoc {
		sourceName = parsedTime.Location().String()
	}

	// convert time spec to target time zone
	targetTime := parsedTime.In(targetLoc)

	return fmt.Sprintf("%s\n\n%s", formatOutput(sourceName, parsedTime), formatOutput(targetName, targetTime)), nil
}
