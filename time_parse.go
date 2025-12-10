package friendlytime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// unitConversion defines time unit conversions to hours
type unitConversion struct {
	patterns []string
	toHours  int
}

var (
	timeUnitsToHours = []unitConversion{
		{patterns: []string{"years", "year"}, toHours: 8760},  // 365 days
		{patterns: []string{"months", "month"}, toHours: 720}, // 30 days
		{patterns: []string{"days", "day"}, toHours: 24},
	}

	simpleTimeUnits = []struct {
		patterns []string
		to       string
	}{
		{patterns: []string{"seconds", "second", "sec"}, to: "s"},
		{patterns: []string{"minutes", "minute", "min"}, to: "m"},
		{patterns: []string{"hours", "hour"}, to: "h"},
	}

	reservedKeywords = []string{
		"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday",
		"yesterday", "last week", "last month", "last year",
	}
)

// ParseTimeRange parses human-readable time range to UNIX timestamps.
// It returns start and end timestamps in seconds since Unix epoch.
//
// The function supports various formats:
//   - Single time values: "1h" (start and end are the same)
//   - Time ranges with "/": "1h/30m" (from 1 hour ago to 30 minutes ago)
//   - Empty start or end: "/now" or "yesterday/" (empty means zero or now)
//   - Relative offsets: "1h/+30m" (from 1 hour ago, plus 30 minutes from that)
//
// Examples:
//   - "1h" -> exactly 1 hour ago
//   - "00:00/12m" -> from today's 00:00 till 12 minutes ago
//   - "4d/2 minutes ago" -> from 4 days ago till 2 minutes ago
//   - "1h/+1m" -> from 1 hour ago till 59 minutes ago (1h ago + 1m)
//   - "03:30/+2h" -> from today's 3:30 till today's 5:30
//   - "1416434697" -> exactly Unix timestamp 1416434697
//   - "14-11-19 22:04:57/" -> from 2014-11-19 22:04:57 till now
//   - "/1416434697" -> from zero time till the specified timestamp
//   - "40 minutes ago/35 min ago" -> from 40 minutes ago till 35 minutes ago
//   - "last monday/yesterday" -> from last monday 00:00:00 till yesterday 00:00:00
//
// Returns:
//   - start: Unix timestamp in seconds for the start of the range
//   - end: Unix timestamp in seconds for the end of the range
//   - error: An error if parsing fails or if end time is before start time
func ParseTimeRange(timeRange string) (int64, int64, error) {
	if timeRange == "" {
		return 0, 0, nil
	}

	now := time.Now()

	var startTime, endTime time.Time

	var err error

	if strings.Contains(timeRange, "/") {
		parts := strings.Split(timeRange, "/")
		if len(parts) != 2 {
			return 0, 0, ErrInvalidTimeRange
		}

		startTime, err = ParseTime(parts[0], now, time.Time{})
		if err != nil {
			return 0, 0, fmt.Errorf("%w: %v", ErrInvalidStartTime, err)
		}

		endTime, err = ParseTime(parts[1], now, startTime)
		if err != nil {
			return 0, 0, fmt.Errorf("%w: %v", ErrInvalidEndTime, err)
		}
	} else {
		startTime, err = ParseTime(timeRange, now, time.Time{})
		if err != nil {
			return 0, 0, err
		}

		endTime = startTime
	}

	// Validate that end time is not before start time
	if !endTime.IsZero() && !startTime.IsZero() && endTime.Before(startTime) {
		return 0, 0, ErrEndBeforeStart
	}

	return startTime.Unix(), endTime.Unix(), nil
}

