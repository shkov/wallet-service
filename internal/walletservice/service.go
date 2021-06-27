package walletservice

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/shkov/wallet-service/internal/account"
)

// Service provides wallet-service functionality.
type Service interface {
	ApplyPayment(ctx context.Context, p *account.Payment) error
	GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error)
	GetAccount(ctx context.Context, id int64) (*account.Account, error)
}

type serviceImpl struct {
	logger log.Logger
}

func newService(logger log.Logger) Service {
	return &serviceImpl{
		logger: logger,
	}
}

// ApplyPayment applies the given payment to the accounts.
func (s *serviceImpl) ApplyPayment(ctx context.Context, p *account.Payment) error {
	return errInternal("ApplyPayment wasn't implemented")
}

// GetPayments returns all payments by the account id.
func (s *serviceImpl) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	return nil, errInternal("GetPayments wasn't implemented")
}

// GetAccount returns an account by the given id.
func (s *serviceImpl) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	return nil, errInternal("GetAccount wasn't implemented")
}
