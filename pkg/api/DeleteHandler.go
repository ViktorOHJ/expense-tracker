package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimSpace(r.URL.Query().Get("id"))
	if idStr == "" {
		JsonError(w, http.StatusBadRequest, "id cannot be empty")
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		JsonError(w, http.StatusBadRequest, "id must be a positive number")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	err = db.DeleteTransaction(ctx, id)
	if errors.Is(err, db.ErrNotFound) {
		JsonError(w, http.StatusNotFound, "transaction not found")
		return
	}
	if err != nil {
		log.Printf("failed to delete transaction: %v", err)
		JsonError(w, http.StatusInternalServerError, "failed to delete transaction")
		return
	}
	JsonResponse(w, http.StatusOK, models.SuccessResponse{
		Message: fmt.Sprintf("Transaction with id %d deleted successfully", id),
	})
}