// ParseTime parses a human-readable time string to a time.Time value.
//
// The function accepts various formats including:
//   - Durations: "1h", "30m", "45s" (relative to now)
//   - Custom units: "5 days ago", "2 months", "1 year"
//   - Weekdays: "last monday", "yesterday"
//   - Time of day: "15:30", "09:00"
//   - Dates: "2006-01-02", "06-01-02 15:04:05"
//   - Unix timestamps: "1416434697"
//   - Relative offsets: "+30m" (relative to startTime), "-15m" (relative to now)
//
// Parameters:
//   - timeStr: The time string to parse
//   - now: The reference time for relative calculations (usually time.Now())
//   - startTime: Used for relative offsets with "+" prefix (e.g., "+30m" means startTime + 30 minutes)
//
// Examples:
//   - "1h" -> 1 hour before now
//   - "00:00" -> today at midnight
//   - "4d" -> 4 days ago
//   - "last monday" -> previous Monday at 00:00:00
//   - "yesterday" -> yesterday at 00:00:00
//   - "+30m" -> 30 minutes after startTime
//
// Returns a time.Time value or an error if the format is not recognized.
func ParseTime(timeStr string, now time.Time, startTime time.Time) (time.Time, error) {
	if timeStr == "" {
		if startTime.IsZero() {
			return startTime, nil
		}
		return now, nil
	}

	timeStr = strings.TrimSpace(timeStr)

	// Try to parse as Unix timestamp
	if unixTime, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
		return parseUnixTimestamp(unixTime)
	}

	lowerTimeStr := strings.ToLower(timeStr)

	if strings.HasPrefix(lowerTimeStr, "last ") || strings.HasPrefix(lowerTimeStr, "yesterday") {
		return parseRelativeTime(lowerTimeStr, now, startTime, false)
	}

	// Handle "N units ago" format - convert units first, then remove "ago"
	if strings.Contains(timeStr, " ago") {
		cleanStr := strings.Replace(timeStr, " ago", "", 1)
		cleanStr, err := convertCustomUnits(cleanStr)
		if err != nil {
			return time.Time{}, fmt.Errorf("failed to convert custom units: %w", err)
		}

		cleanStr = strings.Join(strings.Fields(cleanStr), "")

		return parseRelativeTime(cleanStr, now, startTime, false)
	}

	timeStr, err := convertCustomUnits(timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to convert custom units: %w", err)
	}

	if strings.HasPrefix(timeStr, "+") {
		return parseRelativeTime(timeStr[1:], now, startTime, true)
	}

	if strings.HasPrefix(timeStr, "-") {
		return parseRelativeTime(timeStr[1:], now, startTime, false)
	}

	if duration, err := time.ParseDuration(timeStr); err == nil {
		return now.Add(-duration), nil
	}

	if t, err := time.Parse("15:04", timeStr); err == nil {
		return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location()), nil
	}

	if t, err := time.Parse("06-01-02", timeStr); err == nil {
		return t, nil
	}

	if t, err := time.Parse("2006-01-02", timeStr); err == nil {
		return t, nil
	}

	if t, err := time.Parse("06-01-02 15:04:05", timeStr); err == nil {
		return t, nil
	}

	if t, err := time.Parse("2006-01-02 15:04:05", timeStr); err == nil {
		return t, nil
	}

	if t, err := time.Parse("Mon, 02 Jan 2006 15:04:05", timeStr); err == nil {
		return t, nil
	}

	return time.Time{}, ErrInvalidTimeFormat
}

func convertCustomUnits(timeStr string) (string, error) {
	// Don't convert weekday names or special keywords
	lowerTimeStr := strings.ToLower(timeStr)
	for _, keyword := range reservedKeywords {
		if strings.Contains(lowerTimeStr, keyword) {
			return timeStr, nil
		}
	}

	// Order matters: process longer patterns first to avoid partial matches
	for _, conv := range timeUnitsToHours {
		for _, pattern := range conv.patterns {
			if converted, ok := tryConvertUnit(timeStr, pattern, conv.toHours); ok {
				return converted, nil
			}
		}
	}

	// Handle short forms: d, y (only if not part of a longer word)
	if result := tryConvertShortForm(timeStr, "d", 24, "day"); result != "" {
		return result, nil
	}
	if result := tryConvertShortForm(timeStr, "y", 8760, "year"); result != "" {
		return result, nil
	}

	for _, conv := range simpleTimeUnits {
		for _, pattern := range conv.patterns {
			if converted, ok := tryConvertSimpleUnit(timeStr, pattern, conv.to); ok {
				return converted, nil
			}
		}
	}

	// If no conversion matched, return the original string
	return timeStr, nil
}

