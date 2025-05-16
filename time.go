package luno

import (
	"os"
	"strconv"
	"time"
)

type Time time.Time

func (t *Time) UnmarshalJSON(b []byte) error {
	i, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	if i == 0 {
		*t = Time{}
		return nil
	}
	*t = Time(time.Unix(0, i*1e6))
	return nil
}

func (t Time) MarshalJSON() ([]byte, error) {
	// Check if LUNO_TIME_LEGACY_FORMAT environment variable is set to control
	// whether to use the original string format (when set to "true")
	// or the new millisecond timestamp format (default)
	useLegacyFormat := os.Getenv("LUNO_TIME_LEGACY_FORMAT") == "true"

	if time.Time(t).IsZero() {
		return []byte("0"), nil
	}

	if useLegacyFormat {
		// Return properly quoted string in the original format
		return []byte(`"` + time.Time(t).String() + `"`), nil
	}

	// Return millisecond timestamp (numeric format)
	ms := time.Time(t).UnixNano() / 1e6
	return []byte(strconv.FormatInt(ms, 10)), nil
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
