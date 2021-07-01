package walletservice

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/shkov/wallet-service/internal/account"
)

// loggingMiddleware wraps the given Service and logs errors.
type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func NewLoggingMiddleware(next Service, logger log.Logger) Service {
	return &loggingMiddleware{
		next:   next,
		logger: logger,
	}
}

func (mw *loggingMiddleware) ApplyPayment(ctx context.Context, p *account.PaymentRequest) (*account.Payment, error) {
	startedAt := time.Now()
	out, err := mw.next.ApplyPayment(ctx, p)
	mw.log(ctx, startedAt, "ApplyPayment", err)
	return out, err
}

func (mw *loggingMiddleware) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	startedAt := time.Now()
	out, err := mw.next.GetPayments(ctx, accountID)
	mw.log(ctx, startedAt, "GetPayments", err)
	return out, err
}

func (mw *loggingMiddleware) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	startedAt := time.Now()
	out, err := mw.next.GetAccount(ctx, id)
	mw.log(ctx, startedAt, "GetAccount", err)
	return out, err
}

func (mw *loggingMiddleware) log(ctx context.Context, beginTime time.Time, method string, err error) {
	if err != nil {
		level.Error(mw.logger).Log("method", method, "err", err, "took", time.Since(beginTime))
	}
}
