package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
	"github.com/ViktorOHJ/expense-tracker/pkg/auth"
	"github.com/ViktorOHJ/expense-tracker/pkg/db"
)

type Server struct {
	db              db.DB
	jwtService      *auth.JWTService
	passwordService *auth.PasswordService
}

func NewServer(db db.DB, jwtService *auth.JWTService, passwordService *auth.PasswordService) *Server {
	return &Server{
		db:              db,
		jwtService:      jwtService,
		passwordService: passwordService,
	}
}

func (s *Server) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	// Публичные маршруты
	mux.HandleFunc("/auth/register", s.RegisterHandler)
	mux.HandleFunc("/auth/login", s.LoginHandler)

	// Защищенные маршруты
	mux.HandleFunc("/transactions", s.AuthMiddleware(s.TransactionHandler))
	mux.HandleFunc("/transaction/", s.AuthMiddleware(s.DeleteGetHandler))
	mux.HandleFunc("/categories", s.AuthMiddleware(s.CategoriesHandler))
	mux.HandleFunc("/summary", s.AuthMiddleware(s.SummaryHandler))

	return mux
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
