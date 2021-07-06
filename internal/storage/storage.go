package storage

import (
	"context"
	"errors"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"

	"github.com/shkov/wallet-service/internal/account"
)

// Storage represents the wallet-service storage.
type Storage interface {
	GetAccount(ctx context.Context, id int64) (*account.Account, error)
	GetAccounts(ctx context.Context, ids []int64) ([]*account.Account, error)
	GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error)
	InsertPayment(ctx context.Context, p *account.Payment) error
	ReplaceAccounts(ctx context.Context, aa []*account.Account) error
}

type storageImpl struct {
	db orm.DB
}

// new creates a new Storage.
func newStorageImpl(db orm.DB) *storageImpl {
	return &storageImpl{
		db: db,
	}
}

func (s *storageImpl) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	a := &account.Account{}
	err := s.db.ModelContext(ctx, a).Where(`id = ?`, id).Select()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, account.ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

func (s *storageImpl) GetAccounts(ctx context.Context, ids []int64) ([]*account.Account, error) {
	var accounts []*account.Account
	err := s.db.ModelContext(ctx, &accounts).WhereIn(`id in (?)`, ids).Select()
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *storageImpl) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	payments := make([]*account.Payment, 0)
	err := s.db.ModelContext(ctx, &payments).
		WhereOr("from_account_id = ?", accountID).
		WhereOr("to_account_id = ?", accountID).
		Select()
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (s *storageImpl) InsertPayment(ctx context.Context, p *account.Payment) error {
	_, err := s.db.ModelContext(ctx, p).Insert()
	if err != nil {
		return err
	}
	return nil
}

func (s *storageImpl) ReplaceAccounts(ctx context.Context, aa []*account.Account) error {
	_, err := s.db.ModelContext(ctx, &aa).
		OnConflict(`(id) do update`).
		Set(`balance = excluded.balance`).
		Insert()
	if err != nil {
		return err
	}
	return nil
}
