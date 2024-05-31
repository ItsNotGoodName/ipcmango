package core

import (
	"database/sql"
	"errors"
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

func FlagChannel(c chan<- struct{}) {
	select {
	case c <- struct{}{}:
	default:
	}
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

func DirectorySize(dir string) (int64, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}

	dirSize := int64(0)
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return 0, err
		}

		if info.Mode().IsRegular() {
			dirSize += info.Size()
		}
	}

	return dirSize, nil
}
