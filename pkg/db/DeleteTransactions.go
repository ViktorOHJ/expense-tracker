package db

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (db *PostgresDB) DeleteTransaction(parentCtx context.Context, userID int, transactionID int) error {
	query := `DELETE FROM transactions WHERE id=$1 AND user_id=$2`
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	row, err := db.pool.Exec(ctx, query, transactionID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %v", err)
	}
	if row.RowsAffected() == 0 {
		return ErrNotFound
	}
	log.Printf("Transaction with id %d deleted successfully for user %d\n", transactionID, userID)
	return nil
}
