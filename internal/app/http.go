package app

import (
	"context"
	"log/slog"
	"net/http"
)

func NewHTTPServer(httpServer *http.Server) HTTPServer {
	return HTTPServer{
		server: httpServer,
	}
}

type HTTPServer struct {
	server *http.Server
}

func (s HTTPServer) String() string {
	return "app.HTTPServer"
}

func (s HTTPServer) Serve(ctx context.Context) error {
	slog.Info("Starting HTTP server", "address", s.server.Addr)

	errC := make(chan error, 1)
	go func() { errC <- s.server.ListenAndServe() }()

	select {
	case err := <-errC:
		return err
	case <-ctx.Done():
		return s.server.Shutdown(context.Background())
	}
}
