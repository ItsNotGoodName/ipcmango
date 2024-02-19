package core

import (
	"database/sql"
	"errors"
	"io"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ItsNotGoodName/ipcmanview/internal/models"
)

type contextKey string

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

func NewTimeRange(start, end time.Time) (models.TimeRange, error) {
	if end.Before(start) {
		return models.TimeRange{}, errors.New("invalid time range: end is before start")
	}

	return models.TimeRange{
		Start: start,
		End:   end,
	}, nil
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

func Int64ToNullInt64(a int64) sql.NullInt64 {
	if a == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: a,
		Valid: true,
	}
}

func NewNullInt64(i int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: i,
		Valid: true,
	}
}

func NewNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func ErrorToNullString(err error) sql.NullString {
	if err == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: err.Error(),
		Valid:  true,
	}
}

func NilStringToNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{
			Valid: true,
		}
	}
	return sql.NullString{
		String: *s,
		Valid:  true,
	}
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
