package main

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/ItsNotGoodName/ipcmanview/internal/core"
	"github.com/ItsNotGoodName/ipcmanview/internal/system"
	"github.com/ItsNotGoodName/ipcmanview/pkg/chiext"
	"github.com/go-chi/chi/v5"
	"github.com/phsym/console-slog"
)

type Config struct {
	Address     string
	Servers     []Server
	Certificate system.Certificate
}

type Server struct {
	URL    string
	Mount  []string
	Handle []string
}

func main() {
	start(Config{
		Address: ":3443",
		Certificate: system.Certificate{
			CertFile: "./ipcmanview_data/cert.pem",
			KeyFile:  "./ipcmanview_data/key.pem",
		},
		Servers: []Server{
			{
				URL:    "http://127.0.0.1:8080",
				Mount:  []string{"/api"},
				Handle: []string{"/docs", "/openapi.yaml", "/openapi.json"},
			},
			{
				URL:   "http://127.0.0.1:5173",
				Mount: []string{"/"},
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
			proxy := &httputil.ReverseProxy{
				Rewrite: func(r *httputil.ProxyRequest) {
					r.SetURL(urL)
					r.SetXForwarded()
				},
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}
			for _, route := range server.Mount {
				r.Mount(route, proxy)
			}
			for _, static := range server.Handle {
				r.Handle(static, proxy)
			}
		}(urL)
	}

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	core.Must(server.ListenAndServeTLS(cfg.Certificate.CertFile, cfg.Certificate.KeyFile))
}
