package friendlytime_test

import (
	"fmt"
	"time"

	friendlytime "github.com/Starlexxx/friendly-time"
)

// ExampleParseTimeRange demonstrates parsing various time range formats
func ExampleParseTimeRange() {
	start, end, err := friendlytime.ParseTimeRange("2h/1h")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Range from 2h ago to 1h ago: %d seconds\n", end-start)

	start2, end2, err := friendlytime.ParseTimeRange("1416434697/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("From timestamp to now: start=%d, end>%d\n", start2, end2)
}

// ExampleParseTimeRange_singleTime shows parsing a single time value
func ExampleParseTimeRange_singleTime() {
	start, end, err := friendlytime.ParseTimeRange("1416434697")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Single timestamp: %d (start==end: %v)\n", start, start == end)
}

// ExampleParseTimeRange_emptyRange shows behavior with empty values
func ExampleParseTimeRange_emptyRange() {
	start, end, err := friendlytime.ParseTimeRange("")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Empty range: start=%d, end=%d\n", start, end)
}

// ExampleParseTime_durations demonstrates parsing duration strings
func ExampleParseTime_durations() {
	now := time.Date(2025, 12, 10, 15, 30, 0, 0, time.UTC)

	t1, _ := friendlytime.ParseTime("1h", now, time.Time{})
	fmt.Printf("1h ago: %v\n", t1.Format("15:04"))

	t2, _ := friendlytime.ParseTime("30 minutes ago", now, time.Time{})
	fmt.Printf("30 minutes ago: %v\n", t2.Format("15:04"))

	t3, _ := friendlytime.ParseTime("2 days ago", now, time.Time{})
	fmt.Printf("2 days ago: %v\n", t3.Format("2006-01-02"))
}

// ExampleParseTime_timeOfDay shows parsing time of day
func ExampleParseTime_timeOfDay() {
	now := time.Date(2025, 12, 10, 15, 30, 0, 0, time.UTC)

	t, _ := friendlytime.ParseTime("09:30", now, time.Time{})
	fmt.Printf("Today at 09:30: %v\n", t.Format("2006-01-02 15:04"))
}

// ExampleParseTime_weekdays demonstrates parsing weekday references
func ExampleParseTime_weekdays() {
	now := time.Date(2025, 12, 10, 15, 30, 0, 0, time.UTC)

	t, _ := friendlytime.ParseTime("last monday", now, time.Time{})
	fmt.Printf("Last Monday: %v (weekday: %v)\n", t.Format("2006-01-02"), t.Weekday())
}

// ExampleParseTime_relativeKeywords shows relative time keywords
func ExampleParseTime_relativeKeywords() {
	now := time.Date(2025, 12, 10, 15, 30, 0, 0, time.UTC)

	t1, _ := friendlytime.ParseTime("yesterday", now, time.Time{})
	fmt.Printf("Yesterday: %v\n", t1.Format("2006-01-02"))

	t2, _ := friendlytime.ParseTime("last week", now, time.Time{})
	fmt.Printf("Last week: %v\n", t2.Format("2006-01-02"))
}

// ExampleParseTime_unixTimestamp demonstrates Unix timestamp parsing
func ExampleParseTime_unixTimestamp() {
	now := time.Now()

	t1, _ := friendlytime.ParseTime("1416434697", now, time.Time{})
	fmt.Printf("Unix timestamp (seconds): %v\n", t1.UTC().Format("2006-01-02 15:04:05"))

	t2, _ := friendlytime.ParseTime("1416434697000", now, time.Time{})
	fmt.Printf("Unix timestamp (milliseconds): %v\n", t2.UTC().Format("2006-01-02 15:04:05"))
}

// ExampleParseTime_relativeOffset shows using relative offsets from start time
func ExampleParseTime_relativeOffset() {
	now := time.Date(2025, 12, 10, 15, 30, 0, 0, time.UTC)
	startTime := time.Date(2025, 12, 10, 12, 0, 0, 0, time.UTC)

	t, _ := friendlytime.ParseTime("+1h", now, startTime)
	fmt.Printf("Start time + 1h: %v\n", t.Format("15:04"))
}

// ExampleParseTime_errorHandling demonstrates error handling
func ExampleParseTime_errorHandling() {
	now := time.Now()

	_, err := friendlytime.ParseTime("invalid format", now, time.Time{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// ExampleParseTimeRange_errorHandling shows error handling for ranges
func ExampleParseTimeRange_errorHandling() {
	// Invalid: end before start
	_, _, err := friendlytime.ParseTimeRange("1h/2h")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
