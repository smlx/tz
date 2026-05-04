package parser

import (
	"strings"
	"time"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// TimeSpecAST represents the parsed abstract syntax tree of a time
// specification, optionally including a time of day and a weekday.
type TimeSpecAST struct {
	Time    *TimePart `parser:"@@?"`
	Weekday string    `parser:"@Ident?"`
}

// TimePart represents the parsed time of day, including hour, minute, and an
// optional AM/PM modifier.
type TimePart struct {
	Hour   int    `parser:"@Int"`
	Minute int    `parser:"(':' @Int)?"`
	AMPM   string `parser:"@Ident?"` // am or pm
}

var (
	// timeLexer defines the tokenization rules for parsing time specifications
	timeLexer = lexer.MustSimple([]lexer.SimpleRule{
		{Name: "Ident", Pattern: `[a-zA-Z]+`},
		{Name: "Int", Pattern: `[0-9]+`},
		{Name: "Punct", Pattern: `:`},
		{Name: "whitespace", Pattern: `\s+`},
	})
	// timeParser is the participle parser for TimeSpecAST using the timeLexer
	timeParser = participle.MustBuild[TimeSpecAST](
		participle.Lexer(timeLexer),
		participle.Elide("whitespace"),
		participle.UseLookahead(2),
	)
)

// Evaluate parses a time specification string and applies it to a base time,
// returning the resulting evaluated time. It handles time of day, AM/PM, and
// weekday adjustments.
func Evaluate(timeSpec string, base time.Time) (time.Time, error) {
	if strings.TrimSpace(timeSpec) == "" || timeSpec == "@" {
		return base, nil
	}
	ast, err := timeParser.ParseString("", strings.ToLower(timeSpec))
	if err != nil {
		return base, err
	}
	result := base
	// if a time is specified
	if ast.Time != nil {
		hour := ast.Time.Hour
		if ast.Time.AMPM == "pm" && hour < 12 {
			hour += 12
		} else if ast.Time.AMPM == "am" && hour == 12 {
			hour = 0
		}
		result = time.Date(result.Year(), result.Month(), result.Day(), hour, ast.Time.Minute, 0, 0, result.Location())
		// if the parsed time is before the current time and no weekday is
		// specified, advance to tomorrow
		if ast.Weekday == "" && result.Before(base) {
			result = result.AddDate(0, 0, 1)
		}
	}
	// if a weekday is specified, advance to it
	if ast.Weekday != "" {
		targetWeekday := parseWeekday(ast.Weekday)
		if targetWeekday != -1 {
			daysToAdd := int(targetWeekday) - int(result.Weekday())
			if daysToAdd < 0 || (daysToAdd == 0 && result.Before(base)) {
				daysToAdd += 7
			}
			result = result.AddDate(0, 0, daysToAdd)
		}
	}
	return result, nil
}

// parseWeekday converts a string representation of a weekday (full or
// abbreviated) into its corresponding time.Weekday value, returning -1 if
// unrecognized.
func parseWeekday(s string) time.Weekday {
	switch strings.ToLower(s) {
	case "sunday", "sun":
		return time.Sunday
	case "monday", "mon":
		return time.Monday
	case "tuesday", "tue":
		return time.Tuesday
	case "wednesday", "wed":
		return time.Wednesday
	case "thursday", "thu":
		return time.Thursday
	case "friday", "fri":
		return time.Friday
	case "saturday", "sat":
		return time.Saturday
	}
	return -1
}
