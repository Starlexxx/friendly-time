package friendlytime

import (
	"testing"
	"time"
)

// Benchmark for ParseTime with different input types.
func BenchmarkParseTime_Duration(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("1h", now, time.Time{})
	}
}

func BenchmarkParseTime_CustomUnit(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("30 minutes ago", now, time.Time{})
	}
}

func BenchmarkParseTime_Weekday(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("last monday", now, time.Time{})
	}
}

func BenchmarkParseTime_Yesterday(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("yesterday", now, time.Time{})
	}
}

func BenchmarkParseTime_TimeOfDay(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("15:30", now, time.Time{})
	}
}

func BenchmarkParseTime_Date(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("2025-12-10", now, time.Time{})
	}
}

func BenchmarkParseTime_DateTime(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("2025-12-10 15:30:45", now, time.Time{})
	}
}

func BenchmarkParseTime_UnixTimestamp(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("1416434697", now, time.Time{})
	}
}

func BenchmarkParseTime_RelativePositive(b *testing.B) {
	now := time.Now()
	startTime := now.Add(-2 * time.Hour)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("+30m", now, startTime)
	}
}

func BenchmarkParseTime_Days(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("5 days ago", now, time.Time{})
	}
}

func BenchmarkParseTime_Months(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("2 months ago", now, time.Time{})
	}
}

func BenchmarkParseTime_Years(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = ParseTime("1 year ago", now, time.Time{})
	}
}

// Benchmark for ParseTimeRange.
func BenchmarkParseTimeRange_Simple(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = ParseTimeRange("1h")
	}
}

func BenchmarkParseTimeRange_Range(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = ParseTimeRange("2h/1h")
	}
}

func BenchmarkParseTimeRange_ComplexRange(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = ParseTimeRange("last monday/yesterday")
	}
}

func BenchmarkParseTimeRange_MixedFormats(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = ParseTimeRange("2025-12-01 09:00:00/30 minutes ago")
	}
}

func BenchmarkParseTimeRange_TimeOfDayRange(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = ParseTimeRange("09:00/17:00")
	}
}

func BenchmarkParseTimeRange_RelativeOffset(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _, _ = ParseTimeRange("2h/+30m")
	}
}

// Benchmark for convertCustomUnits.
func BenchmarkConvertCustomUnits_Seconds(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = convertCustomUnits("30 seconds")
	}
}

func BenchmarkConvertCustomUnits_Days(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = convertCustomUnits("5 days")
	}
}

func BenchmarkConvertCustomUnits_Months(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = convertCustomUnits("2 months")
	}
}

func BenchmarkConvertCustomUnits_Years(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = convertCustomUnits("1 year")
	}
}

func BenchmarkConvertCustomUnits_Weekday(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = convertCustomUnits("last monday")
	}
}

// Benchmark for parseWeekday.
func BenchmarkParseWeekday(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = parseWeekday("monday")
	}
}

// Benchmark for getMidnight.
func BenchmarkGetMidnight(b *testing.B) {
	now := time.Now()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = getMidnight(now)
	}
}
