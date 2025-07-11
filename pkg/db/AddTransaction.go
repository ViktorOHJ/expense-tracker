package db

import (
	"context"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func AddTransaction(parentCtx context.Context, t *models.Transaction) (err error) {
	if DB == nil {
		log.Println("db = nil")
		return err
	}
	query := `INSERT INTO transactions (is_income, amount, category_id, note) VALUES ($1, $2, $3, $4)`
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()
	_, err = DB.Exec(ctx, query, t.IsIncome, t.Amount, t.CategoryID, t.Note)
	if err != nil {
		log.Printf("failed to insert transaction: %v", err)
		return err
	}
	return nil
}
