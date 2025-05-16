package main

// Example showing how to handle both old and new Time format
//
// Note: Starting with v0.1.0, you can control the format by setting the
// environment variable LUNO_TIME_LEGACY_FORMAT=true to use the legacy string format.
// This example shows how to handle both formats if you're consuming data from
// different versions or configurations of the library.

import (
	"encoding/json"
	"fmt"
	"time"
)

// TimeCompatHelper helps with migrating from old Time string format to new millisecond format
func TimeCompatHelper(timeField json.RawMessage) (time.Time, error) {
	// Try to parse as millisecond timestamp (new format)
	var ms int64
	err := json.Unmarshal(timeField, &ms)
	if err == nil {
		return time.Unix(0, ms*1e6), nil
	}

	// Try to parse as string (old format)
	var timeStr string
	err = json.Unmarshal(timeField, &timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse time: %w", err)
	}

	// Try parsing as time string
	t, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return t, nil
	}

	// Try various other time formats...
	layouts := []string{
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		"2006-01-02 15:04:05 -0700 MST",
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, timeStr)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("could not parse time from: %s", timeStr)
}
