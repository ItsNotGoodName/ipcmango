package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/ItsNotGoodName/ipcmanview/internal/system"
)

func NewHTTPServer(httpServer *http.Server, certificate *system.Certificate) HTTPServer {
	return HTTPServer{
		server:      httpServer,
		certificate: certificate,
	}
}

type HTTPServer struct {
	server      *http.Server
	certificate *system.Certificate
}

func (s HTTPServer) String() string {
	return "app.HTTPServer"
}

func (s HTTPServer) Serve(ctx context.Context) error {
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
