package core

import (
	"database/sql"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
)

type Key struct {
	ID   int64
	UUID string
}

func SplitAddress(address string) (host string, port string) {
	var err error
	host, port, err = net.SplitHostPort(address)
	if err != nil {
		host = address
	}
	return
}

func Address(host string, port int) string {
	return host + ":" + strconv.Itoa(port)
}

type MultiReadCloser struct {
	io.Reader
	Closers []func() error
}

func (c MultiReadCloser) Close() error {
	var multiErr error
	for _, closer := range c.Closers {
		err := closer()
		if err != nil {
			multiErr = errors.Join(multiErr, err)
		}
	}
	return multiErr
}

func NewNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func Int64ToNullInt64(i int64) sql.NullInt64 {
	if i == 0 {
		return sql.NullInt64{}
	}
	return NewNullInt64(i)
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func StringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return NewNullString(s)
}

func ErrorToNullString(err error) sql.NullString {
	if err == nil {
		return sql.NullString{}
	}
	return NewNullString(err.Error())
}

// https://stackoverflow.com/a/12518877
func FileExists(filePath string) (bool, error) {
	if _, err := os.Stat(filePath); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else {
		return false, err
	}
}

func IgnoreError[T any](data T, err error) T {
	return data
}

func FlagChannel(c chan<- struct{}) {
	select {
	case c <- struct{}{}:
	default:
	}
}

func First(s ...string) string {
	for _, s := range s {
		if s != "" {
			return s
		}
	}
	return ""
}

func Optional[T any](optional *T, defaulT T) T {
	if optional != nil {
		return *optional
	}
	return defaulT
}

func SQLNullToNull[T any](t sql.Null[T]) *T {
	if t.Valid {
		return &t.V
	}
	return nil
}

func NullToSQLNull[T any](t *T) sql.Null[T] {
	if t == nil {
		return sql.Null[T]{
			Valid: false,
		}
	}
	return sql.Null[T]{
		V:     *t,
		Valid: true,
	}
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func Must2[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
