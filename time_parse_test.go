package friendlytime

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a fixed time for testing.
func fixedTime() time.Time {
	// Wednesday, December 10, 2025, 15:30:45
	return time.Date(2025, 12, 10, 15, 30, 45, 0, time.UTC)
}

// Helper function to create midnight for a given date.
func midnight(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func TestParseTime_Durations(t *testing.T) {
	now := fixedTime()

	tests := []struct {
		name     string
		input    string
		expected time.Time
		wantErr  bool
	}{
		// Standard Go durations
		{
			name:     "1 hour ago",
			input:    "1h",
			expected: now.Add(-1 * time.Hour),
		},
		{
			name:     "30 minutes ago",
			input:    "30m",
			expected: now.Add(-30 * time.Minute),
		},
		{
			name:     "45 seconds ago",
			input:    "45s",
			expected: now.Add(-45 * time.Second),
		},
		{
			name:     "2 hours 30 minutes ago",
			input:    "2h30m",
			expected: now.Add(-2*time.Hour - 30*time.Minute),
		},

		// Custom units - seconds
		{
			name:     "10 seconds ago",
			input:    "10 seconds ago",
			expected: now.Add(-10 * time.Second),
		},
		{
			name:     "1 second ago",
			input:    "1 second ago",
			expected: now.Add(-1 * time.Second),
		},
		{
			name:     "5 sec ago",
			input:    "5 sec ago",
			expected: now.Add(-5 * time.Second),
		},

		// Custom units - minutes
		{
			name:     "15 minutes ago",
			input:    "15 minutes ago",
			expected: now.Add(-15 * time.Minute),
		},
		{
			name:     "1 minute ago",
			input:    "1 minute ago",
			expected: now.Add(-1 * time.Minute),
		},
		{
			name:     "20 min ago",
			input:    "20 min ago",
			expected: now.Add(-20 * time.Minute),
		},

		// Custom units - hours
		{
			name:     "3 hours ago",
			input:    "3 hours ago",
			expected: now.Add(-3 * time.Hour),
		},
		{
			name:     "1 hour ago",
			input:    "1 hour ago",
			expected: now.Add(-1 * time.Hour),
		},

		// Custom units - days
		{
			name:     "2 days ago",
			input:    "2 days ago",
			expected: now.Add(-2 * 24 * time.Hour),
		},
		{
			name:     "1 day ago",
			input:    "1 day ago",
			expected: now.Add(-1 * 24 * time.Hour),
		},
		{
			name:     "5d",
			input:    "5d",
			expected: now.Add(-5 * 24 * time.Hour),
		},

		// Custom units - months (approximated as 30 days = 720 hours)
		{
			name:     "1 month ago",
			input:    "1 month ago",
			expected: now.Add(-720 * time.Hour),
		},
		{
			name:     "2 months ago",
			input:    "2 months ago",
			expected: now.Add(-1440 * time.Hour),
		},

		// Custom units - years (approximated as 365 days = 8760 hours)
		{
			name:     "1 year ago",
			input:    "1 year ago",
			expected: now.Add(-8760 * time.Hour),
		},
		{
			name:     "2 years ago",
			input:    "2 years ago",
			expected: now.Add(-17520 * time.Hour),
		},
		{
			name:     "1y",
			input:    "1y",
			expected: now.Add(-8760 * time.Hour),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, now, time.Time{})
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected.Unix(), result.Unix(), "Time mismatch for input: %s", tt.input)
			}
		})
	}
}

func TestParseTime_RelativeKeywords(t *testing.T) {
	// Wednesday, December 10, 2025, 15:30:45
	now := fixedTime()

	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "yesterday",
			input:    "yesterday",
			expected: midnight(2025, 12, 9), // Tuesday at 00:00:00
		},
		{
			name:     "last week",
			input:    "last week",
			expected: midnight(2025, 12, 3), // Wednesday at 00:00:00
		},
		{
			name:     "last month",
			input:    "last month",
			expected: midnight(2025, 11, 10), // November 10 at 00:00:00
		},
		{
			name:     "last year",
			input:    "last year",
			expected: midnight(2024, 12, 10), // December 10, 2024 at 00:00:00
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, now, time.Time{})
			require.NoError(t, err)
			assert.Equal(
				t,
				tt.expected.Unix(),
				result.Unix(),
				"Expected %v, got %v",
				tt.expected,
				result,
			)
		})
	}
}

