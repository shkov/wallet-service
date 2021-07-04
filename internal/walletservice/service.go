package walletservice

import (
	"context"
	"errors"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/shkov/wallet-service/internal/account"
	"github.com/shkov/wallet-service/internal/storage"
)

// Service provides wallet-service functionality.
type Service interface {
	ApplyPayment(ctx context.Context, p *account.PaymentRequest) (*account.Payment, error)
	GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error)
	GetAccount(ctx context.Context, id int64) (*account.Account, error)
}

type serviceImpl struct {
	logger  log.Logger
	storage storage.TransactionalStorage

	now func() time.Time
}

func newService(logger log.Logger, storage storage.TransactionalStorage) Service {
	return &serviceImpl{
		logger:  logger,
		storage: storage,
		now: func() time.Time {
			return time.Now()
		},
	}
}

// ApplyPayment applies the given payment request to the accounts.
func (s *serviceImpl) ApplyPayment(ctx context.Context, r *account.PaymentRequest) (*account.Payment, error) {
	err := account.ValidatePaymentRequest(r)
	if err != nil {
		return nil, errBadRequest("payment is invalid: %v", err)
	}

	payment := r.ToPayment(s.now())

	txFn := func(ctx context.Context, storage storage.Storage) error {
		fromAccount, toAccount, err := s.getAccountsByPayment(ctx, payment)
		if err != nil {
			return errInternal("failed to get accounts: %v", err)
		}

		err = fromAccount.ApplyPayment(payment)
		if err != nil {
			return errBadRequest("failed to apply payment to the sender: %v", err)
		}

		err = toAccount.ApplyPayment(payment)
		if err != nil {
			return errBadRequest("failed to apply payment to the receiver: %v", err)
		}

		err = s.storage.ReplaceAccounts(ctx, []*account.Account{fromAccount, toAccount})
		if err != nil {
			return errInternal("failed to replace accounts: %v", err)
		}

		err = s.storage.InsertPayment(ctx, payment)
		if err != nil {
			return errInternal("failed to insert payment: %v", err)
		}

		return nil
	}

	err = s.storage.ExecTx(ctx, txFn)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

// getAccountsByPayment tries to get accounts from the storage, if they were not found it creates a new ones.
func (s *serviceImpl) getAccountsByPayment(ctx context.Context, p *account.Payment) (*account.Account, *account.Account, error) {
	accounts, err := s.storage.GetAccounts(ctx, []int64{p.From, p.To})
	if err != nil {
		return nil, nil, err
	}

	var fromAccount, toAccount *account.Account
	for _, a := range accounts {
		if a.ID == p.From {
			fromAccount = a
		}
		if a.ID == p.To {
			toAccount = a
		}
	}

	// TODO: add an ability to create accounts, now new ones have a default balance.

	if fromAccount == nil {
		fromAccount = account.Create(p.From, s.now())
	}

	if toAccount == nil {
		toAccount = account.Create(p.To, s.now())
	}

	return fromAccount, toAccount, nil
}

// GetPayments returns all payments by the account id.
func (s *serviceImpl) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	err := account.ValidateAccountID(accountID)
	if err != nil {
		return nil, errBadRequest("provided account id is invalid: %d", accountID)
	}

	payments, err := s.storage.GetPayments(ctx, accountID)
	if err != nil {
		return nil, errInternal("failed to get payments from the storage: %v", err)
	}

	return payments, nil
}

// GetAccount returns an account by the given id.
func (s *serviceImpl) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	err := account.ValidateAccountID(id)
	if err != nil {
		return nil, errBadRequest("provided account id is invalid: %d", id)
	}

	a, err := s.storage.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, account.ErrNotFound) {
			return nil, errNotFound("account %d is not found", id)
		}
		return nil, errInternal("failed to get account from the storage: %v", err)
	}

	return a, nil
}
