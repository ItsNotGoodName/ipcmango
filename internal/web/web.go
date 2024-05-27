//go:build !static

package web

import (
	"net/http"
)

func FS(skipPaths ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return next
	}
}
