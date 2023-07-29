package main

import (
	"context"
	"os"

	"github.com/ItsNotGoodName/ipcmango/internal/build"
	"github.com/ItsNotGoodName/ipcmango/internal/db"
	"github.com/ItsNotGoodName/ipcmango/migrations"
	"github.com/ItsNotGoodName/ipcmango/pkg/background"
	"github.com/ItsNotGoodName/ipcmango/pkg/interrupt"
	"github.com/ItsNotGoodName/ipcmango/sandbox"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	sandbox.Dahua(interrupt.Context())
	return

	ctx, shutdown := context.WithCancel(interrupt.Context())
	defer shutdown()

	// Database
	pool, err := db.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create database connection pool")
	}
	defer pool.Close()

	// Database migrate
	if err := migrations.Migrate(ctx, pool); err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	// sandbox.Jet(ctx, pool)
	<-background.Run(ctx, sandbox.Chi(ctx, shutdown, pool))
	// sandbox.User(ctx, pool)
}

var (
	builtBy    = "unknown"
	commit     = ""
	date       = ""
	version    = "dev"
	repoURL    = "https://github.com/ItsNotGoodName/smtpbridge"
	releaseURL = ""
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	build.Current = build.Build{
		BuiltBy:    builtBy,
		Commit:     commit,
		Date:       date,
		Version:    version,
		RepoURL:    repoURL,
		ReleaseURL: releaseURL,
	}
}
