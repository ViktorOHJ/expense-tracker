package api

import (
	"encoding/json"
	"io"
	"net/http"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
)

func (s *Server) AddHandler(w http.ResponseWriter, r *http.Request) {
	transaction := models.Transaction{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "request reading error")
		return
	}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "invalid json format")
		return
	}

	if transaction.Amount <= 0 {
		JsonError(w, http.StatusBadRequest, "amount must be greater than 0")
		return
	}
	ctx := r.Context()
	exists, err := db.CheckCategory(s.db, ctx, transaction.CategoryID)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "database error during category check")
		return
	}
	if !exists {
		JsonError(w, http.StatusBadRequest, "category does not exist")
		return
	}

	transaction, err = db.AddTransaction(s.db, ctx, &transaction)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error adding transactions")
		return
	}

	resp := models.SuccessResponse{
		Message: "transactions added successfuly",
		Data:    transaction,
	}
	JsonResponse(w, http.StatusCreated, resp)
}
