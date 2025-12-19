package friendlytime

import (
	"errors"
	"testing"
	"time"
)

// FuzzParseTime tests ParseTime with random inputs to ensure it doesn't panic.
func FuzzParseTime(f *testing.F) {
	seeds := []string{
		"1h",
		"30m",
		"45s",
		"yesterday",
		"last monday",
		"15:30",
		"2025-12-10",
		"1416434697",
		"10 seconds ago",
		"2 days ago",
		"1y",
		"5d",
		"+30m",
		"-1h",
		"",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	now := time.Date(2025, 12, 10, 15, 30, 0, 0, time.UTC)
	startTime := time.Date(2025, 12, 10, 12, 0, 0, 0, time.UTC)

	f.Fuzz(func(t *testing.T, input string) {
		result, err := ParseTime(input, now, startTime)
		if err == nil {
			_ = result
		} else if !errors.Is(err, ErrInvalidTimeFormat) &&
			!errors.Is(err, ErrInvalidWeekday) {
			t.Errorf("Expected non-nil error for input: %q", input)
		}
	})
}

// FuzzParseTimeRange tests ParseTimeRange with random inputs.
func FuzzParseTimeRange(f *testing.F) {
	seeds := []string{
		"1h/30m",
		"2h/1h",
		"yesterday/now",
		"last monday/last friday",
		"09:00/17:00",
		"1416434697/1416434797",
		"/now",
		"yesterday/",
		"",
		"1h",
		"2d/+1d",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		start, end, err := ParseTimeRange(input)
		if err == nil {
			validateTimeRangeResult(t, input, start, end)
		} else {
			validateTimeRangeError(t, input, err)
		}
	})
}

// validateTimeRangeResult validates successful ParseTimeRange results.
func validateTimeRangeResult(t *testing.T, input string, start, end int64) {
	t.Helper()

	if input == "" && (start != 0 || end != 0) {
		t.Errorf("Empty input should give 0,0, got %d,%d", start, end)
	}

	if start != 0 && end != 0 && end < start {
		t.Errorf(
			"End (%d) should not be before start (%d) for input: %q",
			end,
			start,
			input,
		)
	}
}

// validateTimeRangeError validates ParseTimeRange errors.
func validateTimeRangeError(t *testing.T, input string, err error) {
	t.Helper()

	if !isExpectedTimeRangeError(err) && err != nil {
		t.Logf("Unexpected error for input %q: %v", input, err)
	}
}

// isExpectedTimeRangeError checks if error is one of the expected types.
func isExpectedTimeRangeError(err error) bool {
	return errors.Is(err, ErrInvalidTimeRange) ||
		errors.Is(err, ErrInvalidStartTime) ||
		errors.Is(err, ErrInvalidEndTime) ||
		errors.Is(err, ErrEndBeforeStart)
}

// FuzzConvertCustomUnits tests the unit conversion function.
func FuzzConvertCustomUnits(f *testing.F) {
	seeds := []string{
		"30 seconds",
		"15 minutes",
		"2 hours",
		"3 days",
		"1 month",
		"2 years",
		"5d",
		"1y",
		"last monday",
		"yesterday",
		"",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		result := convertCustomUnits(input)

		if input != "" && result == "" {
			t.Errorf("Result should not be empty for non-empty input: %q", input)
		}
	})
}
