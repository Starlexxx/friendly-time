package friendlytime

import "errors"

// Error types for time parsing operations.
var (
	// ErrInvalidTimeFormat indicates the time string format is not recognized.
	ErrInvalidTimeFormat = errors.New("invalid time format")

	// ErrInvalidTimeRange indicates the time range format is invalid (e.g., wrong number of separators).
	ErrInvalidTimeRange = errors.New("invalid time range format")

	// ErrInvalidStartTime indicates the start time in a range could not be parsed.
	ErrInvalidStartTime = errors.New("invalid start time")

	// ErrInvalidEndTime indicates the end time in a range could not be parsed.
	ErrInvalidEndTime = errors.New("invalid end time")

	// ErrInvalidWeekday indicates an unrecognized weekday name was provided.
	ErrInvalidWeekday = errors.New("invalid weekday")

	// ErrEndBeforeStart indicates the end time is chronologically before the start time.
	ErrEndBeforeStart = errors.New("end time is before start time")
)
