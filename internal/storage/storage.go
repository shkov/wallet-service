package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/shkov/wallet-service/internal/account"
)

type Storage interface {
	GetAccount(ctx context.Context, id int64) (*account.Account, error)
	GetAccounts(ctx context.Context, ids []int64) ([]*account.Account, error)
	GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error)
	InsertPayment(ctx context.Context, p *account.Payment) error
	ReplaceAccounts(ctx context.Context, aa []*account.Account) error
	Close(ctx context.Context) error
}

type Config struct {
	Host           string
	Port           uint16
	Database       string
	User           string
	Password       string
	ConnectTimeout time.Duration
}

type storageImpl struct {
	conn *pgx.Conn
}

func New(cfg Config) (Storage, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?connect_timeout=%d",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		int(cfg.ConnectTimeout.Seconds()),
	)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	s := &storageImpl{
		conn: conn,
	}

	return s, nil
}

func (s *storageImpl) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *storageImpl) GetAccount(ctx context.Context, id int64) (*account.Account, error) {
	return nil, nil
}

func (s *storageImpl) GetAccounts(ctx context.Context, ids []int64) ([]*account.Account, error) {
	return nil, nil
}

func (s *storageImpl) GetPayments(ctx context.Context, accountID int64) ([]*account.Payment, error) {
	return nil, nil
}

func (s *storageImpl) InsertPayment(ctx context.Context, p *account.Payment) error {
	return nil
}

func (s *storageImpl) ReplaceAccounts(ctx context.Context, aa []*account.Account) error {
	return nil
}
