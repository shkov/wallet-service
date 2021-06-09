package walletservice

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ServerConfig is a server configuration.
type ServerConfig struct {
	Logger          log.Logger
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	MetricPrefix    string
}

// Server is a wallet-service server.
type Server struct {
	cfg *ServerConfig
	srv *http.Server
}

// NewServer creates a new server.
func NewServer(cfg ServerConfig) (*Server, error) {
	var svc Service
	svc = newService(cfg.Logger)
	svc = NewLoggingMiddleware(svc, cfg.Logger)
	svc = NewInstrumentingMiddleware(svc, cfg.MetricPrefix)

	router := http.NewServeMux()
	router.Handle("/metrics", promhttp.Handler())
	router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	router.Handle("/debug/pprof/heap", pprof.Handler("heap"))

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + cfg.Port,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	s := &Server{
		cfg: &cfg,
		srv: srv,
	}
	return s, nil
}

// Serve starts HTTP server and stops it when the provided context is canceled.
func (s *Server) Serve(ctx context.Context) error {
	errChan := make(chan error, 1)
	go func() {
		errChan <- s.srv.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		return err

	case <-ctx.Done():
		ctxShutdown, cancel := context.WithTimeout(context.Background(), s.cfg.ShutdownTimeout)
		defer cancel()
		if err := s.srv.Shutdown(ctxShutdown); err != nil {
			return fmt.Errorf("failed to shutdown server: %w", err)
		}
		return nil
	}
}
