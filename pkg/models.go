package models

import (
	"time"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type Transaction struct {
	ID         int       `json:"id"`
	IsIncome   bool      `json:"is_income"`
	Amount     float64   `json:"amount"`
	CategoryID int       `json:"category_id"`
	Note       string    `json:"note,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type Summary struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	Balance      float64 `json:"balance"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
