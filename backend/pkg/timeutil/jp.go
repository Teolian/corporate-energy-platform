// Package timeutil provides timezone utilities for Japanese power system data.
// Per AGENT_TECH_SPEC.md: Japan has no DST; always use Asia/Tokyo (UTC+09:00).
package timeutil

import (
	"fmt"
	"time"
)

var (
	// TokyoLocation is the IANA timezone for Japan.
	TokyoLocation *time.Location
)

func init() {
	var err error
	TokyoLocation, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(fmt.Sprintf("failed to load Asia/Tokyo timezone: %v", err))
	}
}

// ParseDate parses YYYY-MM-DD string into time.Time in Asia/Tokyo at 00:00:00.
func ParseDate(dateStr string) (time.Time, error) {
	t, err := time.ParseInLocation("2006-01-02", dateStr, TokyoLocation)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}
	return t, nil
}

// FormatDate formats time.Time to YYYY-MM-DD string.
func FormatDate(t time.Time) string {
	return t.In(TokyoLocation).Format("2006-01-02")
}

// HourlySlots generates 24 hourly timestamps for a given date in Asia/Tokyo.
// Returns timestamps from 00:00 to 23:00 (inclusive).
func HourlySlots(date time.Time) []time.Time {
	base := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, TokyoLocation)
	slots := make([]time.Time, 24)
	for i := 0; i < 24; i++ {
		slots[i] = base.Add(time.Duration(i) * time.Hour)
	}
	return slots
}

// FormatISO8601 formats time to ISO8601 with +09:00 offset (required by spec).
func FormatISO8601(t time.Time) string {
	return t.In(TokyoLocation).Format("2006-01-02T15:04:05+09:00")
}