// tryConvertUnit attempts to convert a time string with a unit pattern to hours
func tryConvertUnit(timeStr, pattern string, toHours int) (string, bool) {
	// Try "N pattern" format (e.g., "2 days")
	withSpace := " " + pattern
	if idx := strings.Index(timeStr, withSpace); idx > 0 {
		prefix := timeStr[:idx]

		var amount int
		if n, _ := fmt.Sscanf(prefix, "%d", &amount); n == 1 && amount >= 0 {
			endIdx := idx + len(withSpace)
			if endIdx >= len(timeStr) || !isLetter(rune(timeStr[endIdx])) {
				hours := amount * toHours

				suffix := ""
				if endIdx < len(timeStr) {
					suffix = timeStr[endIdx:]
				}

				result := fmt.Sprintf("%dh%s", hours, suffix)

				return result, true
			}
		}
	}

	// Try "Npattern" format (e.g., "2days")
	if idx := strings.Index(timeStr, pattern); idx > 0 {
		prefix := timeStr[:idx]

		var amount int
		if n, _ := fmt.Sscanf(prefix, "%d", &amount); n == 1 && amount >= 0 {
			endIdx := idx + len(pattern)
			if endIdx >= len(timeStr) || !isLetter(rune(timeStr[endIdx])) {
				hours := amount * toHours

				suffix := ""
				if endIdx < len(timeStr) {
					suffix = timeStr[endIdx:]
				}

				result := fmt.Sprintf("%dh%s", hours, suffix)

				return result, true
			}
		}
	}

	return "", false
}

// tryConvertSimpleUnit attempts to convert a time string with a simple unit pattern
func tryConvertSimpleUnit(timeStr, pattern, to string) (string, bool) {
	// Try "N pattern" format (e.g., "10 seconds")
	// Check for word boundary after pattern (space, end of string, or non-letter)
	withSpace := " " + pattern
	if idx := strings.Index(timeStr, withSpace); idx > 0 {
		// Extract the number before the pattern
		prefix := timeStr[:idx]

		var amount int
		if n, _ := fmt.Sscanf(prefix, "%d", &amount); n == 1 && amount >= 0 {
			endIdx := idx + len(withSpace)
			if endIdx >= len(timeStr) || !isLetter(rune(timeStr[endIdx])) {
				suffix := ""
				if endIdx < len(timeStr) {
					suffix = timeStr[endIdx:]
				}

				result := fmt.Sprintf("%d%s%s", amount, to, suffix)

				return result, true
			}
		}
	}

	// Try "Npattern" format (e.g., "10seconds")
	if idx := strings.Index(timeStr, pattern); idx > 0 {
		prefix := timeStr[:idx]

		var amount int
		if n, _ := fmt.Sscanf(prefix, "%d", &amount); n == 1 && amount >= 0 {
			endIdx := idx + len(pattern)
			if endIdx >= len(timeStr) || !isLetter(rune(timeStr[endIdx])) {
				suffix := ""
				if endIdx < len(timeStr) {
					suffix = timeStr[endIdx:]
				}

				result := fmt.Sprintf("%d%s%s", amount, to, suffix)

				return result, true
			}
		}
	}

	return "", false
}

