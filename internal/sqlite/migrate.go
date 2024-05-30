package sqlite

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/pressly/goose/v3"
)

//go:embed trigger.sql
var trigger string

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(ctx context.Context, db *sql.DB) (*sql.DB, error) {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return nil, err
	}

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return nil, err
	}

	_, err := db.ExecContext(ctx, trigger)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func init() {
	goose.SetLogger(&logger{})
}

type logger struct{}

func (*logger) Fatalf(format string, v ...interface{}) {
	slog.Error(strings.TrimSuffix(fmt.Sprintf(format, v...), "\n"))
	os.Exit(1)
}

func (*logger) Printf(format string, v ...interface{}) {
	slog.Info(strings.TrimSuffix(fmt.Sprintf(format, v...), "\n"))
}
