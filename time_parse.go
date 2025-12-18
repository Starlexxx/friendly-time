// Package friendlytime provides utilities for parsing human-readable time strings and ranges.
package friendlytime

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Time conversion constants.
	hoursPerYear  = 8760 // 365 days * 24 hours
	hoursPerMonth = 720  // 30 days * 24 hours
	hoursPerDay   = 24

	// Timestamp boundaries.
	partsCountInRange          = 2
	timestampMillisecondBorder = 9999999999 // Timestamps > this are treated as milliseconds
	millisecondsPerSecond      = 1000
	nanosecondsPerMillisecond  = 1000000
)

// unitConversion defines time unit conversions to hours.
type unitConversion struct {
	patterns []string
	toHours  int
}

// simpleUnit defines simple time unit conversions.
type simpleUnit struct {
	patterns []string
	to       string
}

// getTimeUnitsToHours returns unit conversions that need to be converted to hours.
func getTimeUnitsToHours() []unitConversion {
	return []unitConversion{
		{patterns: []string{"years", "year"}, toHours: hoursPerYear},
		{patterns: []string{"months", "month"}, toHours: hoursPerMonth},
		{patterns: []string{"days", "day"}, toHours: hoursPerDay},
	}
}

// getSimpleTimeUnits returns simple time unit conversions.
func getSimpleTimeUnits() []simpleUnit {
	return []simpleUnit{
		{patterns: []string{"seconds", "second", "sec"}, to: "s"},
		{patterns: []string{"minutes", "minute", "min"}, to: "m"},
		{patterns: []string{"hours", "hour"}, to: "h"},
	}
}

// getReservedKeywords returns keywords that should not be converted.
func getReservedKeywords() []string {
	return []string{
		"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday",
		"yesterday", "last week", "last month", "last year",
	}
}

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

	if !strings.Contains(timeRange, "/") {
		return parseSingleTime(timeRange, now)
	}

	return parseTimeRangeParts(timeRange, now)
}

// parseSingleTime parses a single time value (no range).
func parseSingleTime(timeRange string, now time.Time) (int64, int64, error) {
	startTime, err := ParseTime(timeRange, now, time.Time{})
	if err != nil {
		return 0, 0, err
	}

	return startTime.Unix(), startTime.Unix(), nil
}

// parseTimeRangeParts parses a time range with "/" separator.
func parseTimeRangeParts(timeRange string, now time.Time) (int64, int64, error) {
	parts := strings.Split(timeRange, "/")
	if len(parts) != partsCountInRange {
		return 0, 0, ErrInvalidTimeRange
	}

	startTime, err := ParseTime(parts[0], now, time.Time{})
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrInvalidStartTime, err)
	}

	endTime, err := ParseTime(parts[1], now, startTime)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: %w", ErrInvalidEndTime, err)
	}

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
func ParseTime(timeStr string, now, startTime time.Time) (time.Time, error) {
	if timeStr == "" {
		return handleEmptyTime(startTime, now), nil
	}

	timeStr = strings.TrimSpace(timeStr)

	// Try to parse as Unix timestamp
	if t, ok := tryParseUnixTimestamp(timeStr); ok {
		return t, nil
	}

	// Try relative time formats
	if t, ok, err := tryParseRelativeFormats(timeStr, now, startTime); ok || err != nil {
		return t, err
	}

	// Try duration and date formats
	return tryParseDateFormats(timeStr, now)
}

// handleEmptyTime returns appropriate time for empty input.
func handleEmptyTime(startTime, now time.Time) time.Time {
	if startTime.IsZero() {
		return startTime
	}

	return now
}

// tryParseUnixTimestamp attempts to parse as Unix timestamp.
func tryParseUnixTimestamp(timeStr string) (time.Time, bool) {
	unixTime, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return time.Time{}, false
	}

	t := parseUnixTimestamp(unixTime)

	return t, true
}

