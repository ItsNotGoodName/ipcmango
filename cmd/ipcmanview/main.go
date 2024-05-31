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
	"github.com/ItsNotGoodName/ipcmanview/internal/smtp"
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

// Options for the CLI. Pass `--http-port` or set the `SERVICE_HTTP_PORT` env var
type Options struct {
	Debug     bool   `doc:"enable debug"`
	Dir       string `doc:"directory to store data" short:"d" default:"ipcmanview_data"`
	HttpHost  string `doc:"http host to listen on"`
	HttpPort  int    `doc:"http port to listen on" default:"8080"`
	HttpsHost string `doc:"https host to listen on"`
	HttpsPort int    `doc:"https port to listen on" default:"8443"`
	SmtpHost  string `doc:"smtp host to listen on"`
	SmtpPort  int    `doc:"smtp port to listen on" default:"1025"`
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
			// Initialize context
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			bus.SetContext(ctx)

			// Create suture root
			root := suture.New("root", suture.Spec{
				EventHook: sutureext.EventHook(),
			})

			// Create data directory
			dir := core.Must2(filepath.Abs(options.Dir))
			core.Must(os.MkdirAll(dir, 0o755))

			// Create database
			db := sqlx.NewDb(core.Must2(sqlite.Migrate(ctx, core.Must2(sqlite.New(filepath.Join(dir, "sqlite.db"))))), sqlite.Driver)

			// Create afero filesystem
			afsDirectory := filepath.Join(dir, "afero")
			core.Must(os.MkdirAll(afsDirectory, 0o755))
			afs := afero.NewBasePathFs(afero.NewOsFs(), afsDirectory)

			// Create quartz scheduler
			scheduler := quartzext.NewServiceScheduler(quartz.NewStdScheduler())
			root.Add(scheduler)

			// Create dahua store
			dahuaStore := dahua.NewStore(db)
			root.Add(dahuaStore)

			// Create file scan job
			dahuaFileScanJobClient := dahua.NewFileScanJobClient(db)
			dahuaFileScanJob := dahua.RegisterFileScanJob(dahuaFileScanJobClient, db, dahuaStore)
			root.Add(dahuaFileScanJobClient)

			// Create dahua event manager
			root.Add(dahua.NewServiceManager(root, dahua.NewDefaultServiceFactory(db, dahuaStore)).Register())

			// Create dahua email worker
			root.Add(dahua.NewEmailWorker(db, afs).Register())

			// Schedule jobs
			core.Must(scheduler.ScheduleJob(
				quartzext.NewJobDetail(dahua.NewSyncVideoInModeJob(db, dahuaStore)),
				core.Must2(quartz.NewCronTrigger("0 0 */8 * * *")), // Every 8 hours
			))
			core.Must(scheduler.ScheduleJob(
				quartzext.NewJobDetail(dahua.NewDeleteOrphanEmailAttachmentsJob(db, afs)),
				core.Must2(quartz.NewCronTrigger("0 0 * * * *")), // Every hour
			))

			// Create router
			router := chi.NewMux()
			router.Use(chiext.Logger())
			router.Use(web.FS("/api", "/openapi", "/docs"))

			// Create api
			api := humachi.New(router, app.NewConfig())

			// Register handlers
			app.Register(api, app.App{
				DB:           db,
				AFS:          afs,
				AFSDirectory: afsDirectory,
				FileScanJob:  dahuaFileScanJob,
				DahuaStore:   dahuaStore,
			})

			// Create HTTP server
			root.Add(app.NewServer(&http.Server{
				Addr:    core.Address(options.HttpHost, options.HttpPort),
				Handler: router,
			}, nil))

			// Generate certificate
			certificate := system.Certificate{
				CertFile: filepath.Join(dir, "cert.pem"),
				KeyFile:  filepath.Join(dir, "key.pem"),
			}
			core.Must(certificate.GenerateIfNotExist())

			// Create HTTPS server
			root.Add(app.NewServer(&http.Server{
				Addr:    core.Address(options.HttpsHost, options.HttpsPort),
				Handler: router,
			}, &certificate))

			// Create SMTP server
			root.Add(smtp.NewServer(db, afs, core.Address(options.SmtpHost, options.SmtpPort)))

			core.Must(system.InitializeSettings(ctx, db))

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
			api := humachi.New(chi.NewMux(), app.NewConfig())
			app.Register(api, app.App{})
			b, _ := json.MarshalIndent(api.OpenAPI(), "", "  ")
			fmt.Println(string(b))
		},
	})

	// Run the CLI. When passed no commands, it starts the server
	cli.Run()
}