func TestParseTime_Weekdays(t *testing.T) {
	// Wednesday, December 10, 2025, 15:30:45
	now := fixedTime()

	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "last monday",
			input:    "last monday",
			expected: midnight(2025, 12, 8), // Monday, Dec 8 at 00:00:00
		},
		{
			name:     "last tuesday",
			input:    "last tuesday",
			expected: midnight(2025, 12, 9), // Tuesday, Dec 9 at 00:00:00
		},
		{
			name:     "last wednesday",
			input:    "last wednesday",
			expected: midnight(2025, 12, 3), // Wednesday, Dec 3 at 00:00:00 (last week)
		},
		{
			name:     "last thursday",
			input:    "last thursday",
			expected: midnight(2025, 12, 4), // Thursday, Dec 4 at 00:00:00
		},
		{
			name:     "last friday",
			input:    "last friday",
			expected: midnight(2025, 12, 5), // Friday, Dec 5 at 00:00:00
		},
		{
			name:     "last saturday",
			input:    "last saturday",
			expected: midnight(2025, 12, 6), // Saturday, Dec 6 at 00:00:00
		},
		{
			name:     "last sunday",
			input:    "last sunday",
			expected: midnight(2025, 12, 7), // Sunday, Dec 7 at 00:00:00
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, now, time.Time{})
			require.NoError(t, err)
			assert.Equal(
				t,
				tt.expected.Unix(),
				result.Unix(),
				"Expected %v, got %v",
				tt.expected,
				result,
			)
		})
	}
}

func TestParseTime_TimeOfDay(t *testing.T) {
	now := fixedTime()

	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "midnight",
			input:    "00:00",
			expected: time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "noon",
			input:    "12:00",
			expected: time.Date(2025, 12, 10, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "specific time",
			input:    "09:30",
			expected: time.Date(2025, 12, 10, 9, 30, 0, 0, time.UTC),
		},
		{
			name:     "evening time",
			input:    "18:45",
			expected: time.Date(2025, 12, 10, 18, 45, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, now, time.Time{})
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Unix(), result.Unix())
		})
	}
}

func TestParseTime_Dates(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "short date format YY-MM-DD",
			input:    "14-11-19",
			expected: time.Date(2014, 11, 19, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "full date format YYYY-MM-DD",
			input:    "2025-12-10",
			expected: time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "short datetime YY-MM-DD HH:MM:SS",
			input:    "14-11-19 22:04:57",
			expected: time.Date(2014, 11, 19, 22, 4, 57, 0, time.UTC),
		},
		{
			name:     "full datetime YYYY-MM-DD HH:MM:SS",
			input:    "2025-12-10 15:30:45",
			expected: time.Date(2025, 12, 10, 15, 30, 45, 0, time.UTC),
		},
		{
			name:     "RFC822 style",
			input:    "Mon, 07 Aug 2006 12:34:56",
			expected: time.Date(2006, 8, 7, 12, 34, 56, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, time.Now(), time.Time{})
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Unix(), result.Unix())
		})
	}
}

func TestParseTime_UnixTimestamps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64 // Unix timestamp in seconds
	}{
		{
			name:     "unix timestamp in seconds",
			input:    "1416434697",
			expected: 1416434697,
		},
		{
			name:     "unix timestamp in milliseconds",
			input:    "1416434697000",
			expected: 1416434697,
		},
		{
			name:     "zero timestamp",
			input:    "0",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, time.Now(), time.Time{})
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.Unix())
		})
	}
}

func TestParseTime_RelativeToStartTime(t *testing.T) {
	now := fixedTime()
	startTime := now.Add(-2 * time.Hour)

	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "positive offset from start time",
			input:    "+30m",
			expected: startTime.Add(30 * time.Minute),
		},
		{
			name:     "positive offset 1 hour from start time",
			input:    "+1h",
			expected: startTime.Add(1 * time.Hour),
		},
		{
			name:     "negative offset from now",
			input:    "-15m",
			expected: now.Add(-15 * time.Minute),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, now, startTime)
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Unix(), result.Unix())
		})
	}
}

