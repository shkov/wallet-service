package walletservice

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Storage interface {
}

// ServerConfig is a server configuration.
type ServerConfig struct {
	Logger          log.Logger
	Storage         Storage
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
	router.Handle("/api/v1/", makeHandler(svc))

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

func makeHandler(svc Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	router := mux.NewRouter()

	router.Path("/api/v1/accounts/{id}").Methods(http.MethodGet).Handler(kithttp.NewServer(
		makeGetAccountEndpoint(svc),
		decodeGetAccountRequest,
		encodeGetAccountResponse,
		opts...,
	))

	router.Path("/api/v1/payments/{account_id}").Methods(http.MethodGet).Handler(kithttp.NewServer(
		makeGetPaymentsEndpoint(svc),
		decodeGetPaymentsRequest,
		encodeGetPaymentsResponse,
		opts...,
	))

	router.Path("/api/v1/payments").Methods(http.MethodPost).Handler(kithttp.NewServer(
		makeApplyPaymentEndpoint(svc),
		decodeApplyPaymentRequest,
		encodeApplyPaymentResponse,
		opts...,
	))

	return router
}

func makeGetAccountEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getAccountRequest)
		resp, err := svc.GetAccount(ctx, req.id)
		if err != nil {
			return nil, err
		}
		return getAccountResponse{account: resp}, nil
	}
}

func makeGetPaymentsEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getPaymentsRequest)
		resp, err := svc.GetPayments(ctx, req.accountID)
		if err != nil {
			return nil, err
		}
		return getPaymentsResponse{payments: resp}, nil
	}
}

func makeApplyPaymentEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(applyPaymentRequest)
		err := svc.ApplyPayment(ctx, req.input)
		if err != nil {
			return nil, err
		}
		return applyPaymentResponse{}, nil
	}
}
