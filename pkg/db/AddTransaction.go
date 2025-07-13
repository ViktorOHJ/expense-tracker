package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func AddTransaction(parentCtx context.Context, t *models.Transaction) (int, error) {
	if DB == nil {
		log.Println("db = nil")
		return 0, errors.New("DB = nil")
	}
	query := `INSERT INTO transactions (is_income, amount, category_id, note) VALUES ($1, $2, $3, $4)
RETURNING id
	`
	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()
	var id int
	err := DB.QueryRow(ctx, query, t.IsIncome, t.Amount, t.CategoryID, t.Note).Scan(&id)
	if err != nil {
		log.Printf("failed to insert transaction: %v", err)
		return 0, fmt.Errorf("failed to insert transaction: %v", err)
	}
	return id, nil
}
