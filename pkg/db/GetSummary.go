package db

import (
	"context"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetSummary(db *pgxpool.Pool, parentCtx context.Context, from, to time.Time) (models.Summary, error) {
	query := `
		SELECT
			SUM(CASE WHEN is_income THEN amount ELSE 0 END) AS total_income,
			SUM(CASE WHEN NOT is_income THEN amount ELSE 0 END) AS total_expense,
			SUM(CASE WHEN is_income THEN amount ELSE 0 END) - SUM(CASE WHEN NOT is_income THEN amount ELSE 0 END) AS balance
		FROM transactions
		WHERE created_at >= $1 AND created_at <= $2`

	ctx, cancel := context.WithTimeout(parentCtx, 5*time.Second)
	defer cancel()
	var summary models.Summary
	err := db.QueryRow(ctx, query, from, to).Scan(&summary.TotalIncome, &summary.TotalExpense, &summary.Balance)
	if err != nil {
		log.Printf("failed to retrieve summary: %v", err)
		return models.Summary{}, err
	}
	log.Printf("Summary retrieved successfully: %+v", summary)
	return summary, nil
}
