package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/fxamacker/cbor/v2"
)

func NewDuration(duration time.Duration) Duration {
	return Duration{
		Duration: duration,
	}
}

type Duration struct {
	Duration time.Duration
}

func (dst *Duration) Scan(src any) error {
	switch src := src.(type) {
	case string:
		duration, err := time.ParseDuration(src)
		if err != nil {
			return err
		}
		*dst = Duration{
			Duration: duration,
		}
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src Duration) Value() (driver.Value, error) {
	return src.Duration.String(), nil
}

func (l *Duration) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(l.Duration.String())
	return b, err
}

func (l *Duration) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*l = Duration{
		Duration: duration,
	}
	return nil
}

func (l *Duration) MarshalCBOR() ([]byte, error) {
	return cbor.Marshal(l.Duration.String())
}

func (l *Duration) UnmarshalCBOR(data []byte) error {
	var s string
	if err := cbor.Unmarshal(data, &s); err != nil {
		return err
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*l = Duration{
		Duration: duration,
	}
	return nil
}

func (l *Duration) Schema(r huma.Registry) *huma.Schema {
	return &huma.Schema{
		Type:     huma.TypeString,
		Examples: []any{"0s"},
	}
}
