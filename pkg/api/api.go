package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	db *pgxpool.Pool
}

func NewServer(db *pgxpool.Pool) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) InitRoutes() {
	http.HandleFunc("/transactions", s.TransactionHandler)
	http.HandleFunc("/transaction/", s.DeleteGetHandler)
	http.HandleFunc("/categories", s.CategoriesHandler)
	http.HandleFunc("/summary", s.SummaryHandler)
}

func (s *Server) DeleteGetHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.TransactionByIdHandler(w, r)
	case http.MethodDelete:
		s.DeleteHandler(w, r)
	default:
		JsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) TransactionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.AddHandler(w, r)
	case http.MethodGet:
		s.GetHandler(w, r)
	default:
		JsonError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

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
