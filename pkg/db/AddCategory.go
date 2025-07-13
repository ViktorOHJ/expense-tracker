package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func AddCategory(parentContext context.Context, c *models.Category) (int, error) {
	if DB == nil {
		return 0, errors.New("DB = nil")
	}
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	ctx, cancel := context.WithTimeout(parentContext, 5*time.Second)
	defer cancel()
	var id int
	err := DB.QueryRow(ctx, query, c.Name, c.Description).Scan(&id)
	if err != nil {
		log.Printf("failed to insert transaction: %v", err)
		return 0, fmt.Errorf("failed to insert transaction: %v", err)
	}
	return id, nil
}
