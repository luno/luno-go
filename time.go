package luno

import (
	"encoding/json"
	"strconv"
	"time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) error {
	// Handle null case
	if string(b) == "null" || len(b) == 0 {
		*t = Time{}
		return nil
	}

	// Try to parse as integer (milliseconds since epoch)
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err == nil {
		if i == 0 {
			*t = Time{}
			return nil
		}
		*t = Time(time.Unix(0, i*1e6))
		return nil
	}

	// Try to parse as string
	var timeStr string
	if err := json.Unmarshal(b, &timeStr); err != nil {
		return err
	}

	// Try RFC3339 format
	parsed, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		*t = Time(parsed)
		return nil
	}

	// Try Go's default format
	parsed, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", timeStr)
	if err == nil {
		*t = Time(parsed)
		return nil
	}

	return err
}

func (t Time) MarshalJSON() ([]byte, error) {
	// Omit if zero
	if time.Time(t).IsZero() {
		return []byte("null"), nil
	}

	// Format as RFC3339 string
	timeStr := time.Time(t).Format(time.RFC3339)
	return json.Marshal(timeStr)
}

func (t Time) String() string {
	return time.Time(t).String()
}

func (t Time) QueryValue() string {
	if time.Time(t).IsZero() {
		return ""
	}
	return strconv.FormatInt(time.Time(t).UnixNano()/1e6, 10)
}
