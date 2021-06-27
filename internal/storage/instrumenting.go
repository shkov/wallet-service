package storage

import (
	"strconv"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
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

func (mw *instrumentingMiddleware) record(beginTime time.Time, method string, err error) {
	labels := []string{"method", method, "error", strconv.FormatBool(err != nil)}
	mw.histogram.With(labels...).Observe(time.Since(beginTime).Seconds())
}
