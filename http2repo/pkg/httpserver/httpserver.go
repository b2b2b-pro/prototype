package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/b2b2b-pro/prototype/http2repo/config"
	"go.uber.org/zap"
)

// TODO: ? константы перенести в конфиг?
const (
	_defaultReadTimeout = 5 * time.Second

// _defaultWriteTimeout    = 5 * time.Second
// _defaultShutdownTimeout = 3 * time.Second
)

type HTTPServer struct {
	s             *http.Server
	notify        chan error
	serverCtx     context.Context
	serverStopCtx context.CancelFunc
}

func New(cfg *config.Config, h http.Handler) *HTTPServer {
	hs := &HTTPServer{
		s: &http.Server{
			Addr:              ":" + cfg.ConfigHTTP.Port,
			Handler:           h,
			ReadHeaderTimeout: _defaultReadTimeout,
		},
		notify: make(chan error, 1),
	}

	return hs
}

func (hs *HTTPServer) Start() {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	hs.serverCtx = serverCtx
	hs.serverStopCtx = serverStopCtx

	go func() {
		// так включается https: hs.notify <- hs.s.ListenAndServeTLS("server.crt", "server.key")
		hs.notify <- hs.s.ListenAndServe()
		close(hs.notify)
	}()
}

func (hs *HTTPServer) Notify() <-chan error {
	return hs.notify
}

func (hs *HTTPServer) Stop() {
	zap.S().Info("HTTP Server получил команду остановиться")
	// TODO вынести в config таймаут
	shutdownCtx, _ := context.WithTimeout(hs.serverCtx, 45*time.Second)

	go func() {
		<-shutdownCtx.Done()
		if shutdownCtx.Err() == context.DeadlineExceeded {
			zap.S().Fatal("graceful shutdown timed out.. forcing exit.")
		}
	}()

	err := hs.s.Shutdown(shutdownCtx)
	if err != nil {
		zap.S().Fatal(err)
	}
	hs.serverStopCtx()
}
