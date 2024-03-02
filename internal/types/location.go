package types

import (
	"database/sql/driver"
	"fmt"
	"time"
)

func NewLocation(loc *time.Location) Location {
	if loc == nil {
		loc = time.Local
	}
	return Location{
		Location: loc,
	}
}

// Location cannot be nil.
type Location struct {
	*time.Location
}

func (dst *Location) Scan(src any) error {
	switch src := src.(type) {
	case string:
		loc, err := time.LoadLocation(string(src))
		if err != nil {
			return err
		}
		*dst = Location{loc}
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src Location) Value() (driver.Value, error) {
	return src.Location.String(), nil
}

func (l *Location) MarshalJSON() ([]byte, error) {
	return []byte(l.Location.String()), nil
}

func (l *Location) UnmarshalJSON(data []byte) error {
	loc, err := time.LoadLocation(string(data))
	if err != nil {
		return err
	}
	*l = Location{loc}
	return nil
}

func (src Location) MarshalText() (text []byte, err error) {
	return []byte(src.String()), nil
}

func (dst *Location) UnmarshalText(text []byte) error {
	var err error
	dst.Location, err = time.LoadLocation(string(text))
	return err
}