// tryParseRelativeFormats attempts to parse relative time formats.
func tryParseRelativeFormats(timeStr string, now, startTime time.Time) (time.Time, bool, error) {
	lowerTimeStr := strings.ToLower(timeStr)

	// Check for "last" keywords and "yesterday"
	if strings.HasPrefix(lowerTimeStr, "last ") || strings.HasPrefix(lowerTimeStr, "yesterday") {
		t, err := parseRelativeTime(lowerTimeStr, now, startTime, false)

		return t, true, err
	}

	// Handle "N units ago" format
	if strings.Contains(timeStr, " ago") {
		t, err := parseAgoFormat(timeStr, now, startTime)

		return t, true, err
	}

	// Handle + and - prefixes
	if t, ok, err := tryParsePrefixedTime(timeStr, now, startTime); ok {
		return t, true, err
	}

	return time.Time{}, false, nil
}

// parseAgoFormat handles "N units ago" format.
func parseAgoFormat(timeStr string, now, startTime time.Time) (time.Time, error) {
	cleanStr := strings.Replace(timeStr, " ago", "", 1)

	cleanStr = convertCustomUnits(cleanStr)

	cleanStr = strings.Join(strings.Fields(cleanStr), "")

	return parseRelativeTime(cleanStr, now, startTime, false)
}

// tryParsePrefixedTime handles + and - prefixed times.
func tryParsePrefixedTime(timeStr string, now, startTime time.Time) (time.Time, bool, error) {
	if strings.HasPrefix(timeStr, "+") {
		converted := convertCustomUnits(timeStr[1:])

		t, err := parseRelativeTime(converted, now, startTime, true)

		return t, true, err
	}

	if strings.HasPrefix(timeStr, "-") {
		converted := convertCustomUnits(timeStr[1:])

		t, err := parseRelativeTime(converted, now, startTime, false)

		return t, true, err
	}

	return time.Time{}, false, nil
}

// tryParseDateFormats attempts to parse various date and time formats.
func tryParseDateFormats(timeStr string, now time.Time) (time.Time, error) {
	converted := convertCustomUnits(timeStr)

	// Try duration
	if duration, err := time.ParseDuration(converted); err == nil {
		return now.Add(-duration), nil
	}

	// Try time of day
	if t, ok := tryParseTimeOfDay(converted, now); ok {
		return t, nil
	}

	// Try standard date formats
	formats := []string{
		"06-01-02",
		"2006-01-02",
		"06-01-02 15:04:05",
		"2006-01-02 15:04:05",
		"Mon, 02 Jan 2006 15:04:05",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, converted); err == nil {
			return t, nil
		}
	}

	return time.Time{}, ErrInvalidTimeFormat
}

// tryParseTimeOfDay attempts to parse time of day format (HH:MM).
func tryParseTimeOfDay(timeStr string, now time.Time) (time.Time, bool) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		return time.Time{}, false
	}

	return time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		t.Hour(),
		t.Minute(),
		0,
		0,
		now.Location(),
	), true
}

func convertCustomUnits(timeStr string) string {
	if isReservedKeyword(timeStr) {
		return timeStr
	}

	if result := tryConvertToHours(timeStr); result != "" {
		return result
	}

	if result := tryConvertShortForms(timeStr); result != "" {
		return result
	}

	if result := tryConvertSimpleUnits(timeStr); result != "" {
		return result
	}

	return timeStr
}

// isReservedKeyword checks if the string contains reserved keywords.
func isReservedKeyword(timeStr string) bool {
	lowerTimeStr := strings.ToLower(timeStr)
	for _, keyword := range getReservedKeywords() {
		if strings.Contains(lowerTimeStr, keyword) {
			return true
		}
	}

	return false
}

// tryConvertToHours converts time units to hours.
func tryConvertToHours(timeStr string) string {
	for _, conv := range getTimeUnitsToHours() {
		for _, pattern := range conv.patterns {
			if converted, ok := tryConvertUnit(timeStr, pattern, conv.toHours); ok {
				return converted
			}
		}
	}

	return ""
}

// tryConvertShortForms handles short form conversions (d, y).
func tryConvertShortForms(timeStr string) string {
	if result := tryConvertShortForm(timeStr, "d", hoursPerDay, "day"); result != "" {
		return result
	}

	if result := tryConvertShortForm(timeStr, "y", hoursPerYear, "year"); result != "" {
		return result
	}

	return ""
}

