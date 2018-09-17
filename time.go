package luno

import (
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
	return []byte(t.String()), nil
}

func (t Time) String() string {
	return strconv.FormatInt(time.Time(t).UnixNano()/1e6, 10)
}
