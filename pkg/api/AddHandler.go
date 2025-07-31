package api

import (
	"encoding/json"
	"io"
	"net/http"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func (s *Server) AddHandler(w http.ResponseWriter, r *http.Request) {

	user := GetUserFromContext(r.Context())
	if user == nil {
		JsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

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

	exists, err := s.db.CheckCategory(ctx, user.UserID, transaction.CategoryID)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "database error during category check")
		return
	}
	if !exists {
		JsonError(w, http.StatusBadRequest, "category does not exist or access denied")
		return
	}

	transaction, err = s.db.AddTransaction(ctx, user.UserID, &transaction)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error adding transaction")
		return
	}

	resp := models.SuccessResponse{
		Message: "transaction added successfully",
		Data:    transaction,
	}
	JsonResponse(w, http.StatusCreated, resp)
}
