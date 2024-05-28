package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ItsNotGoodName/ipcmanview/internal/app"
	"github.com/ItsNotGoodName/ipcmanview/internal/build"
	"github.com/ItsNotGoodName/ipcmanview/internal/bus"
	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/dahua"
	"github.com/ItsNotGoodName/ipcmanview/internal/sqlite"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/internal/web"
	"github.com/ItsNotGoodName/ipcmanview/pkg/chiext"
	"github.com/ItsNotGoodName/ipcmanview/pkg/quartzext"
	"github.com/ItsNotGoodName/ipcmanview/pkg/sutureext"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/phsym/console-slog"
	"github.com/reugn/go-quartz/quartz"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/thejerf/suture/v4"

	_ "github.com/k0kubun/pp/v3"
)

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var
type Options struct {
	Debug    bool   `doc:"enable debug"`
	Dir      string `doc:"directory to store data" short:"d" default:"ipcmanview_data"`
	Host     string `doc:"host to listen on"`
	Port     int    `doc:"port to listen on" short:"p" default:"8888"`
	SmtpHost string `doc:"smtp host to listen on"`
	SmtpPort int    `doc:"smtp port to listen on" default:"1025"`
}

func main() {
	godotenv.Load()

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

		stopC := make(chan struct{})

		// Tell the CLI how to start your application
		hooks.OnStart(func() {
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

			// Create a new router
			router := chi.NewMux()
			router.Use(chiext.Logger())
			router.Use(web.FS("/api", "/openapi", "/docs"))

			// Create api
			api := humachi.New(router, app.NewHumaConfig())

			// Register handlers
			app.Register(api, db, afs, afsPath, dahuaStore)

			// Create the httpServer
			httpServer := app.NewHTTPServer(&http.Server{
				Addr:    core.Address("", options.Port),
				Handler: router,
			})
			root.Add(httpServer)

			// Create the smtpServer
			smtpServer := app.NewSMTPServer(db, afs, core.Address("", options.SmtpPort))
			root.Add(smtpServer)

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

	cli.Root().Version = build.Current.Version

	cli.Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec",
		Run: func(cmd *cobra.Command, args []string) {
			api := humachi.New(chi.NewMux(), app.NewHumaConfig())
			app.Register(api, nil, nil, "", nil)
			b, _ := json.MarshalIndent(api.OpenAPI(), "", "  ")
			fmt.Println(string(b))
		},
	})

	// Run the CLI. When passed no commands, it starts the server
	cli.Run()
}
