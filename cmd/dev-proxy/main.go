package main

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/pkg/chiext"
	"github.com/go-chi/chi/v5"
	"github.com/phsym/console-slog"
)

type Config struct {
	Address string
	Servers []Server
}

type Server struct {
	URL    string
	Routes []string
}

func main() {
	start(Config{
		Address: ":3000",
		Servers: []Server{
			{
				URL:    "http://127.0.0.1:8888",
				Routes: []string{"/api"},
			},
			{
				URL:    "http://127.0.0.1:5173",
				Routes: []string{"/"},
			},
		},
	})
}

func start(cfg Config) {
	slog.SetDefault(slog.New(console.NewHandler(os.Stderr, &console.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	r := chi.NewMux()
	r.Use(chiext.Logger())

	for _, server := range cfg.Servers {
		urL := core.Must2(url.Parse(server.URL))
		func(urL *url.URL) {
			for _, route := range server.Routes {
				r.Mount(route, &httputil.ReverseProxy{
					Rewrite: func(r *httputil.ProxyRequest) {
						r.SetURL(urL)
						r.SetXForwarded()
					},
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
					},
				})
			}
		}(urL)
	}

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	server.ListenAndServe()
}
