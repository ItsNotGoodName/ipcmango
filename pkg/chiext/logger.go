package chiext

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func Logger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&DefaultLogFormatter{})
}

// DefaultLogFormatter is a simple logger that implements a LogFormatter.
type DefaultLogFormatter struct{}

// NewLogEntry creates a new LogEntry for the request.
func (l *DefaultLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	attrs := []any{}

	// attrs = append(attrs, slog.String("method", r.Method))
	reqID := middleware.GetReqID(r.Context())
	if reqID != "" {
		attrs = append(attrs, slog.String("request", reqID))
	}
	attrs = append(attrs, slog.String("from", r.RemoteAddr))

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	msg := fmt.Sprintf("%s %s://%s%s %s", r.Method, scheme, r.Host, r.RequestURI, r.Proto)

	return &defaultLogEntry{
		DefaultLogFormatter: l,
		attrs:               attrs,
		msg:                 msg,
	}
}

type defaultLogEntry struct {
	*DefaultLogFormatter
	attrs []any
	msg   string
}

func (l *defaultLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	attrs := append(l.attrs,
		slog.Int("status", status),
		slog.Int("bytes", bytes),
		slog.String("elapsed", elapsed.String()),
	)

	switch {
	case status < 200:
		slog.Info(l.msg, attrs...)
	case status < 300:
		slog.Info(l.msg, attrs...)
	case status < 400:
		slog.Info(l.msg, attrs...)
	case status < 500:
		slog.Info(l.msg, attrs...)
	default:
		slog.Error(l.msg, attrs...)
	}
}

func (l *defaultLogEntry) Panic(v interface{}, stack []byte) {
	middleware.PrintPrettyStack(v)
}
