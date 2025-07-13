package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
)

func JsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	resp, err := json.Marshal(data)
	if err != nil {
		JsonError(w, http.StatusBadRequest, fmt.Sprintf("json serialization error: %v", err))
		return
	}
	w.Write(resp)
}

func JsonError(w http.ResponseWriter, status int, errorMessage string) {
	JsonResponse(w, status, models.ErrorResponse{Message: errorMessage})
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	transaction := models.Transaction{}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "request reading error")
		return
	}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		JsonError(w, http.StatusBadRequest, "ivalid json format")
		return
	}
	ctx := context.Background()
	id, err := db.AddTransaction(ctx, &transaction)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error adding transactions")
		return
	}
	transaction.CreatedAt = time.Now()
	transaction.ID = id
	resp := models.SuccessResponse{
		Message: fmt.Sprintf("transactions with id: %d added successfuly", id),
		Data:    transaction,
	}

	JsonResponse(w, http.StatusCreated, resp)
}
