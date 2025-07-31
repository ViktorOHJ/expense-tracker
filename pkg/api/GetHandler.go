package api

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {

	user := GetUserFromContext(r.Context())
	if user == nil {
		JsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

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
			JsonError(w, http.StatusBadRequest, "invalid type parameter")
			return
		}
	}

	categoryID := strings.TrimSpace(q.Get("category_id"))
	var categoryInt *int
	if categoryID != "" {
		c, err := strconv.Atoi(categoryID)
		if err != nil {
			JsonError(w, http.StatusBadRequest, "invalid category_id format")
			return
		}
		categoryInt = &c
	}

	// Парсинг дат (без изменений)
	from := strings.TrimSpace(q.Get("from"))
	var fromTime, toTime *time.Time
	if from != "" {
		f, err := time.Parse("2006-01-02", from)
		if err != nil {
			JsonError(w, http.StatusBadRequest, "invalid date format for 'from'")
			return
		}
		fromTime = &f
	}

	to := strings.TrimSpace(q.Get("to"))
	if to != "" {
		t, err := time.Parse("2006-01-02", to)
		if err != nil {
			JsonError(w, http.StatusBadRequest, "invalid date format for 'to'")
			return
		}
		toTime = &t
	}

	limitStr := strings.TrimSpace(q.Get("limit"))
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		JsonError(w, http.StatusBadRequest, "invalid limit parameter")
		return
	}

	pageStr := strings.TrimSpace(q.Get("offset"))
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 0 {
		JsonError(w, http.StatusBadRequest, "invalid offset parameter")
		return
	}
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit

	transactions, err := s.db.GetTransactions(r.Context(), user.UserID, typeBool, categoryInt, fromTime, toTime, limit, offset)
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
