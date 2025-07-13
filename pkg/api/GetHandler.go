package api

import (
	"net/http"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	transactions := []models.Transaction{}

	err := db.GetTransactions(r.Context(), &transactions)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error retrieving transactions")
		return
	}
	resp := models.SuccessResponse{
		Message: "transactions retrieved successfully",
		Data:    transactions,
	}
	JsonResponse(w, http.StatusOK, resp)
}
