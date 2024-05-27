package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ItsNotGoodName/ipcmanview/internal/app"
	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/pkg/chiext"
	"github.com/ItsNotGoodName/ipcmanview/pkg/quartzext"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/phsym/console-slog"
	"github.com/reugn/go-quartz/quartz"
	"github.com/spf13/afero"
	"github.com/thejerf/suture/v4"

	_ "github.com/k0kubun/pp/v3"
)

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var
type Options struct {
	Debug    bool   `doc:"Enable debug"`
	Host     string `doc:"Host to listen on"`
	Port     int    `doc:"Port to listen on" short:"p" default:"8888"`
	SmtpHost string `doc:"SMTP host to listen on"`
	SmtpPort int    `doc:"SMTP port to listen on" default:"1025"`
	Dir      string `doc:"Directory to store data" short:"d" default:"ipcmanview_data"`
}

func main() {
	// Create a CLI app which takes a port option
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		// Init logger
		level := slog.LevelInfo
		if options.Debug {
			level = slog.LevelDebug
		}
		slog.SetDefault(slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{
			Level: level,
		})))

		// Create suture root
		root := suture.New("root", suture.Spec{
			EventHook: sutureext.EventHook(),
		})

		// Create data directory
		dir := core.Must2(filepath.Abs(options.Dir))
		core.Must(os.MkdirAll(dir, 0755))

		// Create sqlite database
		sqliteDB := core.Must2(sqlite.New(filepath.Join(dir, "sqlite.db")))
		core.Must(sqlite.Migrate(sqliteDB))
		db := sqlx.NewDb(sqliteDB, sqlite.Driver)

		// Create afero filesystem
		afsPath := filepath.Join(dir, "afero")
		core.Must(os.MkdirAll(afsPath, 0755))
		afs := afero.NewBasePathFs(afero.NewOsFs(), afsPath)

		// Create quartz scheduler
		scheduler := quartzext.NewServiceScheduler(quartz.NewStdScheduler())
		root.Add(scheduler)

		// Create dahuaStore
		dahuaStore := dahua.NewStore()
		root.Add(dahuaStore)

		// Create dahuaEventManager
		dahuaEventManager := dahua.NewEventManager(root, db).Register()
		root.Add(dahuaEventManager)

		dahua.RegisterEmailToEndpoints(db, afs)

		// Schedule jobs
		core.Must(scheduler.ScheduleJob(
			quartzext.NewJobDetail(dahua.NewSyncVideoInModeJob(db, dahuaStore)),
			core.Must2(quartz.NewCronTrigger("0 0 */8 * * *")), // Every 8 hours
		))
		core.Must(scheduler.ScheduleJob(
			quartzext.NewJobDetail(dahua.NewDeleteOrphanEmailAttachmentsJob(db, afs)),
			core.Must2(quartz.NewCronTrigger("0 0 * * * *")), // Every hour
		))

		// Create a new router & API
		router := chi.NewMux()
		router.Use(middleware.RequestLogger(&chiext.DefaultLogFormatter{}))
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		// Register handlers
		app.Register(api, db, afs, dahuaStore)

		// Create the httpServer
		httpServer := app.NewHTTPServer(&http.Server{
			Addr:    core.Address("", options.Port),
			Handler: router,
		})
		root.Add(httpServer)

		// Create the smtpServer
		smtpServer := app.NewSMTPServer(db, afs, core.Address("", options.SmtpPort))
		root.Add(smtpServer)

		stopC := make(chan struct{})

		// Tell the CLI how to start your application
		hooks.OnStart(func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			core.Must(system.InitializeSettings(ctx, db))

			bus.SetContext(ctx)

			errC := root.ServeBackground(ctx)

			select {
			case <-stopC:
				cancel()
			case <-errC:
				return
			}

			<-errC
			<-stopC
		})

		// Tell the CLI how to stop your server
		hooks.OnStop(func() {
			stopC <- struct{}{}
			stopC <- struct{}{}
		})
	})

	// Run the CLI. When passed no commands, it starts the server
	cli.Run()
}
