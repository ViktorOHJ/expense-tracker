package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	models "github.com/ViktorOHJ/expense-tracker/pkg"
)

func InitAPI() {
	http.HandleFunc("/transactions", TansactionHandler)
	http.HandleFunc("/categories", CategoriesHandler)
}

func TansactionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddHandler(w, r)
	case http.MethodGet:
		GetHandler(w, r)
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
