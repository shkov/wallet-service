package walletservice

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/shkov/wallet-service/internal/account"
)

var _ Service = (*instrumentingMiddleware)(nil)

// instrumentingMiddleware wraps the given Service and records metrics.
type instrumentingMiddleware struct {
	next      Service
	histogram metrics.Histogram
}

// NewInstrumentingMiddleware creates a new instrumenting middleware.
func NewInstrumentingMiddleware(next Service, prefix string) Service {
	return &instrumentingMiddleware{
		next: next,
		histogram: kitprometheus.NewHistogramFrom(
			prometheus.HistogramOpts{
				Name:    prefix + "_queries",
				Buckets: prometheus.ExponentialBuckets(0.025, 2, 8),
			},
			[]string{"method", "error"},
		),
	}
}

func (mw *instrumentingMiddleware) ApplyPayment(ctx context.Context, p *account.Payment) error {
	startedAt := time.Now()
	err := mw.next.ApplyPayment(ctx, p)
	mw.record(ctx, startedAt, "ApplyPayment", err)
	return err
}

func (mw *instrumentingMiddleware) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	startedAt := time.Now()
	out, err := mw.next.GetPayments(ctx, accountID)
	mw.record(ctx, startedAt, "GetPayments", err)
	return out, err
}

func (mw *instrumentingMiddleware) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	startedAt := time.Now()
	out, err := mw.next.GetAccount(ctx, id)
	mw.record(ctx, startedAt, "GetAccount", err)
	return out, err
}

func (mw *instrumentingMiddleware) record(ctx context.Context, beginTime time.Time, method string, err error) {
	mw.histogram.With(
		"method", method,
		"error", strconv.FormatBool(err != nil),
	).Observe(time.Since(beginTime).Seconds())
}
