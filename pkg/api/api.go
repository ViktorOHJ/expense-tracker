package api

import (
	"net/http"
)

func InitAPI() {
	http.HandleFunc("/transactions", AddHandler)
	http.HandleFunc("/categories", CategoriesHandler)
}
