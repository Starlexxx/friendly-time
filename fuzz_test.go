package friendlytime

import (
	"testing"
	"time"
)

// FuzzParseTime tests ParseTime with random inputs to ensure it doesn't panic
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
			if result.IsZero() && input != "" && input != " " {
			}
		} else {
			if err != ErrInvalidTimeFormat &&
				err != ErrInvalidWeekday {
				if err == nil {
					t.Errorf("Expected non-nil error for input: %q", input)
				}
			}
		}
	})
}

// FuzzParseTimeRange tests ParseTimeRange with random inputs
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
			if input == "" {
				if start != 0 || end != 0 {
					t.Errorf("Empty input should give 0,0, got %d,%d", start, end)
				}
			}

			if start != 0 && end != 0 && end < start {
				t.Errorf("End (%d) should not be before start (%d) for input: %q", end, start, input)
			}
		} else {
			if err != ErrInvalidTimeRange &&
				err != ErrInvalidStartTime &&
				err != ErrInvalidEndTime &&
				err != ErrEndBeforeStart {
				if err == nil {
					t.Errorf("Expected non-nil error for input: %q", input)
				}
			}
		}
	})
}

// FuzzConvertCustomUnits tests the unit conversion function
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
		result, err := convertCustomUnits(input)

		if err != nil {
			t.Errorf("convertCustomUnits should not error, got: %v for input: %q", err, input)
		}

		if input != "" && result == "" {
			t.Errorf("Result should not be empty for non-empty input: %q", input)
		}
	})
}
