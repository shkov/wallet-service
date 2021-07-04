package storage

import (
	"context"
	"errors"
	"time"

	"github.com/go-pg/pg/v10"

	"github.com/shkov/wallet-service/internal/account"
)

// Storage represents the wallet-service storage.
type Storage interface {
	GetAccount(ctx context.Context, id int64) (*account.Account, error)
	GetAccounts(ctx context.Context, ids []int64) ([]*account.Account, error)
	GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error)
	InsertPayment(ctx context.Context, p *account.Payment) error
	ReplaceAccounts(ctx context.Context, aa []*account.Account) error
	Close() error
}

// Config is a Storage configuration.
type Config struct {
	Host         string
	Port         string
	Database     string
	User         string
	Password     string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type storageImpl struct {
	db *pg.DB
}

// new creates a new Storage.
func newStorageImpl(cfg Config) *storageImpl {
	return &storageImpl{
		db: pg.Connect(&pg.Options{
			Addr:         cfg.Host + ":" + cfg.Port,
			User:         cfg.User,
			Password:     cfg.Password,
			Database:     cfg.Database,
			DialTimeout:  cfg.DialTimeout,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
		}),
	}
}

func (s *storageImpl) Close() error {
	return s.db.Close()
}

func (s *storageImpl) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	a := &account.Account{}
	err := s.db.WithContext(ctx).Model(a).Where(`id = ?`, id).Select()
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
	err := s.db.WithContext(ctx).Model(&accounts).WhereIn(`id in (?)`, ids).Select()
	if err != nil {
		return nil, err
	}
	return accounts, nil
}

func (s *storageImpl) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	payments := make([]*account.Payment, 0)
	err := s.db.WithContext(ctx).
		Model(&payments).
		WhereOr("from_account_id = ?", accountID).
		WhereOr("to_account_id = ?", accountID).
		Select()
	if err != nil {
		return nil, err
	}

	return payments, nil
}

func (s *storageImpl) InsertPayment(ctx context.Context, p *account.Payment) error {
	_, err := s.db.WithContext(ctx).Model(p).Insert()
	if err != nil {
		return err
	}
	return nil
}

func (s *storageImpl) ReplaceAccounts(ctx context.Context, aa []*account.Account) error {
	_, err := s.db.WithContext(ctx).
		Model(&aa).
		OnConflict(`(id) do update`).
		Set(`balance = excluded.balance`).
		Insert()
	if err != nil {
		return err
	}
	return nil
}
