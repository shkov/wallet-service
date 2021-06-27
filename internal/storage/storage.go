package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
)

type Storage interface {
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

func (s *storageImpl) Close() error {
	return s.conn.Close(context.Background())
}
