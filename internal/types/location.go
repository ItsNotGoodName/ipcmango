package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fxamacker/cbor/v2"
)

func NewLocation(location *time.Location) Location {
	return Location{
		Location: location,
	}
}

type Location struct {
	Location *time.Location
}

func (dst *Location) Scan(src any) error {
	switch src := src.(type) {
	case string:
		loc, err := time.LoadLocation(src)
		if err != nil {
			return err
		}
		*dst = Location{
			Location: loc,
		}
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src Location) Value() (driver.Value, error) {
	return src.Location.String(), nil
}

func (l *Location) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(l.Location.String())
	return b, err
}

func (l *Location) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	loc, err := time.LoadLocation(s)
	if err != nil {
		return err
	}
	*l = Location{
		Location: loc,
	}
	return nil
}

func (l *Location) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(l.Location.String())
}

func (l *Location) UnmarshalCBOR(data []byte) error {
	var s string
	if err := cbor.Unmarshal(data, &s); err != nil {
		return err
	}

	loc, err := time.LoadLocation(s)
	if err != nil {
		return err
	}
	*l = Location{
		Location: loc,
	}
	return nil
}

func (l *Location) Schema(r huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type:     huma.TypeString,
		Examples: []any{"UTC"},
	}
}
