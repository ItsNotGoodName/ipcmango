package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

var timeFormats = []string{
	"2006-01-02 15:04:05.000000",
	"2006-01-02 15:04:05.000000 +0000 UTC",
	"2006-01-02 15:04:05",
}

func NewTime(time time.Time) Time {
	return Time{
		Time: time,
	}
}

// Time will always UTC.
type Time struct {
	time.Time
}

func (dst *Time) Scan(src any) error {
	switch src := src.(type) {
	case time.Time:
		dst.Time = src.UTC()
		return nil
	case string:
		for _, f := range timeFormats {
			t, err := time.ParseInLocation(f, src, time.UTC)
			if err != nil {
				continue
			}
			dst.Time = t
			return nil
		}

		return fmt.Errorf("parsing time %s", src)
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src Time) Value() (driver.Value, error) {
	return src.Time.UTC().Format(timeFormats[0]), nil
}