func TestParseTime_EmptyString(t *testing.T) {
	now := fixedTime()

	t.Run("empty string with zero start time", func(t *testing.T) {
		result, err := ParseTime("", now, time.Time{})
		require.NoError(t, err)
		assert.True(t, result.IsZero())
	})

	t.Run("empty string with non-zero start time", func(t *testing.T) {
		result, err := ParseTime("", now, now)
		require.NoError(t, err)
		assert.Equal(t, now.Unix(), result.Unix())
	})
}

func TestParseTime_InvalidFormats(t *testing.T) {
	now := fixedTime()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid text",
			input: "invalid",
		},
		{
			name:  "random string",
			input: "foobar",
		},
		{
			name:  "malformed date",
			input: "2025-13-45",
		},
		{
			name:  "invalid weekday",
			input: "last funday",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTime(tt.input, now, time.Time{})
			require.Error(t, err)
		})
	}
}

func TestParseTimeRange_SingleTime(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantStart int64
		wantEnd   int64
	}{
		{
			name:      "single timestamp",
			input:     "1416434697",
			wantStart: 1416434697,
			wantEnd:   1416434697,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := ParseTimeRange(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.wantStart, start)
			assert.Equal(t, tt.wantEnd, end)
		})
	}

	t.Run("single duration", func(t *testing.T) {
		before := time.Now().Add(-1 * time.Hour).Unix()
		start, end, err := ParseTimeRange("1h")
		after := time.Now().Add(-1 * time.Hour).Unix()

		require.NoError(t, err)
		assert.Equal(t, start, end, "start and end should be equal")

		// Allow for time passage during test execution
		assert.GreaterOrEqual(t, start, before-1, "start should be approximately 1h ago")
		assert.LessOrEqual(t, start, after+1, "start should be approximately 1h ago")
	})
}

func TestParseTimeRange_Ranges(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "duration range",
			input:       "2h/1h",
			description: "from 2 hours ago to 1 hour ago",
		},
		{
			name:        "time of day range",
			input:       "09:00/17:00",
			description: "from 9am to 5pm today",
		},
		{
			name:        "mixed range",
			input:       "yesterday/12:00",
			description: "from yesterday midnight to today noon",
		},
		{
			name:        "relative positive offset",
			input:       "2h/+30m",
			description: "from 2 hours ago to 1.5 hours ago",
		},
		{
			name:        "timestamp to duration",
			input:       "1416434697/1h",
			description: "from specific timestamp to 1 hour ago",
		},
		{
			name:        "empty start to timestamp",
			input:       "/1416434697",
			description: "from beginning to specific timestamp",
		},
		{
			name:        "timestamp to empty end",
			input:       "1416434697/",
			description: "from specific timestamp to now",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := ParseTimeRange(tt.input)
			require.NoError(t, err, "Failed to parse: %s", tt.description)

			// Basic sanity checks
			if tt.input != "/1416434697" { // Skip validation for ranges starting from zero
				assert.LessOrEqual(
					t,
					start,
					end,
					"Start should be before or equal to end: %s",
					tt.description,
				)
			}
		})
	}
}

