# Friendly Time

[![Go Reference](https://pkg.go.dev/badge/github.com/Starlexxx/friendly-time.svg)](https://pkg.go.dev/github.com/Starlexxx/friendly-time)
[![Go Report Card](https://goreportcard.com/badge/github.com/Starlexxx/friendly-time)](https://goreportcard.com/report/github.com/Starlexxx/friendly-time)
[![Coverage](https://img.shields.io/badge/coverage-96.2%25-brightgreen)](https://github.com/Starlexxx/friendly-time)

A Go library for parsing human-readable time expressions into Unix timestamps. Perfect for building time-based queries, filters, and analytics tools.

## Features

**Human-Readable Formats**: Parse natural language time expressions like "yesterday", "last monday", "2 hours ago"  
**Flexible**: Supports multiple date/time formats and timezones  
**Zero Dependencies**: Pure Go implementation

## Installation

```shell
go get github.com/Starlexxx/friendly-time
```

## Quick Start

```go
import friendlytime "github.com/Starlexxx/friendly-time"

// Parse a time range
start, end, err := friendlytime.ParseTimeRange("1h/30m")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("From %d to %d\n", start, end)

// Parse a single time
now := time.Now()
t, err := friendlytime.ParseTime("2 hours ago", now, time.Time{})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Time: %v\n", t)
```

## Supported Formats

### Durations (Relative to Now)

Standard Go durations:

- `1h`, `30m`, `45s` - hours, minutes, seconds
- `2h30m` - combined durations

Custom units with "ago":

- `10 seconds ago`, `1 second ago`
- `15 minutes ago`, `1 minute ago`
- `3 hours ago`
- `2 days ago`, `1 day ago`
- `1 month ago`, `2 months ago`
- `1 year ago`, `2 years ago`

Short forms:

- `5d` - 5 days ago
- `1y` - 1 year ago

### Relative Keywords

- `yesterday` - yesterday at 00:00:00
- `last week` - 7 days ago at 00:00:00
- `last month` - 1 month ago at 00:00:00
- `last year` - 1 year ago at 00:00:00

### Weekdays

- `last monday`, `last tuesday`, ..., `last sunday` - previous occurrence at 00:00:00

### Time of Day

- `00:00` - today at midnight
- `15:30` - today at 3:30 PM
- `09:00` - today at 9:00 AM

### Dates

- `2025-12-10` - specific date (YYYY-MM-DD)
- `14-11-19` - specific date (YY-MM-DD)
- `2025-12-10 15:30:45` - date with time
- `Mon, 02 Jan 2006 15:04:05` - RFC822 style

### Unix Timestamps

- `1416434697` - seconds since epoch
- `1416434697000` - milliseconds since epoch

### Relative Offsets

- `+30m` - 30 minutes after startTime parameter
- `-15m` - 15 minutes before now

### Time Ranges

Combine any of the above formats with `/`:

**Duration ranges:**

- `2h/1h` - from 2 hours ago to 1 hour ago
- `1d/12h` - from 1 day ago to 12 hours ago

**Mixed ranges:**

- `yesterday/12:00` - from yesterday midnight to today noon
- `last monday/yesterday` - from last Monday to yesterday
- `00:00/17:00` - from midnight to 5 PM today

**With timestamps:**

- `1416434697/` - from timestamp to now
- `/1416434697` - from zero time to timestamp

**Relative offsets:**

- `2h/+30m` - from 2h ago to 1h30m ago
- `09:00/+8h` - from 9 AM to 5 PM

## API Reference

### ParseTimeRange

```go
func ParseTimeRange(timeRange string) (start, end int64, err error)
```

Parses a human-readable time range string into Unix timestamps (seconds since epoch).

**Parameters:**

- `timeRange`: A time range string (e.g., "1h/30m", "yesterday/now")

**Returns:**

- `start`: Unix timestamp for the start of the range
- `end`: Unix timestamp for the end of the range
- `err`: Error if parsing fails

**Example:**

```go
start, end, err := friendlytime.ParseTimeRange("1h/30m")
// Returns timestamps for "1 hour ago" to "30 minutes ago"
```

### ParseTime

```go
func ParseTime(timeStr string, now time.Time, startTime time.Time) (time.Time, error)
```

Parses a human-readable time string into a `time.Time` value.

**Parameters:**

- `timeStr`: Time string to parse
- `now`: Reference time for relative calculations (usually `time.Now()`)
- `startTime`: Used for relative offsets with "+" prefix

**Returns:**

- `time.Time`: Parsed time value
- `error`: Error if format is not recognized

**Example:**

```go
now := time.Now()
t, err := friendlytime.ParseTime("2 hours ago", now, time.Time{})
```

## Error Types

The library defines several error types for better error handling:

```go
var (
    ErrInvalidTimeFormat  // Unrecognized time format
    ErrInvalidTimeRange   // Invalid range format
    ErrInvalidStartTime   // Start time couldn't be parsed
    ErrInvalidEndTime     // End time couldn't be parsed
    ErrInvalidWeekday     // Unrecognized weekday name
    ErrEndBeforeStart     // End time is before start time
)
```

## Advanced Usage

### Working with Timezones

```go
// Parse time in specific timezone
loc, _ := time.LoadLocation("America/New_York")
now := time.Now().In(loc)
t, _ := friendlytime.ParseTime("yesterday", now, time.Time{})
// Returns yesterday at midnight in New York timezone
```

### Using Relative Offsets

```go
now := time.Now()
start, _ := friendlytime.ParseTime("2h", now, time.Time{})
// start is 2 hours ago

// Use start as reference for end time
end, _ := friendlytime.ParseTime("+30m", now, start)
// end is 1h30m ago (start + 30 minutes)
```

### Batch Processing

```go
timeExpressions := []string{"1h", "yesterday", "last monday"}
for _, expr := range timeExpressions {
    t, err := friendlytime.ParseTime(expr, time.Now(), time.Time{})
    if err != nil {
        log.Printf("Failed to parse %s: %v", expr, err)
        continue
    }
    fmt.Printf("%s -> %v\n", expr, t)
}
```

## Testing

Run all tests:

```shell
go test -v
```

Run with coverage:

```shell
go test -cover
```

Run benchmarks:

```shell
go test -bench=. -benchmem
```

Run fuzz tests:

```shell
go test -fuzz=FuzzParseTime -fuzztime=30s
```

## License

MIT License - see [LICENSE](LICENSE) file for details
