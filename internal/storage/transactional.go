package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
)

type TransactionalStorage interface {
	Storage
	ExecTx(ctx context.Context, fn func(context.Context, Storage) error) error
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

type transactionalStorage struct {
	*storageImpl
	conn *pg.DB
}

// NewTransactional creates a new transactional storage.
func NewTransactional(cfg Config) TransactionalStorage {
	conn := pg.Connect(&pg.Options{
		Addr:         cfg.Host + ":" + cfg.Port,
		User:         cfg.User,
		Password:     cfg.Password,
		Database:     cfg.Database,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
	return &transactionalStorage{
		storageImpl: newStorageImpl(conn),
		conn:        conn,
	}
}

func (ts *transactionalStorage) Close() error {
	return ts.conn.Close()
}

// Exec wraps the execution of fn into a postgres transaction.
func (ts *transactionalStorage) ExecTx(ctx context.Context, fn func(context.Context, Storage) error) (err error) {
	tx, err := ts.conn.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			return
		}

		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			err = fmt.Errorf("failed to rollback transaction: %w", rollbackErr)
		}
	}()

	err = fn(ctx, newStorageImpl(tx))
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
