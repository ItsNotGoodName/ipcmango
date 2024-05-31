package quartzext

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/reugn/go-quartz/logger"
	"github.com/reugn/go-quartz/quartz"
)

func NewServiceScheduler(s quartz.Scheduler) ServiceScheduler {
	logger.SetDefault(&Logger{})
	return ServiceScheduler{
		Scheduler: s,
	}
}

type ServiceScheduler struct {
	quartz.Scheduler
}

func (s ServiceScheduler) String() string {
	return "quartzext.ServiceScheduler"
}

func (s ServiceScheduler) Serve(ctx context.Context) error {
	s.Start(ctx)
	s.Wait(context.Background())
	return nil
}

func NewJobDetail(job quartz.Job) *quartz.JobDetail {
	return quartz.NewJobDetail(job, quartz.NewJobKey(job.Description()))
}

type Logger struct{}

// Debug implements logger.Logger.
func (l *Logger) Debug(msg any) {
	slog.Debug(fmt.Sprint(msg))
}

// Debugf implements logger.Logger.
func (l *Logger) Debugf(format string, args ...any) {
	slog.Debug(fmt.Sprintf(format, args...))
}

// Enabled implements logger.Logger.
func (l *Logger) Enabled(level logger.Level) bool {
	slogLevel := slog.LevelDebug
	switch level {
	case logger.LevelTrace:
		slogLevel = slog.LevelDebug
	case logger.LevelDebug:
		slogLevel = slog.LevelDebug
	case logger.LevelInfo:
		slogLevel = slog.LevelInfo
	case logger.LevelWarn:
		slogLevel = slog.LevelWarn
	case logger.LevelError:
		slogLevel = slog.LevelError
	case logger.LevelOff:
		return false
	}
	return slog.Default().Enabled(nil, slogLevel)
}

// Error implements logger.Logger.
func (l *Logger) Error(msg any) {
	slog.Error(fmt.Sprint(msg))
}

// Errorf implements logger.Logger.
func (l *Logger) Errorf(format string, args ...any) {
	slog.Error(fmt.Sprintf(format, args...))
}

// Info implements logger.Logger.
func (l *Logger) Info(msg any) {
	slog.Error(fmt.Sprint(msg))
}

// Infof implements logger.Logger.
func (l *Logger) Infof(format string, args ...any) {
	slog.Info(fmt.Sprintf(format, args...))
}

// Trace implements logger.Logger.
func (l *Logger) Trace(msg any) {
	slog.Debug(fmt.Sprint(msg))
}

// Tracef implements logger.Logger.
func (l *Logger) Tracef(format string, args ...any) {
	slog.Debug(fmt.Sprintf(format, args...))
}

// Warn implements logger.Logger.
func (l *Logger) Warn(msg any) {
	slog.Warn(fmt.Sprint(msg))
}

// Warnf implements logger.Logger.
func (l *Logger) Warnf(format string, args ...any) {
	slog.Warn(fmt.Sprintf(format, args...))
}

var _ logger.Logger = (*Logger)(nil)