func TestParseTimeRange_ComplexScenarios(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errType error
	}{
		{
			name:    "valid range with keywords",
			input:   "last monday/yesterday",
			wantErr: false,
		},
		{
			name:    "valid range with mixed formats",
			input:   "2025-12-01/2025-12-10",
			wantErr: false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: false,
		},
		{
			name:    "invalid format - too many slashes",
			input:   "1h/2h/3h",
			wantErr: true,
			errType: ErrInvalidTimeRange,
		},
		{
			name:    "invalid start time",
			input:   "invalid/1h",
			wantErr: true,
			errType: ErrInvalidStartTime,
		},
		{
			name:    "invalid end time",
			input:   "1h/invalid",
			wantErr: true,
			errType: ErrInvalidEndTime,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := ParseTimeRange(tt.input)
			if tt.wantErr {
				require.Error(t, err)

				if tt.errType != nil {
					assert.True(
						t,
						errors.Is(err, tt.errType),
						"Expected error type %v, got %v",
						tt.errType,
						err,
					)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestParseTimeRange_EndBeforeStart(t *testing.T) {
	_, _, err := ParseTimeRange("1h/2h")
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrEndBeforeStart), "Expected ErrEndBeforeStart")
}

func TestConvertCustomUnits(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "seconds to s",
			input:    "30 seconds",
			expected: "30s",
		},
		{
			name:     "minutes to m",
			input:    "15 minutes",
			expected: "15m",
		},
		{
			name:     "hours to h",
			input:    "2 hours",
			expected: "2h",
		},
		{
			name:     "days to hours",
			input:    "3 days",
			expected: "72h",
		},
		{
			name:     "1 day to hours",
			input:    "1 day",
			expected: "24h",
		},
		{
			name:     "months to hours",
			input:    "2 months",
			expected: "1440h",
		},
		{
			name:     "years to hours",
			input:    "1 year",
			expected: "8760h",
		},
		{
			name:     "preserve weekdays",
			input:    "last monday",
			expected: "last monday",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertCustomUnits(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMidnight(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected time.Time
	}{
		{
			name:     "afternoon to midnight",
			input:    time.Date(2025, 12, 10, 15, 30, 45, 0, time.UTC),
			expected: time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "morning to midnight",
			input:    time.Date(2025, 12, 10, 8, 15, 30, 0, time.UTC),
			expected: time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "already midnight",
			input:    time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2025, 12, 10, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getMidnight(tt.input)
			assert.Equal(t, tt.expected.Unix(), result.Unix())
		})
	}
}

func TestParseWeekday(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Weekday
		wantErr  bool
	}{
		{name: "sunday", input: "sunday", expected: time.Sunday},
		{name: "monday", input: "monday", expected: time.Monday},
		{name: "tuesday", input: "tuesday", expected: time.Tuesday},
		{name: "wednesday", input: "wednesday", expected: time.Wednesday},
		{name: "thursday", input: "thursday", expected: time.Thursday},
		{name: "friday", input: "friday", expected: time.Friday},
		{name: "saturday", input: "saturday", expected: time.Saturday},
		{name: "invalid", input: "notaday", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseWeekday(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, errors.Is(err, ErrInvalidWeekday))
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Edge cases and regression tests.
func TestParseTime_EdgeCases(t *testing.T) {
	now := fixedTime()

	t.Run("very large duration", func(t *testing.T) {
		result, err := ParseTime("1000h", now, time.Time{})
		require.NoError(t, err)

		expected := now.Add(-1000 * time.Hour)
		assert.Equal(t, expected.Unix(), result.Unix())
	})

	t.Run("zero duration", func(t *testing.T) {
		result, err := ParseTime("0s", now, time.Time{})
		require.NoError(t, err)
		assert.Equal(t, now.Unix(), result.Unix())
	})

	t.Run("multiple spaces in 'ago' format", func(t *testing.T) {
		result, err := ParseTime("5  minutes  ago", now, time.Time{})
		require.NoError(t, err)

		expected := now.Add(-5 * time.Minute)
		assert.Equal(t, expected.Unix(), result.Unix())
	})
}

func TestParseTimeRange_EdgeCases(t *testing.T) {
	t.Run("same start and end time", func(t *testing.T) {
		start, end, err := ParseTimeRange("2025-12-10 12:00:00/2025-12-10 12:00:00")
		require.NoError(t, err)
		assert.Equal(t, start, end)
	})

	t.Run("very wide range", func(t *testing.T) {
		start, end, err := ParseTimeRange("2020-01-01/2025-12-31")
		require.NoError(t, err)
		assert.Less(t, start, end)

		diff := end - start
		expectedDiff := int64(6 * 365 * 24 * 60 * 60)
		assert.InDelta(t, expectedDiff, diff, float64(2*24*60*60))
	})
}

// Test thread safety (if library is used concurrently).
func TestParseTime_Concurrent(t *testing.T) {
	now := fixedTime()

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_, err := ParseTime("1h", now, time.Time{})
				assert.NoError(t, err)
			}

			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// Test formats without spaces.
func TestParseTime_NoSpaceFormats(t *testing.T) {
	now := fixedTime()

	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "3days without space",
			input:    "3days",
			expected: now.Add(-3 * 24 * time.Hour),
		},
		{
			name:     "2months without space",
			input:    "2months",
			expected: now.Add(-1440 * time.Hour), // 2 * 30 days * 24 hours
		},
		{
			name:     "1year without space",
			input:    "1year",
			expected: now.Add(-8760 * time.Hour), // 365 days * 24 hours
		},
		{
			name:     "10seconds without space",
			input:    "10seconds",
			expected: now.Add(-10 * time.Second),
		},
		{
			name:     "5minutes without space",
			input:    "5minutes",
			expected: now.Add(-5 * time.Minute),
		},
		{
			name:     "2hours without space",
			input:    "2hours",
			expected: now.Add(-2 * time.Hour),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, now, time.Time{})
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Unix(), result.Unix())
		})
	}
}

