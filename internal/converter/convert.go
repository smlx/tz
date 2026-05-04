package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/smlx/tz/internal/location"
	"github.com/smlx/tz/internal/parser"
)

// Convert performs timezone conversion based on target, source, and time specification.
func Convert(target, source string, timeSpec []string, now time.Time) (string, error) {
	// resolve source
	var sourceLoc *time.Location
	var err error
	if source == "" || source == "@" {
		sourceLoc = time.Local
	} else {
		sourceLoc, _, err = location.Find(source)
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
	// format and return output
	zone, offsetSec := targetTime.Zone()
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
	formattedTime := targetTime.Format(format)
	return fmt.Sprintf("%s (%s %s %s)\n%s", targetName, targetLoc.String(), zone, offset, formattedTime), nil
}
