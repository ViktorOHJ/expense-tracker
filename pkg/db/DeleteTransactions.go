package db

import (
	"context"
	"fmt"
	"log"
)

func DeleteTransaction(ctx context.Context, id int) error {
	query := `DELETE FROM transactions WHERE id=$1`
	row, err := DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %v", err)
	}
	if row.RowsAffected() == 0 {
		return ErrNotFound
	}
	log.Printf("Transaction with id %d deleted successfully\n", id)
	return nil
}