// Test millisecond timestamps.
func TestParseTime_MillisecondTimestamps(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "millisecond timestamp",
			input:    "1416434697000",
			expected: time.Unix(1416434697, 0),
		},
		{
			name:     "millisecond timestamp with partial milliseconds",
			input:    "1416434697123",
			expected: time.Unix(1416434697, 123000000),
		},
		{
			name:     "second timestamp (should not be treated as milliseconds)",
			input:    "1416434697",
			expected: time.Unix(1416434697, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, time.Now(), time.Time{})
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Unix(), result.Unix())
			// For millisecond precision, also check nanoseconds
			if len(tt.input) > 10 {
				assert.Equal(t, tt.expected.Nanosecond(), result.Nanosecond())
			}
		})
	}
}

// Test timezone handling.
func TestParseTime_Timezones(t *testing.T) {
	utc := time.Date(2025, 12, 10, 15, 30, 45, 0, time.UTC)
	est := time.Date(2025, 12, 10, 15, 30, 45, 0, time.FixedZone("EST", -5*3600))
	pst := time.Date(2025, 12, 10, 15, 30, 45, 0, time.FixedZone("PST", -8*3600))

	tests := []struct {
		name string
		now  time.Time
		desc string
	}{
		{
			name: "UTC timezone",
			now:  utc,
			desc: "should work correctly with UTC",
		},
		{
			name: "EST timezone",
			now:  est,
			desc: "should work correctly with EST",
		},
		{
			name: "PST timezone",
			now:  pst,
			desc: "should work correctly with PST",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime("1h", tt.now, time.Time{})
			require.NoError(t, err)

			expected := tt.now.Add(-1 * time.Hour)
			assert.Equal(t, expected.Unix(), result.Unix(), tt.desc)

			result2, err := ParseTime("09:00", tt.now, time.Time{})
			require.NoError(t, err)
			assert.Equal(t, tt.now.Location(), result2.Location(), "timezone should be preserved")
		})
	}
}

// Test boundary values.
func TestParseTime_BoundaryValues(t *testing.T) {
	now := fixedTime()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "very large year value",
			input:   "100 years ago",
			wantErr: false,
		},
		{
			name:    "very small duration",
			input:   "1s",
			wantErr: false,
		},
		{
			name:    "zero value",
			input:   "0s",
			wantErr: false,
		},
		{
			name:    "negative unix timestamp",
			input:   "-1",
			wantErr: false, // Should be treated as negative timestamp (before 1970)
		},
		{
			name:    "very large timestamp",
			input:   "9999999999",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTime(tt.input, now, time.Time{})
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// Test case sensitivity.
func TestParseTime_CaseSensitivity(t *testing.T) {
	now := fixedTime()

	tests := []struct {
		name  string
		input string
	}{
		{name: "uppercase YESTERDAY", input: "YESTERDAY"},
		{name: "mixed case YeStErDaY", input: "YeStErDaY"},
		{name: "uppercase LAST MONDAY", input: "LAST MONDAY"},
		{name: "mixed case Last Monday", input: "Last Monday"},
		{name: "uppercase LAST WEEK", input: "LAST WEEK"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, now, time.Time{})
			require.NoError(t, err, "Should handle case-insensitive input: %s", tt.input)
			assert.False(t, result.IsZero(), "Result should not be zero time")
		})
	}
}

// Test multiple spaces and whitespace variations.
func TestParseTime_WhitespaceVariations(t *testing.T) {
	now := fixedTime()

	tests := []struct {
		name     string
		input    string
		expected time.Time
	}{
		{
			name:     "multiple spaces",
			input:    "5    minutes    ago",
			expected: now.Add(-5 * time.Minute),
		},
		{
			name:     "leading spaces",
			input:    "  10 seconds ago",
			expected: now.Add(-10 * time.Second),
		},
		{
			name:     "trailing spaces",
			input:    "10 seconds ago  ",
			expected: now.Add(-10 * time.Second),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseTime(tt.input, now, time.Time{})
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Unix(), result.Unix())
		})
	}
}

