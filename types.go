package tangerino

import (
	"fmt"
	"time"
)

// UnixMilliTime is an absolute timestamp stored as milliseconds since the Unix epoch,
// exactly as received from the API. It provides helpers to access the raw value or
// convert to standard Go types without losing precision.
type UnixMilliTime int64

// Raw returns the original Unix millisecond value as received from the API.
func (t UnixMilliTime) Raw() int64 {
	return int64(t)
}

// Time converts the value to a time.Time in UTC.
func (t UnixMilliTime) Time() time.Time {
	return time.UnixMilli(int64(t)).UTC()
}

// Format formats the timestamp using the given layout (same syntax as time.Time.Format).
//
// Example:
//
//	t.Format("02/01/2006")       // "01/01/2025"
//	t.Format("15:04")            // "09:00"
//	t.Format(time.RFC3339)
func (t UnixMilliTime) Format(layout string) string {
	return t.Time().Format(layout)
}

// String returns the timestamp formatted as "2006-01-02 15:04:05 UTC".
func (t UnixMilliTime) String() string {
	return t.Format("2006-01-02 15:04:05 UTC")
}

// DayOffset is a time offset from midnight stored as milliseconds, as received from
// the API. It is used for shift and interval fields in work schedule timetables.
// Values may exceed 86400000 (24 h) when a shift extends past midnight into the next day.
type DayOffset int64

// Raw returns the original millisecond value as received from the API.
func (d DayOffset) Raw() int64 {
	return int64(d)
}

// Duration converts the value to a time.Duration.
func (d DayOffset) Duration() time.Duration {
	return time.Duration(d) * time.Millisecond
}

// String formats the offset as "HH:MM".
// Hours are not capped at 23, so a shift ending at 28 h is shown as "28:00".
func (d DayOffset) String() string {
	totalSec := int64(d) / 1000
	h := totalSec / 3600
	m := (totalSec % 3600) / 60
	return fmt.Sprintf("%02d:%02d", h, m)
}