// isLetter checks if a rune is a letter
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// tryConvertShortForm attempts to convert short form units (d, y) to hours
func tryConvertShortForm(timeStr, shortForm string, toHours int, longForm string) string {
	if !strings.Contains(timeStr, shortForm) || strings.Contains(timeStr, longForm) {
		return ""
	}

	// Try to extract "Nshortform" pattern (e.g., "5d", "1y")
	var amount int

	searchPattern := fmt.Sprintf("%%d%s", shortForm)
	if n, _ := fmt.Sscanf(timeStr, searchPattern, &amount); n == 1 && amount >= 0 {
		target := fmt.Sprintf("%d%s", amount, shortForm)

		idx := strings.Index(timeStr, target)
		if idx >= 0 {
			endIdx := idx + len(target)
			if endIdx >= len(timeStr) || !isLetter(rune(timeStr[endIdx])) {
				hours := amount * toHours

				return strings.Replace(timeStr, target, fmt.Sprintf("%dh", hours), 1)
			}
		}
	}

	return ""
}

func parseRelativeTime(durationStr string, now, startTime time.Time, isPositive bool) (time.Time, error) {
	durationStr = strings.ToLower(strings.TrimSpace(durationStr))
	switch durationStr {
	case "yesterday":
		// Return yesterday at midnight
		return getMidnight(now.AddDate(0, 0, -1)), nil
	case "last week":
		// Return 7 days ago at midnight
		return getMidnight(now.AddDate(0, 0, -7)), nil
	case "last month":
		// Return 1 month ago at midnight
		return getMidnight(now.AddDate(0, -1, 0)), nil
	case "last year":
		// Return 1 year ago at midnight
		return getMidnight(now.AddDate(-1, 0, 0)), nil
	default:
		if strings.HasPrefix(durationStr, "last ") {
			return parseLastWeekday(durationStr, now)
		}

		duration, err := time.ParseDuration(strings.ReplaceAll(durationStr, " ", ""))
		if err != nil {
			return time.Time{}, err
		}

		if isPositive {
			if startTime.IsZero() {
				return now.Add(duration), nil
			}

			return startTime.Add(duration), nil
		}

		return now.Add(-duration), nil
	}
}

// getMidnight returns the time at midnight (00:00:00) for the given date
func getMidnight(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func parseLastWeekday(durationStr string, now time.Time) (time.Time, error) {
	weekdayStr := strings.TrimPrefix(durationStr, "last ")

	weekday, err := parseWeekday(weekdayStr)
	if err != nil {
		return time.Time{}, err
	}

	// Calculate the date for the last occurrence of the specified weekday
	daysAgo := int(now.Weekday() - weekday)
	if daysAgo <= 0 {
		daysAgo += 7
	}

	lastWeekday := now.AddDate(0, 0, -daysAgo)

	return getMidnight(lastWeekday), nil
}

// parseUnixTimestamp parses a Unix timestamp, handling both seconds and milliseconds
func parseUnixTimestamp(unixTime int64) (time.Time, error) {
	// Detect if it's milliseconds (13 digits or more) vs seconds (10 digits or less)
	// Timestamps > 9999999999 are either after year 2286 in seconds, or milliseconds
	if unixTime > 9999999999 {
		// Treat as milliseconds
		sec := unixTime / 1000
		nsec := (unixTime % 1000) * 1000000

		return time.Unix(sec, nsec), nil
	}

	// Treat as seconds (including negative timestamps for dates before 1970)
	return time.Unix(unixTime, 0), nil
}

func parseWeekday(weekdayStr string) (time.Weekday, error) {
	switch weekdayStr {
	case "sunday":
		return time.Sunday, nil
	case "monday":
		return time.Monday, nil
	case "tuesday":
		return time.Tuesday, nil
	case "wednesday":
		return time.Wednesday, nil
	case "thursday":
		return time.Thursday, nil
	case "friday":
		return time.Friday, nil
	case "saturday":
		return time.Saturday, nil
	default:
		return time.Sunday, fmt.Errorf("%w: %s", ErrInvalidWeekday, weekdayStr)
	}
}
