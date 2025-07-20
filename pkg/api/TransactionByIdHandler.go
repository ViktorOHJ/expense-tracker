package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
)

func (s *Server) TransactionByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(r.URL.Query().Get("id"))
	if idStr == "" {
		JsonError(w, http.StatusBadRequest, "id cannot be empty")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 0 {
		JsonError(w, http.StatusBadRequest, "id must be a positive number")
		return
	}

	transaction, err := db.GetTransactionByID(s.db, r.Context(), id)

	if errors.Is(err, db.ErrNotFound) {
		JsonError(w, http.StatusNotFound, fmt.Sprintf("transaction with id %d not found", id))
		return
	}
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error retrieving transaction")
		return
	}

	resp := models.SuccessResponse{
		Message: fmt.Sprintf("transaction with id: %d successfully retrieved", id),
		Data:    transaction,
	}
	JsonResponse(w, http.StatusOK, resp)
}
