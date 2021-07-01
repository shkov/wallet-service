package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/shkov/wallet-service/internal/account"
)

// instrumentingMiddleware wraps Storage and records metrics.
type instrumentingMiddleware struct {
	next      Storage
	histogram metrics.Histogram
}

func NewInstrumentingMiddleware(next Storage, prefix string) Storage {
	return &instrumentingMiddleware{
		next: next,
		histogram: kitprometheus.NewHistogramFrom(
			prometheus.HistogramOpts{
				Name:    prefix + "_storage_queries",
				Buckets: prometheus.ExponentialBuckets(0.01, 2, 7),
			},
			[]string{"method", "error"},
		),
	}
}

func (mw *instrumentingMiddleware) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	createdAt := time.Now()
	out, err := mw.next.GetAccount(ctx, id)
	mw.record(createdAt, "GetAccount", err)
	return out, err
}

func (mw *instrumentingMiddleware) GetAccounts(ctx context.Context, ids []int64) ([]*account.Account, error) {
	createdAt := time.Now()
	out, err := mw.next.GetAccounts(ctx, ids)
	mw.record(createdAt, "GetAccounts", err)
	return out, err
}

func (mw *instrumentingMiddleware) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	createdAt := time.Now()
	out, err := mw.next.GetPayments(ctx, accountID)
	mw.record(createdAt, "GetPayments", err)
	return out, err
}

func (mw *instrumentingMiddleware) InsertPayment(ctx context.Context, p *account.Payment) error {
	createdAt := time.Now()
	err := mw.next.InsertPayment(ctx, p)
	mw.record(createdAt, "InsertPayment", err)
	return err
}

func (mw *instrumentingMiddleware) ReplaceAccounts(ctx context.Context, aa []*account.Account) error {
	createdAt := time.Now()
	err := mw.next.ReplaceAccounts(ctx, aa)
	mw.record(createdAt, "ReplaceAccounts", err)
	return err
}

func (mw *instrumentingMiddleware) Close(ctx context.Context) error {
	createdAt := time.Now()
	err := mw.next.Close(ctx)
	mw.record(createdAt, "Close", err)
	return err
}

func (mw *instrumentingMiddleware) record(beginTime time.Time, method string, err error) {
	labels := []string{"method", method, "error", strconv.FormatBool(err != nil)}
	mw.histogram.With(labels...).Observe(time.Since(beginTime).Seconds())
}
