package meeting

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/smlx/tz/internal/location"
)

// Plan creates a tabular meeting planner centered around `now` for the given
// locations.
func Plan(locations []string, now time.Time) (string, error) {
	// parse and resolve location names to their respective timezones
	locs := make([]*time.Location, len(locations))
	origNames := make([]string, len(locations))
	for i, name := range locations {
		loc, origName, err := location.Find(name)
		if err != nil {
			return "", fmt.Errorf("failed to find location %s: %w", name, err)
		}
		locs[i] = loc
		origNames[i] = origName
	}
	// initialize a tabwriter to format the output into aligned columns
	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 0, 0, 2, ' ', 0)
	// local location details
	localZone, _ := now.In(time.Local).Zone()
	localOffset := offsetString(now.In(time.Local))
	// write header rows
	header0 := []string{"", ""}
	header1 := []string{"UTC", "Local"}
	header2 := []string{"", fmt.Sprintf("%s %s", localZone, localOffset)}
	for i := range locations {
		// calculate offset from UTC for the current time to display in header
		zone, _ := now.In(locs[i]).Zone()
		offset := offsetString(now.In(locs[i]))
		header0 = append(header0, origNames[i])
		header1 = append(header1, locs[i].String())
		header2 = append(header2, fmt.Sprintf("%s %s", zone, offset))
	}
	header0 = append(header0, "")
	header1 = append(header1, "")
	header2 = append(header2, "")
	fmt.Fprintln(tw, strings.Join(header0, "\t"))
	fmt.Fprintln(tw, strings.Join(header1, "\t"))
	fmt.Fprintln(tw, strings.Join(header2, "\t"))
	// write rows (now - 12h to now + 12h)
	baseHour := now.UTC().Truncate(time.Hour)
	for i := -12; i <= 12; i++ {
		t := baseHour.Add(time.Duration(i) * time.Hour)
		row := []string{t.Format("15:04"), t.In(time.Local).Format("15:04")}
		for _, loc := range locs {
			row = append(row, t.In(loc).Format("15:04"))
		}
		row = append(row, "")
		fmt.Fprintln(tw, strings.Join(row, "\t"))
		// if this is the current hour, insert the exact current time after it
		if i == 0 {
			exactRow := []string{
				now.In(time.UTC).Format("15:04"),
				now.In(time.Local).Format("15:04"),
			}
			for _, loc := range locs {
				exactRow = append(exactRow, now.In(loc).Format("15:04"))
			}
			exactRow = append(exactRow, "<-- current time")
			fmt.Fprintln(tw, strings.Join(exactRow, "\t"))
		}
	}
	// ensure all data is written to the underlying builder
	tw.Flush()
	// apply bold to the row containing the exact current time
	output := sb.String()
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if strings.Contains(line, "<-- current time") {
			lines[i] = fmt.Sprintf("\033[1m%s\033[0m", line)
		}
	}
	return strings.Join(lines, "\n"), nil
}

// offsetString calculates and formats the UTC offset for a given time as a
// string (e.g., UTC+10)
func offsetString(t time.Time) string {
	_, offset := t.Zone()
	hours := offset / 3600
	mins := (offset % 3600) / 60
	if mins < 0 {
		mins = -mins
	}
	if mins == 0 {
		return fmt.Sprintf("UTC%+d", hours)
	}
	return fmt.Sprintf("UTC%+d:%02d", hours, mins)
}
