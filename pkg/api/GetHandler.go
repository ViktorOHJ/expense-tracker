package api

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	txType := strings.TrimSpace(q.Get("type"))
	var typeBool *bool
	if txType != "" {
		switch txType {
		case "true":
			typeBool = new(bool)
			*typeBool = true
		case "false":
			typeBool = new(bool)
			*typeBool = false
		default:
			log.Printf("invalid type parameter: %s", txType)
			JsonError(w, http.StatusBadRequest, "invalid type parameter")
			return
		}
	}
	categoryID := strings.TrimSpace(q.Get("category_id"))
	var categoryInt *int
	if categoryID != "" {
		c, err := strconv.Atoi(categoryID)
		if err != nil {
			log.Printf("error converting category_id to int: %v", err)
			JsonError(w, http.StatusBadRequest, "invalid category_id format")
			return
		}
		categoryInt = &c
	}

	from := strings.TrimSpace(q.Get("from"))
	var fromTime, toTime *time.Time
	if from != "" {
		f, err := time.Parse("2006-01-02", from)
		if err != nil {
			log.Printf("error parsing 'from' date: %v", err)
			JsonError(w, http.StatusBadRequest, "invalid date format for 'from'")
			return
		}
		fromTime = &f
	}
	to := strings.TrimSpace(q.Get("to"))
	if to != "" {
		t, err := time.Parse("2006-01-02", to)
		if err != nil {
			log.Printf("error parsing 'to' date: %v", err)
			JsonError(w, http.StatusBadRequest, "invalid date format for 'to'")
			return
		}
		toTime = &t
	}

	transactions, err := db.GetTransactions(r.Context(), typeBool, categoryInt, fromTime, toTime)
	if err != nil {
		JsonError(w, http.StatusInternalServerError, "error retrieving transactions")
		return
	}
	var resp models.SuccessResponse
	if len(transactions) == 0 {
		resp = models.SuccessResponse{
			Message: "no transactions found",
			Data:    []models.Transaction{},
		}
	} else {
		resp = models.SuccessResponse{
			Message: "transactions listed successfully",
			Data:    transactions,
		}
	}
	JsonResponse(w, http.StatusOK, resp)
}
