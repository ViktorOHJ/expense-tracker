package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func AddCategory(parentContext context.Context, c *models.Category) (models.Category, error) {
	if DB == nil {
		return models.Category{}, errors.New("DB = nil")
	}
	query := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING *`
	ctx, cancel := context.WithTimeout(parentContext, 5*time.Second)
	defer cancel()
	category := models.Category{}
	err := DB.QueryRow(ctx, query, c.Name, c.Description).Scan(&category.ID, &category.Name, &category.Description)
	if err != nil {
		log.Printf("failed to insert transaction: %v", err)
		return models.Category{}, fmt.Errorf("failed to insert transaction: %v", err)
	}
	return category, nil
}