// tryConvertSimpleUnits converts simple time units.
func tryConvertSimpleUnits(timeStr string) string {
	for _, conv := range getSimpleTimeUnits() {
		for _, pattern := range conv.patterns {
			if converted, ok := tryConvertSimpleUnit(timeStr, pattern, conv.to); ok {
				return converted
			}
		}
	}

	return ""
}

// tryConvertUnit attempts to convert a time string with a unit pattern to hours.
func tryConvertUnit(timeStr, pattern string, toHours int) (string, bool) {
	// Try "N pattern" format (e.g., "2 days")
	if result, ok := tryConvertWithPattern(timeStr, " "+pattern, toHours); ok {
		return result, true
	}

	// Try "Npattern" format (e.g., "2days")
	if result, ok := tryConvertWithPattern(timeStr, pattern, toHours); ok {
		return result, true
	}

	return "", false
}

// tryConvertWithPattern is a helper for tryConvertUnit.
func tryConvertWithPattern(timeStr, pattern string, toHours int) (string, bool) {
	idx := strings.Index(timeStr, pattern)
	if idx <= 0 {
		return "", false
	}

	prefix := timeStr[:idx]

	var amount int
	if n, _ := fmt.Sscanf(prefix, "%d", &amount); n != 1 || amount < 0 {
		return "", false
	}

	endIdx := idx + len(pattern)
	if endIdx < len(timeStr) && isLetter(rune(timeStr[endIdx])) {
		return "", false
	}

	hours := amount * toHours

	suffix := ""
	if endIdx < len(timeStr) {
		suffix = timeStr[endIdx:]
	}

	return fmt.Sprintf("%dh%s", hours, suffix), true
}

// tryConvertSimpleUnit attempts to convert a time string with a simple unit pattern.
func tryConvertSimpleUnit(timeStr, pattern, to string) (string, bool) {
	// Try "N pattern" format (e.g., "10 seconds")
	if result, ok := tryConvertSimpleWithPattern(timeStr, " "+pattern, to); ok {
		return result, true
	}

	// Try "Npattern" format (e.g., "10seconds")
	if result, ok := tryConvertSimpleWithPattern(timeStr, pattern, to); ok {
		return result, true
	}

	return "", false
}

// tryConvertSimpleWithPattern is a helper for tryConvertSimpleUnit.
func tryConvertSimpleWithPattern(timeStr, pattern, to string) (string, bool) {
	idx := strings.Index(timeStr, pattern)
	if idx <= 0 {
		return "", false
	}

	prefix := timeStr[:idx]

	var amount int
	if n, _ := fmt.Sscanf(prefix, "%d", &amount); n != 1 || amount < 0 {
		return "", false
	}

	endIdx := idx + len(pattern)
	if endIdx < len(timeStr) && isLetter(rune(timeStr[endIdx])) {
		return "", false
	}

	suffix := ""
	if endIdx < len(timeStr) {
		suffix = timeStr[endIdx:]
	}

	return fmt.Sprintf("%d%s%s", amount, to, suffix), true
}

// isLetter checks if a rune is a letter.
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// tryConvertShortForm attempts to convert short form units (d, y) to hours.
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

func parseRelativeTime(
	durationStr string,
	now, startTime time.Time,
	isPositive bool,
) (time.Time, error) {
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

// getMidnight returns the time at midnight (00:00:00) for the given date.
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

// parseUnixTimestamp parses a Unix timestamp, handling both seconds and milliseconds.
func parseUnixTimestamp(unixTime int64) time.Time {
	// Detect if it's milliseconds (13 digits or more) vs seconds (10 digits or less)
	// Timestamps > 9999999999 are either after year 2286 in seconds, or milliseconds
	if unixTime > timestampMillisecondBorder {
		// Treat as milliseconds
		sec := unixTime / millisecondsPerSecond
		nsec := (unixTime % millisecondsPerSecond) * nanosecondsPerMillisecond

		return time.Unix(sec, nsec)
	}

	// Treat as seconds (including negative timestamps for dates before 1970)
	return time.Unix(unixTime, 0)
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
