package utils

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type TransactionFunc func(ctx context.Context) error

func WithTransaction(db *gorm.DB, ctx context.Context, fn TransactionFunc) error {
	tx := db.Begin().WithContext(ctx)
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Sử dụng defer để recover từ panic
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-throw panic sau khi rollback
		}
	}()

	if err := fn(ctx); err != nil {
		if rollbackErr := tx.Rollback().Error; rollbackErr != nil {
			return fmt.Errorf("rollback error: %v (original error: %w)", rollbackErr, err)
		}
		return err
	}

	if commitErr := tx.Commit().Error; commitErr != nil {
		return fmt.Errorf("commit error: %w", commitErr)
	}

	return nil
}