// Test all weekdays systematically.
func TestParseTime_AllWeekdaysSystematic(t *testing.T) {
	testDates := []struct {
		name string
		date time.Time
	}{
		{"From Sunday", time.Date(2025, 12, 7, 15, 30, 45, 0, time.UTC)},     // Sunday
		{"From Monday", time.Date(2025, 12, 8, 15, 30, 45, 0, time.UTC)},     // Monday
		{"From Tuesday", time.Date(2025, 12, 9, 15, 30, 45, 0, time.UTC)},    // Tuesday
		{"From Wednesday", time.Date(2025, 12, 10, 15, 30, 45, 0, time.UTC)}, // Wednesday
		{"From Thursday", time.Date(2025, 12, 11, 15, 30, 45, 0, time.UTC)},  // Thursday
		{"From Friday", time.Date(2025, 12, 12, 15, 30, 45, 0, time.UTC)},    // Friday
		{"From Saturday", time.Date(2025, 12, 13, 15, 30, 45, 0, time.UTC)},  // Saturday
	}

	weekdays := []string{
		"sunday",
		"monday",
		"tuesday",
		"wednesday",
		"thursday",
		"friday",
		"saturday",
	}

	for _, td := range testDates {
		t.Run(td.name, func(t *testing.T) {
			for _, wd := range weekdays {
				result, err := ParseTime("last "+wd, td.date, time.Time{})
				require.NoError(t, err, "Failed to parse 'last %s' from %s", wd, td.name)

				assert.True(t, result.Before(td.date), "last %s should be before %s", wd, td.name)

				assert.Equal(t, 0, result.Hour())
				assert.Equal(t, 0, result.Minute())
				assert.Equal(t, 0, result.Second())
			}
		})
	}
}

// Test ParseTimeRange with all edge cases.
func TestParseTimeRange_MoreEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		errType     error
		description string
	}{
		{
			name:        "both empty",
			input:       "/",
			wantErr:     false,
			description: "empty start and end should be valid",
		},
		{
			name:        "millisecond timestamps",
			input:       "1416434697000/1416434797000",
			wantErr:     false,
			description: "should handle millisecond timestamps",
		},
		{
			name:        "mixed second and millisecond",
			input:       "1416434697/1416434797000",
			wantErr:     false,
			description: "should handle mixed timestamp formats",
		},
		{
			name:        "relative with weekday",
			input:       "last monday/last friday",
			wantErr:     true,
			errType:     ErrEndBeforeStart,
			description: "last monday to last friday should fail (friday is before monday when run on Wednesday)",
		},
		{
			name:        "positive offset without start time",
			input:       "+1h",
			wantErr:     false,
			description: "positive offset without range should work from now",
		},
		{
			name:        "negative offset",
			input:       "-2h",
			wantErr:     false,
			description: "negative offset should work",
		},
		{
			name:        "complex time of day range",
			input:       "00:00/23:59",
			wantErr:     false,
			description: "full day range should work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end, err := ParseTimeRange(tt.input)
			if tt.wantErr {
				require.Error(t, err, tt.description)

				if tt.errType != nil {
					assert.True(
						t,
						errors.Is(err, tt.errType),
						"Expected error type %v, got %v",
						tt.errType,
						err,
					)
				}
			} else {
				require.NoError(t, err, tt.description)
				t.Logf("Start: %d, End: %d, Diff: %d seconds", start, end, end-start)
			}
		})
	}
}

// Test date parsing edge cases.
func TestParseTime_DateEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "leap year date",
			input:   "2024-02-29",
			wantErr: false,
		},
		{
			name:    "invalid leap year date",
			input:   "2023-02-29",
			wantErr: true,
		},
		{
			name:    "end of year",
			input:   "2025-12-31",
			wantErr: false,
		},
		{
			name:    "start of year",
			input:   "2025-01-01",
			wantErr: false,
		},
		{
			name:    "invalid month",
			input:   "2025-13-01",
			wantErr: true,
		},
		{
			name:    "invalid day",
			input:   "2025-11-31",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTime(tt.input, time.Now(), time.Time{})
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
