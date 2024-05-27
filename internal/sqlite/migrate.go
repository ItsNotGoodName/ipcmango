package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func Migrate(sqlDB *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return err
	}

	return nil
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
