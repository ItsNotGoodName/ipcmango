package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ItsNotGoodName/ipcmanview/internal/system"
)

func NewServer(httpServer *http.Server, certificate *system.Certificate) Server {
	return Server{
		server:      httpServer,
		certificate: certificate,
	}
}

type Server struct {
	server      *http.Server
	certificate *system.Certificate
}

func (s Server) String() string {
	return "app.Server"
}

func (s Server) Serve(ctx context.Context) error {
	slog.Info("Starting HTTP server", "address", s.server.Addr)

	errC := make(chan error, 1)
	go func() {
		if s.certificate != nil {
			errC <- s.server.ListenAndServeTLS(s.certificate.CertFile, s.certificate.KeyFile)
		} else {
			errC <- s.server.ListenAndServe()
		}
	}()

	select {
	case err := <-errC:
		return err
	case <-ctx.Done():
		return s.server.Shutdown(context.Background())
	}
}
