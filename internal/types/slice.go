package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

func NewSlice[T any](slice []T) Slice[T] {
	return Slice[T]{
		V: slice,
	}
}

type Slice[T any] struct {
	V []T
}

func (dst *Slice[T]) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("cannot scan nil")
	}

	switch src := src.(type) {
	case string:
		return json.Unmarshal([]byte(src), &dst.V)
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src Slice[T]) Value() (driver.Value, error) {
	b, err := json.Marshal(src.V)
	if err != nil {
		return nil, err
	}
	return string(b), nil
}
