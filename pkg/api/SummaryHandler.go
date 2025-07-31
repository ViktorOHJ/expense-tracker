package api

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func (s *Server) SummaryHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		JsonError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	fromStr := strings.TrimSpace(r.URL.Query().Get("from"))
	toStr := strings.TrimSpace(r.URL.Query().Get("to"))

	if fromStr == "" || toStr == "" {
		JsonError(w, http.StatusBadRequest, "from and to parameters are required")
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		JsonError(w, http.StatusBadRequest, fmt.Sprintf("invalid from date format: %v", err))
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		JsonError(w, http.StatusBadRequest, fmt.Sprintf("invalid to date format: %v", err))
		return
	}

	if to.Before(from) {
		JsonError(w, http.StatusBadRequest, "to date must be after from date")
		return
	}

	summary, err := s.db.GetSummary(r.Context(), user.UserID, from, to)
	if err != nil {
		log.Printf("error retrieving summary: %v", err)
		JsonError(w, http.StatusInternalServerError, "error retrieving summary")
		return
	}

	resp := models.SuccessResponse{
		Message: "Summary retrieved successfully",
		Data:    summary,
	}
	JsonResponse(w, http.StatusOK, resp)
}
