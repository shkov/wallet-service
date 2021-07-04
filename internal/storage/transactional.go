package storage

import (
	"context"
	"fmt"
)

type TransactionalStorage interface {
	Storage
	ExecTx(ctx context.Context, fn func(context.Context, Storage) error) error
}

type transactionalStorage struct {
	*storageImpl
}

func NewTransactional(cfg Config) TransactionalStorage {
	return &transactionalStorage{
		storageImpl: newStorageImpl(cfg),
	}
}

// Exec wraps the execution of fn into a postgres transaction.
func (ts *transactionalStorage) ExecTx(ctx context.Context, fn func(context.Context, Storage) error) (err error) {
	tx, err := ts.storageImpl.db.Begin()
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

	err = fn(ctx, ts.storageImpl)
	if err != nil {
		return fmt.Errorf("failed to exec fn: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
