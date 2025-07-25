package db

import (
	"context"
	"fmt"
	"log"
	"time"
)

func (db *PostgresDB) DeleteTransaction(parentCtx context.Context, id int) error {
	query := `DELETE FROM transactions WHERE id=$1`
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()

	row, err := db.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %v", err)
	}
	if row.RowsAffected() == 0 {
		return ErrNotFound
	}
	log.Printf("Transaction with id %d deleted successfully\n", id)
	return nil
}
